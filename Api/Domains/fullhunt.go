package Domains

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

// 用于保护 addedURLs
func GetEnInfoFullhunt(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "hosts").Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Fullhunt"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	addedURLs := make(map[string]bool)
	for aa, _ := range respons {
		ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + respons[aa].String() + "\"" + "}"
		url := gjson.Parse(ResponseJia).Get("hostname").String()
		DomainsIP.Domains = append(DomainsIP.Domains, url)
		// 检查是否已存在相同的 URL
		if !addedURLs[url] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
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

	Utils.DomainTableShow(keyword, data, "Fullhunt")

	return ensInfos, ensOutMap

}

func Fullhunt(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Fullhunt 威胁平台查询\n")
	urls := "https://fullhunt.io/api/v1/domain/" + domain + "/subdomains"
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"X-Api-Key":  {options.LongConfig.Cookies.Fullhunt},
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
		gologger.Errorf("Fullhunt 威胁平台链接访问失败尝试切换代理\n")
		return ""
	}
	if resp.StatusCode() == 404 {
		//gologger.Labelf("Fullhunt 威胁平台未发现域名\n")
		return ""
	} else if gjson.Get(string(resp.Body()), "status").Int() == 400 {
		gologger.Errorf("Fullhunt 威胁平台次数用完\n")
		return ""
	}
	num := gjson.Get(string(resp.Body()), "hosts").Array()
	if len(num) == 0 {
		return ""
	}
	res, ensOutMap := GetEnInfoFullhunt(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Fullhunt", options)

	return "Success"
}
