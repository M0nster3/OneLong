// Package commoncrawl logic
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
	"time"
)

// 用于保护 addedURLs
func GetEnInfoCommoncrawl(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Commoncrawl"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	for aa, _ := range respons {
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(respons[aa].String()))
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

	Utils.DomainTableShow(keyword, data, "commoncrawl")

	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func Commoncrawl(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Commoncrawl  API 查询 \n")
	//gologger.Labelf("只实现普通Api 如果是企业修改Api接口 免费的每月250次\n")
	urls := "https://index.commoncrawl.org/collinfo.json"
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	client.Header = http.Header{
		"User-Agent": {Utils.RandUA()},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
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
		if resp.RawResponse == nil || resp.StatusCode() == 503 {
			resp, _ = clientR.Get(urls)
			time.Sleep(4 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if err != nil || resp.StatusCode() == 503 {
		gologger.Errorf("Commoncrawl API 链接访问失败尝试切换代理\n")
		return ""
	}
	buff := gjson.Parse(string(resp.Body())).Array()
	var result []string
	addedURLs := make(map[string]bool)

	for aa, item := range buff {
		if aa == 4 {
			break
		}
		// 从当前条目获取域名
		cdx := item.Get("cdx-api").String()
		url := fmt.Sprintf("%s?url=*.%s", cdx, domain)
		clientss := resty.New()
		clienta := clientss.R()
		clienta.Header = http.Header{
			"User-Agent": {Utils.RandUA()},
			"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
		}

		clienta.Header.Set("Content-Type", "application/json")
		clienta.Header.Del("Cookie")
		clienta.URL = url
		time.Sleep(3 * time.Second)
		respa, err := clienta.Get(url)
		for add := 1; add < 20; add += 1 {
			if err != nil || respa.StatusCode() == 503 {
				clients := resty.New()
				clientaa := clients.R()
				clientaa.Header = http.Header{
					"User-Agent": {Utils.RandUA()},
					"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
					//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
				}

				clientaa.Header.Set("Content-Type", "application/json")
				respa, _ = clientaa.Get(url)
				time.Sleep(3 * time.Second)
			} else if resp.Body() != nil {
				break
			}
		}
		if respa.StatusCode() == 503 {
			fmt.Printf("503")
			return ""
		}
		if respa.StatusCode() == 404 {
			return ""
		}

		hostname := `(?:[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?\.)+` + regexp.QuoteMeta(domain)
		// 编译正则表达式
		re := regexp.MustCompile(hostname)

		// 查找匹配的内容
		matches := re.FindAllStringSubmatch(strings.TrimSpace(string(respa.Body())), -1)
		for _, bu := range matches {
			if !addedURLs[bu[0]] {
				// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
				result = append(result, bu[0])
				addedURLs[bu[0]] = true
			}
		}

	}

	passive_dns := "{\"passive_dns\":["
	var add int
	for add = 0; add < len(result); add++ {
		passive_dns += "{\"hostname\"" + ":" + "\"" + result[add] + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, result[add])
	}
	passive_dns = passive_dns + "]}"
	if len(DomainsIP.Domains) == 0 {

		gologger.Labelf("Commoncrawl API 未查询到域名 %s\n", domain)
		return ""
	}
	res, ensOutMap := GetEnInfoCommoncrawl(passive_dns, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Commoncrawl Api", options)

	return "Success"
}
