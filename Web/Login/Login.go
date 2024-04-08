package Login

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/Login/WaybackarchiveLogin"
	"OneLong/Web/Login/alienvaultLogin"
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"sync"
	"time"
)

func gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}
func GetEnInfo(response string) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	//ensInfos.SType = "Commoncrawl"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}

	for aa, _ := range respons {
		ensInfos.Infos["Login"] = append(ensInfos.Infos["Login"], gjson.Parse(respons[aa].String()))
	}

	return ensInfos, ensOutMap

}
func Login(domains []string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	//color.RGBStyleFromString("244,211,49").Println("\n--------------------探测网站后台--------------------")

	gologger.Infof("扫描网站后台，当前共有存活子域%d个\n", len(domains))
	for domainint, domain := range domains {

		gologger.Infof("当前扫描第%d个  %s\n", domainint+1, domain)
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.ReplaceAll(domain, "http://", "")
		if !strings.Contains(domain, "\\") {
			alienvaultLogin.AlienvaultLogin(domain, options, DomainsIP)
			WaybackarchiveLogin.WaybackarchiveLogin(domain, options, DomainsIP)
			//CommoncrawlLogin.CommoncrawlLogin(domain, options, DomainsIP)
		}
	}

	ParseLoginurl(options, DomainsIP)
}

func ParseLoginurl(options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	var wg sync.WaitGroup
	DomainsIP.LoginUrl = Utils.SetStr(DomainsIP.LoginUrl)
	for _, domain := range DomainsIP.LoginUrl {
		wg.Add(1)
		domain := domain
		go func() {
			urls := domain
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
			//if strings.Contains(urls, "mail") {
			//	fmt.Print(11111)
			//}
			resp, err := clientR.Get(urls)
			if err != nil {
				for attempt := 0; attempt < 4; attempt++ {
					if attempt == 2 { // 在第三次尝试时切换到HTTPS
						urls = strings.ReplaceAll(urls, "http://", "https://")
						if strings.Contains(urls, "https://") {
							urls = strings.ReplaceAll(urls, ":80/", ":443/")
						}
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
			}

			if err != nil {

				wg.Done()
			} else if resp.StatusCode() != 404 {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
				title := doc.Find("title").Text()
				title = strings.Replace(title, "\n", "", -1)
				title = strings.Trim(title, " ")
				DomainsIP.LoginTitle = append(DomainsIP.LoginTitle, title)
				DomainsIP.LoginUrlA = append(DomainsIP.LoginUrlA, urls)
				wg.Done()
			} else {
				wg.Done()
			}

		}()

	}
	wg.Wait()
	passive_dns := "{\"passive_dns\":["
	var add int

	for add = 0; add < len(DomainsIP.LoginUrlA); add++ {
		if add < len(DomainsIP.LoginUrlA) && add < len(DomainsIP.LoginTitle) {
			passive_dns += "{\"hostname\"" + ":" + "\"" + DomainsIP.LoginUrlA[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.LoginTitle[add] + "\"" + "},"
		} else if add < len(DomainsIP.LoginUrlA) {
			passive_dns += "{\"hostname\"" + ":" + "\"" + DomainsIP.LoginUrlA[add] + "\"" + "},"
		}

	}
	passive_dns = passive_dns + "]}"

	res, ensOutMap := GetEnInfo(passive_dns)

	outputfile.MergeOutPut(res, ensOutMap, "Login", options)
}
