package ZoomEye

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"fmt"
	"strings"

	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	responselist := gjson.Get(response, "list").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "ZoomEye"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))oomEye
	for _, item := range responselist {
		// 从当前条目获取域名
		responsdomain := item.Get("name").String()

		// 获取当前条目的所有 IP 地址
		ips := item.Get("ip").Array() // 假设每个条目下的 "ip" 是一个数组

		// 为了构建 JSON 字符串，我们先创建 IP 地址的字符串数组
		var ipStrs []string
		for _, ip := range ips {
			ipStrs = append(ipStrs, fmt.Sprintf("\"%s\"", ip.String()))
			DomainsIP.Domains = append(DomainsIP.Domains, ip.String())
		}
		// 将 IP 地址数组转换为一个字符串，以逗号分隔
		ipStr := strings.Join(ipStrs, ",")

		// 构建包含 hostname 和所有 IP 地址的 JSON 字符串
		responseJia := fmt.Sprintf("{\"hostname\": \"%s\", \"address\": [%s]}", responsdomain, ipStr)
		DomainsIP.Domains = append(DomainsIP.Domains, responsdomain)

		// 将构建的 JSON 字符串解析为 gjson.Result 并追加到 ensInfos.Infos["Urls"]
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(responseJia))
	}
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func ZoomEye(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	gologger.Infof("ZoomEye 威胁平台查询\n")

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"API-KEY":    {options.ENConfig.Cookies.Zoomeye},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()
	currentPage := 1
	result := "{\"list\": ["
	for {

		urls := fmt.Sprintf("https://api.zoomeye.org/domain/search?q=%s&type=1&s=1000&page=%d", domain, currentPage)
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
		if resp.Body() == nil || gjson.Get(string(resp.Body()), "total").Int() == 0 {
			gologger.Labelf("ZoomEye 威胁平台未发现域名\n")
			return ""
		}
		responselist := gjson.Get(string(resp.Body()), "list").Array()
		for _, item := range responselist {
			result = result + item.String()
		}
		currentPage++
		if len(responselist) == 0 || currentPage == 10 {
			result = result + "], }"
			break
		}

	}

	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "ZoomEye", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
