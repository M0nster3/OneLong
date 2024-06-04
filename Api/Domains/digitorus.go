package Domains

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strings"
	//"strconv"
	//"strings"
	"time"
)

// 用于保护 addedURLs
func GetEnInfoDigitorus(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Parse(response).Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Digitorus"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	addedURLs := make(map[string]bool)
	for bb, _ := range respons {
		ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + respons[bb].String() + "\"" + "}"
		url := gjson.Parse(ResponseJia).Get("hostname").String()

		// 检查是否已存在相同的 URL
		if !addedURLs[url] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
			DomainsIP.Domains = append(DomainsIP.Domains, url)
			addedURLs[url] = true
		} else {
			fmt.Printf("111")
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

	Utils.DomainTableShow(keyword, data, "Digitorus")

	return ensInfos, ensOutMap

}

func Digitorus(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Digitorus  API 证书查询域名 \n")
	//gologger.Labelf("只实现普通Api 如果是企业修改Api接口 免费的每月250次\n")
	urls := "https://certificatedetails.com/" + domain
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		//"X-Key":      {options.LongConfig.Cookies.Binaryedge},
	}

	client.Header.Set("Content-Type", "application/json")
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
		gologger.Errorf("Digitorus API 链接访问失败尝试切换代理\n")
		return ""
	}
	if string(resp.Body()) == "" || resp.StatusCode() == 404 {
		//gologger.Labelf("Digitorus API 未查询到域名 %s\n", domain)
		return ""
	}
	pattern := `\[c:2+\|t:2+\|false\]\s*(.*)`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 查找匹配的内容
	matches := re.FindAllStringSubmatch(string(resp.Body()), -1)
	if len(matches) > 0 {
		var ia5StringContent []string
		// 遍历匹配项
		ia5StringContent = append(ia5StringContent, "[")
		for _, match := range matches {
			zuo := strings.ReplaceAll(match[1], "IA5String &#39;", "\"")
			you := strings.ReplaceAll(zuo, "&#39;", "\"")
			// 第一个捕获组匹配的内容即为你需要的字符串
			ia5StringContent = append(ia5StringContent, you+",")
		}
		ia5StringContent = append(ia5StringContent, "]")
		respjoin := strings.Join(ia5StringContent, " ")
		res, ensOutMap := GetEnInfoDigitorus(respjoin, DomainsIP)

		outputfile.MergeOutPut(res, ensOutMap, "Digitorus 证书查询", options)
	}

	return "Success"
}
