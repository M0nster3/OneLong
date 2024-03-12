package netlas

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
	"strconv"
	"strings"

	//"strconv"
	//"strings"
	"time"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "passive_dns").Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Netlas"
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

func Netlas(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {

	gologger.Infof("Netlas 威胁平台查询\n")
	//urls := "https://leakix.net/api/subdomains/" + domain
	endpoint := "https://app.netlas.io/api/domains_count/"
	params := url.Values{}
	countQuery := fmt.Sprintf("domain:*.%s AND NOT domain:%s", domain, domain)
	params.Set("q", countQuery)
	urls := endpoint + "?" + params.Encode()

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"api-key":    {options.ENConfig.Cookies.Netlas},
	}

	client.Header.Del("Cookie")

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
			time.Sleep(2 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if gjson.Get(string(resp.Body()), "count").Int() == 0 {
		gologger.Labelf("Netlas 威胁平台未发现域名\n")
		return ""
	}
	if strings.Contains(string(resp.Body()), "You can wait while daily rate limit will ") {
		gologger.Errorf("请切换IP 当日访问上限")
	}
	count := gjson.Get(string(resp.Body()), "count").Int()
	var address []string
	var hostname []string
	var ipss string
	for i := 0; i < int(count); i += 20 {
		endpoint := "https://app.netlas.io/api/domains/"
		params := url.Values{}
		offset := strconv.Itoa(i)
		query := fmt.Sprintf("domain:(domain:*.%s AND NOT domain:%s)", domain, domain)
		params.Set("q", query)
		params.Set("source_type", "include")
		params.Set("start", offset)
		params.Set("fields", "*")
		apiUrl := endpoint + "?" + params.Encode()
		client.Header = http.Header{
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
			"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
			"api-key":    {"R5yHhXQgud0eDV34IR8TUck3AchS99dS"},
		}
		clientR := client.R()

		clientR.URL = apiUrl
		resp, _ := clientR.Send()
		buff := gjson.Get(string(resp.Body()), "items.#.data").Array()
		for _, item := range buff {
			ips := item.Get("a").Array()
			hostnames := item.Get("domain").String()
			for _, ip := range ips {
				ipss = ipss + ip.String() + "\n"
			}
			address = append(address, ipss)
			hostname = append(hostname, hostnames)
			ipss = ""

		}
	}

	passive_dns := "{\"passive_dns\":["
	var add int
	for add = 0; add < len(hostname); add++ {
		passive_dns += "{\"hostname\"" + ":" + "\"" + hostname[add] + "\"" + "," + "\"address\"" + ":" + "\"" + address[add] + "\"" + "},"
		ips := strings.Split(address[add], "\n")
		for _, ip := range ips {
			DomainsIP.IP = append(DomainsIP.IP, ip)
		}
		DomainsIP.Domains = append(DomainsIP.Domains, hostname[add])

	}
	passive_dns = passive_dns + "]}"
	res, ensOutMap := GetEnInfo(passive_dns, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Netlas", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
