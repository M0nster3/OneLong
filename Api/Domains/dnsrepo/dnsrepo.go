package dnsrepo

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
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
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func Dnsrepo(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Dnsrepo API 域名查询 \n")
	urls := "https://dnsrepo.noc.org/?search=" + domain
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
	resp, err := clientR.Get(urls)
	for add := 1; add < 4; add += 1 {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if err != nil {
		gologger.Errorf("Dnsrepo API 链接访问失败尝试切换代理\n")
		return ""
	}
	if resp.Body() == nil || strings.Contains(string(resp.Body()), "nothing found") {
		gologger.Labelf("Dnsrepo API 未查询到域名 %s\n", domain)
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		fmt.Println("HTML 解析错误:", err)
		return ""
	}

	// 查找具有特定 class 的元素并获取其内容
	var Hostname []string
	var Address []string
	doc.Find(".table-responsive").Each(func(i int, s *goquery.Selection) {

		s.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {

			// 查找当前 tr 元素下的所有 td 元素，并提取文本内容
			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				s.Find("a[href^='/?domain=']").Each(func(i int, s *goquery.Selection) {
					href, _ := s.Attr("href")
					href = strings.ReplaceAll(href, "/?domain=", "")
					Hostname = append(Hostname, href)
				})
				s.Find("a[href^='/?ip=']").Each(func(i int, s *goquery.Selection) {
					href, _ := s.Attr("href")
					href = strings.ReplaceAll(href, "/?ip=", "")
					re := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)

					// 查找匹配的内容
					matches := re.FindAllStringSubmatch(href, -1)
					for _, bu := range matches {
						Address = append(Address, bu[0])
					}
				})

			})
		})

	})
	Hostname = Utils.SetStr(Hostname)
	Address = Utils.SetStr(Address)
	var result string
	result = "{\"passive_dns\":["
	var add int
	if len(Hostname) < len(Address) {
		for add = 0; add < len(Hostname); add++ {
			result += "{\"hostname\"" + ":" + "\"" + Hostname[add] + "\"" + "," + "\"address\"" + ":" + "\"" + Address[add] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, Hostname[add])
			DomainsIP.IP = append(DomainsIP.IP, Address[add])
		}
		for ii := add; ii < len(Address); ii++ {
			result += "{\"address\"" + ":" + "\"" + Address[ii] + "\"" + "},"
			DomainsIP.IP = append(DomainsIP.IP, Address[ii])
		}

	} else {
		for add = 0; add < len(Address); add++ {
			result += "{\"hostname\"" + ":" + "\"" + Hostname[add] + "\"" + "," + "\"address\"" + ":" + "\"" + Address[add] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, Hostname[add])
			DomainsIP.IP = append(DomainsIP.IP, Address[add])
		}
		for ii := add; ii < len(Hostname); ii++ {
			result += "{\"hostname\"" + ":" + "\"" + Hostname[ii] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, Hostname[ii])
		}
	}

	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Dnsrepo API 域名查询", options)
	return "Success"
}
