package Urlscan

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"net/http"
	"sync"

	//"strconv"
	//"strings"
	"time"
)

var mu sync.Mutex // 用于保护 addedURLs
func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Google"
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
	mu.Lock()
	//命令输出展示
	color.RGBStyleFromString("199,21,133").Println("\nUrlscan 查询子域名")
	var data [][]string
	var keyword []string
	for _, y := range getENMap() {
		for _, ss := range y.keyWord {
			if ss == "数据关联" {
				continue
			}
			keyword = append(keyword, ss)
		}

		for _, res := range ensInfos.Infos["Urls"] {
			results := gjson.GetMany(res.Raw, y.field...)
			var str []string
			for _, s := range results {
				str = append(str, s.String())
			}
			data = append(data, str)
		}

	}

	Utils.TableShow(keyword, data)
	mu.Unlock()
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func Urlscan(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {

	urls := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)
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

	client.Header.Del("Cookie")

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
		gologger.Errorf("Urlscan链接访问失败尝试切换代理\n")
		return ""
	}
	if len(gjson.Get(string(resp.Body()), "results").Array()) == 0 {
		gologger.Labelf("Urlscan未发现域名 %s\n", domain)
		return ""
	}
	hostname := gjson.Get(string(resp.Body()), "results.#.task.domain").Array()
	address := gjson.Get(string(resp.Body()), "results.#.page.ip").Array()
	// 查找匹配的内容
	var add int
	result1 := "{\"passive_dns\":["
	for add = 0; add < len(address); add++ {
		result1 += "{\"hostname\"" + ":" + "\"" + hostname[add].String() + "\"" + "," + "\"address\"" + ":" + "\"" + address[add].String() + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, hostname[add].String())
		DomainsIP.IP = append(DomainsIP.IP, address[add].String())
	}
	for ii := add; ii < len(hostname); ii++ {
		result1 += "{\"hostname\"" + ":" + "\"" + hostname[ii].String() + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, hostname[ii].String())
	}
	result1 = result1 + "]}"

	res, ensOutMap := GetEnInfo(result1, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Urlscan", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
