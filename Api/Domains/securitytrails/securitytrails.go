package securitytrails

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Securitytrails"
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

func Securitytrails(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Securitytrails 搜索域名 \n")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	urls := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", domain)

	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"APIKEY":     {options.ENConfig.Cookies.Securitytrails},
	}
	client.Header.Set("Content-Type", "application/json")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()
	response, err := clientR.Get(urls)
	for add := 1; add < 4; add += 1 {
		if response.RawResponse == nil {
			response, _ = clientR.Get(urls)
			time.Sleep(1 * time.Second)
		} else if response.Body() != nil {
			break
		}
	}
	if err != nil {
		gologger.Errorf("Securitytrails API 链接访问失败尝试切换代理\n")
		return ""
	}
	if len(gjson.Get(string(response.Body()), "subdomains").Array()) == 0 {
		gologger.Labelf("Securitytrails 未发现域名 %s\n", domain)
		return ""
	}
	Hostname := gjson.Get(string(response.Body()), "subdomains").Array()
	// 查找具有特定 class 的元素并获取其内容
	//var Hostname []string

	var result string
	result = "{\"passive_dns\":["
	for i := 0; i < len(Hostname); i++ {
		result += "{\"hostname\"" + ":" + "\"" + Hostname[i].String() + "." + domain + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, Hostname[i].String()+"."+domain)
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Securitytrails", options)

	return "Success"
}
