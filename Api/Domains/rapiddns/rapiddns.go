package rapiddns

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
	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Rapiddns"
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

func Rapiddns(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Rapiddns Api搜索域名 \n")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	urls := "https://rapiddns.io/subdomain/" + domain + "?full=1"
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	}

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()
	clientR.URL = urls
	response, _ := clientR.Send()
	for {
		if response.RawResponse == nil {
			response, _ = clientR.Send()
			time.Sleep(1 * time.Second)
		} else if response.Body() != nil {
			break
		}
	}
	//Total: <span style="color: #39cfca; ">0
	if strings.Contains(string(response.Body()), "Total: <span style=\"color: #39cfca; \">0") {
		gologger.Labelf("Rapiddns Api未发现域名 %s\n", domain)
		return ""
	}
	host := regexp.MustCompile(`<td>((?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,})</td>`)
	hostrea := host.FindAllStringSubmatch(strings.TrimSpace(string(response.Body())), -1)
	var Hostname []string
	// 使用 HTML 解析器解析 HTML 内容
	for _, bu := range hostrea {
		Hostname = append(Hostname, bu[1])
	}
	Hostname = Utils.SetStr(Hostname)
	// 查找具有特定 class 的元素并获取其内容
	//var Hostname []string

	var result string
	result = "{\"passive_dns\":["
	for i := 0; i < len(Hostname); i++ {
		result += "{\"hostname\"" + ":" + "\"" + Hostname[i] + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, Hostname[i])
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Rapiddns", options)

	return "Success"
}
