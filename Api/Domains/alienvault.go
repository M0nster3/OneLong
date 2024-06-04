package Domains

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
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
func GetEnInfoAlienvault(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "alienvault"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	addedURLs := make(map[string]bool)
	for aa, _ := range respons {
		if strings.Contains(respons[aa].String(), "address") {
			re := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)
			ip := gjson.Get(respons[aa].String(), "address").String()
			matches := re.FindAllStringSubmatch(strings.TrimSpace(ip), -1)
			if len(matches) > 0 {
				for _, bu := range matches {
					if !addedURLs[bu[0]] {
						// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
						DomainsIP.IP = append(DomainsIP.IP, bu[0])
						addedURLs[bu[0]] = true
					}
					break
				}
			}

		}
		if strings.Contains(respons[aa].String(), "hostname") {
			hostname := gjson.Get(respons[aa].String(), "hostname").String()
			if !addedURLs[hostname] {
				// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
				DomainsIP.Domains = append(DomainsIP.Domains, hostname)
				addedURLs[hostname] = true
			}

		}
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

	Utils.DomainTableShow(keyword, data, "alienvault")

	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func Alienvault(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Alienvault\n")
	urls := "https://otx.alienvault.com/api/v1/indicators/domain/" + domain + "/passive_dns"
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
		gologger.Errorf("Alienvault API 链接访问失败尝试切换代理\n")
		return ""
	}
	count := gjson.GetBytes(resp.Body(), "count").Int()
	if count == 0 {
		//gologger.Labelf("Alienvault Api 未发现域名 %s\n", domain)
		return ""
	}

	res, ensOutMap := GetEnInfoAlienvault(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "alienvault", options)

	return ""
}
