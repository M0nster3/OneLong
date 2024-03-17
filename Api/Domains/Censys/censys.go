package Censys

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

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Parse(response).Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Censys"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
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
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func Censys(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
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
	add := 0
	for {
		requestBody := map[string]interface{}{
			"q":        domain,
			"per_page": 100,
			"cursor":   cursor,
		}
		username := options.ENConfig.Cookies.CensysToken
		password := options.ENConfig.Cookies.CensysSecret
		client.SetBasicAuth(username, password)

		//强制延时1s
		time.Sleep(1 * time.Second)
		//加入随机延迟
		time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
		clientR := client.R()

		clientR.URL = urls
		resp, _ := clientR.
			SetBody(requestBody).
			Post(urls)
		for {
			if resp.RawResponse == nil {
				resp, _ = clientR.
					SetBody(requestBody).
					Post(urls)
				time.Sleep(1 * time.Second)
			} else if resp.Body() != nil {
				break
			}
		}
		if resp.StatusCode() == 401 {
			gologger.Labelf("Censys 空间探测 Token 错误\n")
			return ""
		} else if strings.Contains(string(resp.Body()), "our account is limited to 10 page") {
			gologger.Labelf("Censys 空间探测当前Token只能查询10页 \n")
			return ""
		} else if gjson.Get(string(resp.Body()), "result.total").Int() == 0 {
			gologger.Labelf("Censys 空间探测未发现域名 %s\n", domain)
			return ""
		}
		hostname := gjson.Get(string(resp.Body()), "result.hits.#.names").Array()
		for aa, _ := range hostname {
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

		}
		next := gjson.Get(string(resp.Body()), "result.links.next").String()
		if next == "" || add > 10 {
			break
		}
		add += 1
		cursor = next
	}
	result = result + "]"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Censys", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
