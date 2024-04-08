package Domains

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"strings"

	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

// 用于保护 addedURLs
func GetEnInfoBevigil(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "subdomains").Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Anubis"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	addedURLs := make(map[string]bool)
	for aa, _ := range respons {
		ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + respons[aa].String() + "\"" + "}"
		url := gjson.Parse(ResponseJia).Get("hostname").String()

		// 检查是否已存在相同的 URL
		if !addedURLs[url] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
			DomainsIP.Domains = append(DomainsIP.Domains, url)
			addedURLs[url] = true
		}

	}

	//命令输出展示
	var data [][]string
	var keyword []string
	for _, y := range GetENMap() {
		for _, ss := range y.KeyWord {
			if ss == "数据关联" {
				continue
			}
			keyword = append(keyword, ss)
		}

		for _, res := range ensInfos.Infos["Urls"] {
			results := gjson.GetMany(res.Raw, y.Field...)
			var str []string
			for _, s := range results {
				str = append(str, s.String())
			}
			data = append(data, str)
		}

	}

	Utils.DomainTableShow(keyword, data, "bevigil")

	return ensInfos, ensOutMap

}

func Bevigil(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Bevigil 空间探测\n")
	urls := "https://osint.bevigil.com/api/" + domain + "/subdomains/"
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent":     {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":         {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"X-Access-Token": {options.LongConfig.Cookies.Bevigil},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(3 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
	for add := 1; add < 4; add += 1 {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(3 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}

	if err != nil {
		gologger.Errorf("Bevigil 空间探测链接访问失败尝试切换代理\n")
		return ""
	}
	if len(gjson.Get(string(resp.Body()), "subdomains").Array()) == 0 || len(gjson.Get(string(resp.Body()), "subdomains").Array()) == 1 {
		//gologger.Labelf("Bevigil 空间探测未发现域名 %s\n", domain)
		return ""
	}
	if strings.Contains(string(resp.Body()), "Access Token Invalid") {
		gologger.Labelf("Anubis Cookie 不正确\n")
		return ""
	}
	res, ensOutMap := GetEnInfoBevigil(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Bevigil", options)

	return "Success"
}
