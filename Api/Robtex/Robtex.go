package Robtex

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

	//"strconv"
	//"strings"
	"time"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "passive_dns").Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Robtex	"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))oomEye

	for aa, _ := range respons {
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(respons[aa].String()))
	}
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func Robtex(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	gologger.Infof("Robtex Api查詢\n")

	urls := fmt.Sprintf("https://freeapi.robtex.com/pdns/forward/%s", domain)
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":       {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"Content-Type": {"application/x-ndjson"},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, _ := clientR.Send() //ratelimited
	for {
		if resp.RawResponse == nil {
			resp, _ = clientR.Send()
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if len(resp.Body()) == 0 {
		gologger.Labelf("Robtex Api 未发现域名\n")
		return ""
	}
	if strings.Contains(string(resp.Body()), "ratelimited") {
		gologger.Labelf("请稍后使用Robtex Api查詢 已限速 \n")
		return ""
	}
	var hostname []string
	var address []string
	result := "{\"passive_dns\":[" + string(resp.Body()) + "]}"
	responselist := gjson.Get(result, "passive_dns.#.rrdata").Array()
	for aa, _ := range responselist {
		ip := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)
		host := regexp.MustCompile(`(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`)
		if ip.FindAllStringSubmatch(strings.TrimSpace(responselist[aa].String()), -1) != nil {
			address = append(address, responselist[aa].String())
		} else if host.FindAllStringSubmatch(strings.TrimSpace(responselist[aa].String()), -1) != nil {
			hostname = append(hostname, responselist[aa].String())
		}

	}
	var add int
	result1 := "{\"passive_dns\":["
	if len(hostname) < len(address) {
		for add = 0; add < len(hostname); add++ {
			result1 += "{\"hostname\"" + ":" + "\"" + hostname[add] + "\"" + "," + "\"address\"" + ":" + "\"" + address[add] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, hostname[add])
			DomainsIP.IP = append(DomainsIP.IP, address[add])
		}
		for ii := add; ii < len(address); ii++ {
			result1 += "{\"address\"" + ":" + "\"" + address[ii] + "\"" + "},"
			DomainsIP.IP = append(DomainsIP.IP, address[ii])
		}

	} else {
		for add = 0; add < len(address); add++ {
			result1 += "{\"hostname\"" + ":" + "\"" + hostname[add] + "\"" + "," + "\"address\"" + ":" + "\"" + address[add] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, hostname[add])
			DomainsIP.IP = append(DomainsIP.IP, address[add])
		}
		for ii := add; ii < len(hostname); ii++ {
			result1 += "{\"hostname\"" + ":" + "\"" + hostname[ii] + "\"" + "},"
			DomainsIP.Domains = append(DomainsIP.Domains, hostname[ii])
		}
	}

	result1 = result1 + "]}"
	res, ensOutMap := GetEnInfo(result1, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Robtex Api", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
