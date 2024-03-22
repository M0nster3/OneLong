package Login

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/Login/CommoncrawlLogin"
	"OneLong/Web/Login/alienvaultLogin"
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
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
	for _, domain := range domains {
		CommoncrawlLogin.CommoncrawlLogin(domain, options, DomainsIP)
		alienvaultLogin.AlienvaultLogin(domain, options, DomainsIP)
		ParseLoginurl(options, DomainsIP)
	}

}

func ParseLoginurl(options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	DomainsIP.LoginUrl = Utils.SetStr(DomainsIP.LoginUrl)
	for _, aa := range DomainsIP.LoginUrl {
		urls := aa
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
		time.Sleep(1 * time.Second)
		//加入随机延迟
		time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
		clientR := client.R()

		clientR.URL = urls
		resp, err := clientR.Get(urls)
		if err != nil {
			break
		}
		title := gettitle(string(resp.Body()))
		DomainsIP.LoginTitle = append(DomainsIP.LoginTitle, title)
		DomainsIP.LoginUrlA = append(DomainsIP.LoginUrlA, urls)
	}
	passive_dns := "{\"passive_dns\":["
	var add int
	for add = 0; add < len(DomainsIP.LoginUrlA); add++ {
		passive_dns += "{\"hostname\"" + ":" + "\"" + DomainsIP.LoginUrlA[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.LoginTitle[add] + "\"" + "},"

	}
	passive_dns = passive_dns + "]}"

	res, ensOutMap := GetEnInfo(passive_dns)

	outputfile.MergeOutPut(res, ensOutMap, "Login", options)
}
