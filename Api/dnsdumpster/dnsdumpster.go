// Package dnsdumpster logic
package dnsdumpster

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	//respons := gjson.Get(response, "events").Array()
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "[", "")
	//respons := gjson.Parse(response).Array()
	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Digitorus"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
	for aa, _ := range respons {
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(respons[aa].String()))
	}
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func PostBuffer(domain string, options *Utils.ENOptions, token string, client *resty.Client, DomainsIP *outputfile.DomainsIP) string {
	urls := "https://dnsdumpster.com/"
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		//client.SetProxy(options.Proxy)
		client.SetProxy("http://192.168.203.111:808")

	}
	client.Header = http.Header{
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":       {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Referer":      {"https://dnsdumpster.com/"},
		"X-CSRF-Token": {token},
	}

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()
	requestBody := fmt.Sprintf("csrfmiddlewaretoken=%s&targetip=%s&user=free", token, domain)
	response, err := clientR.SetBody(requestBody).Post(urls)
	if strings.Contains(string(response.Body()), "Check your query and try again") {
		return "Kill"
	}
	// 使用 HTML 解析器解析 HTML 内容
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		fmt.Println("HTML 解析错误:", err)
		return "Kill"
	}

	// 查找具有特定 class 的元素并获取其内容
	var Hostname []string
	var Address []string

	doc.Find(".table-responsive").Each(func(i int, s *goquery.Selection) {
		s.Find(".col-md-3").Each(func(j int, td3 *goquery.Selection) {

			// 将每个单元格的内容追加到 content 变量中
			//(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}
			//pattern := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
			//IpV4 := `\b(?:\d{1,3}\.){3}\d{1,3}(?:(?!\d)(?:\.\w+)+)?`

			// 编译正则表达式
			re := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)

			// 查找匹配的内容
			matches := re.FindAllStringSubmatch(strings.TrimSpace(td3.Text()), -1)
			for _, bu := range matches {

				Hostname = append(Hostname, bu[0])

			}

		})
		s.Find(".col-md-4").Each(func(j int, td4 *goquery.Selection) {
			// 将每个单元格的内容追加到 content 变量中

			hostname := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
			// 编译正则表达式
			re := regexp.MustCompile(hostname)

			// 查找匹配的内容
			matches := re.FindAllStringSubmatch(strings.TrimSpace(td4.Text()), -1)
			for _, bu := range matches {
				Address = append(Address, bu[0])

			}

		})

	})
	var result string
	result = "{\"passive_dns\":["
	for i := 0; i < len(Hostname) && i < len(Address); i++ {
		result += "{\"address\"" + ":" + "\"" + Hostname[i] + "\"" + "," + "\"hostname\"" + ":" + "\"" + Address[i] + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, Address[i])
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Dnsdumpster Dns查询", options)
	//res, ensOutMap := GetEnInfo(respjoin)
	//
	//outputfile.MergeOutPut(res, ensOutMap, "Dnsdumpster DNS反差", options)
	return ""
}

func Dnsdumpster(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	gologger.Infof("Dnsdumpster  API DNS反差域名 \n")
	urls := "https://dnsdumpster.com/"
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
		//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
	}

	client.Header.Set("Content-Type", "application/json")
	//client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, _ := clientR.Send()
	for {
		if resp.RawResponse == nil {
			resp, _ = clientR.Send()
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	SetCookie := resp.Header().Get("Set-Cookie")
	csrftoken := strings.Split(SetCookie, ";")
	csrftoken[0] = strings.ReplaceAll(csrftoken[0], "csrftoken=", "")
	res := PostBuffer(domain, options, csrftoken[0], client, DomainsIP)
	if res == "Kill" {
		gologger.Labelf("Dnsdumpster  Api 未发现域名\n")
	}

	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
