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
	"sync"

	//"strconv"
	//"strings"
	"time"
)

// 用于保护 addedURLs
func GetEnInfoCensys(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Parse(response).Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Censys"
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

	Utils.DomainTableShow(keyword, data, "censys")

	return ensInfos, ensOutMap

}

func Censys(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Censys 空间探测\n")
	urls := "https://search.censys.io/api/v2/certificates/search"
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
	var cursor string
	result := "["
	var wg1 sync.WaitGroup
	for add := 1; add < 10; add += 1 {
		requestBody := map[string]interface{}{
			"q":        domain,
			"per_page": 100,
			"cursor":   cursor,
		}
		username := options.LongConfig.Cookies.CensysToken
		password := options.LongConfig.Cookies.CensysSecret
		client.SetBasicAuth(username, password)

		//强制延时1s
		time.Sleep(3 * time.Second)
		//加入随机延迟
		time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
		clientR := client.R()

		clientR.URL = urls
		resp, err := clientR.
			SetBody(requestBody).
			Post(urls)
		for aa := 1; aa < 4; aa += 1 {
			if resp.RawResponse == nil {
				resp, _ = clientR.
					SetBody(requestBody).
					Post(urls)
				time.Sleep(3 * time.Second)
			} else if resp.Body() != nil {
				break
			}
		}

		if err != nil {
			gologger.Errorf("Censys 空间探测访问失败尝试切换代理\n")
			return ""
		}
		if resp.StatusCode() == 401 {
			gologger.Labelf("Censys 空间探测 Token 错误\n")
			return ""
		} else if strings.Contains(string(resp.Body()), "our account is limited to 10 page") {
			gologger.Labelf("Censys 空间探测当前Token只能查询10页 \n")
			return ""
		} else if resp.StatusCode() == 403 {
			gologger.Labelf("Censys 空间探测Cookie失效\n")
			return ""
		} else if gjson.Get(string(resp.Body()), "result.total").Int() == 0 {
			//gologger.Labelf("Censys 空间探测未发现域名 %s\n", domain)
			return ""
		}
		hostname := gjson.Get(string(resp.Body()), "result.hits.#.names").Array()
		for aa, _ := range hostname {
			wg1.Add(1)
			go func() {
				zuo := strings.ReplaceAll(hostname[aa].String(), "[", "")
				zhong := strings.ReplaceAll(zuo, "*.", "")
				you := strings.ReplaceAll(zhong, "]", "")
				host := regexp.MustCompile(`(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`)
				ip := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)
				// 编译正则表达式
				hostrea := host.FindAllStringSubmatch(strings.TrimSpace(you), -1)
				iprea := ip.FindAllStringSubmatch(strings.TrimSpace(you), -1)
				if hostrea != nil && iprea == nil {
					result += you + ","
				}
				wg1.Done()
			}()
		}
		wg1.Wait()
		next := gjson.Get(string(resp.Body()), "result.links.next").String()
		if next == "" {
			break
		}
		add += 1
		cursor = next
	}
	result = result + "]"
	res, ensOutMap := GetEnInfoCensys(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Censys", options)

	return "Success"
}
