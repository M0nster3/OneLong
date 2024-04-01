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
	"net/url"
	"strings"
	//"strconv"
	//"strings"
	"time"
)

// 用于保护 addedURLs
func GetEnInfoWaybackarchive(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	var result []string
	responselist := strings.Split(response, "\n")
	//responselist := gjson.Get(response, "result.records.#.domain").Array()
	for _, urlString := range responselist {
		u, err := url.Parse(urlString)
		if err != nil {
			//fmt.Printf("Error parsing URL '%s': %v\n", urlString, err)
			continue
		}
		result = append(result, u.Hostname())
	}
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "waybackarchive"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	addedURLs := make(map[string]bool)
	for aa, _ := range result {
		ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + result[aa] + "\"" + "}"
		urls := gjson.Parse(ResponseJia).Get("hostname").String()

		// 检查是否已存在相同的 URL
		if !addedURLs[urls] {
			DomainsIP.Domains = append(DomainsIP.Domains, urls)
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
			addedURLs[urls] = true
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

	Utils.DomainTableShow(keyword, data, "waybackarchive")

	return ensInfos, ensOutMap

}

func Waybackarchive(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("waybackarchive 历史快照查询\n")

	urls := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original&collapse=urlkey", domain)
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(3 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
	for attempt := 0; attempt < 4; attempt++ {
		if attempt == 2 { // 在第三次尝试时切换到HTTPS
			break
		}
		resp, err = client.R().Get(urls)
		if err != nil || resp == nil || resp.RawResponse == nil {
			time.Sleep(3 * time.Second) // 在重试前等待
			continue
		}
		// 如果得到有效响应，处理响应
		if resp.Body() != nil {
			// 处理响应的逻辑
			break
		}
	}
	if err != nil {
		gologger.Errorf("waybackarchive 历史快照链接访问失败尝试切换代理\n")
		return ""
	}
	if resp.Body() == nil {
		//gologger.Labelf("waybackarchive 历史快照未发现域名 %s\n", domain)
		return ""
	} else if resp.Size() == 0 {
		//gologger.Labelf("waybackarchive 历史快照未发现域名 %s\n", domain)
		return ""
	}
	res, ensOutMap := GetEnInfoWaybackarchive(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "waybackarchive 历史快照", options)

	return "Success"
}
