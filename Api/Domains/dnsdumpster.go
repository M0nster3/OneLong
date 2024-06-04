// Package dnsdumpster logic
package Domains

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

// 用于保护 addedURLs
func GetEnInfoPostBuffer(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Digitorus"
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

	Utils.DomainTableShow(keyword, data, "dnsdumpster")

	return ensInfos, ensOutMap

}

func PostBuffer(domain string, options *Utils.LongOptions, token string, client *resty.Client, DomainsIP *outputfile.DomainsIP) string {
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
	time.Sleep(3 * time.Second)
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

			// 编译正则表达式
			re := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)

			// 查找匹配的内容
			matches := re.FindAllStringSubmatch(strings.TrimSpace(td3.Text()), -1)
			if len(matches) > 0 {
				for _, bu := range matches {
					Hostname = append(Hostname, bu[0])
				}
			}

		})
		s.Find(".col-md-4").Each(func(j int, td4 *goquery.Selection) {
			// 将每个单元格的内容追加到 content 变量中

			hostname := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
			// 编译正则表达式
			re := regexp.MustCompile(hostname)

			// 查找匹配的内容
			matches := re.FindAllStringSubmatch(strings.TrimSpace(td4.Text()), -1)
			if len(matches) > 0 {
				for _, bu := range matches {
					Address = append(Address, bu[0])

				}
			}

		})

	})
	var result string
	result = "{\"passive_dns\":["
	for i := 0; i < len(Hostname) && i < len(Address); i++ {
		result += "{\"address\"" + ":" + "\"" + Hostname[i] + "\"" + "," + "\"hostname\"" + ":" + "\"" + Address[i] + "\"" + "},"
		Address[i] = strings.TrimRight(Address[i], ".")
		DomainsIP.IP = append(DomainsIP.IP, Hostname[i])
		DomainsIP.Domains = append(DomainsIP.Domains, Address[i])
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfoPostBuffer(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Dnsdumpster Dns查询", options)

	return ""
}

func Dnsdumpster(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Dnsdumpster  API DNS反差域名 \n")
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
		//"X-Key":      {options.LongConfig.Cookies.Binaryedge},
	}

	client.Header.Set("Content-Type", "application/json")
	//client.Header.Del("Cookie")

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
		gologger.Errorf("Dnsdumpster API 链接访问失败尝试切换代理\n")
		return ""
	}
	SetCookie := resp.Header().Get("Set-Cookie")
	csrftoken := strings.Split(SetCookie, ";")
	csrftoken[0] = strings.ReplaceAll(csrftoken[0], "csrftoken=", "")
	res := PostBuffer(domain, options, csrftoken[0], client, DomainsIP)
	if res == "Kill" {
		//gologger.Labelf("Dnsdumpster  Api 未发现域名 %s\n", domain)
		return ""
	}

	return "Success"
}
