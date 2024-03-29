package Racent

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"fmt"
	"github.com/gookit/color"
	"sync"

	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

var mu sync.Mutex // 用于保护 addedURLs
func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Get(response, "data.list.#.dnsnames").Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Racent"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
	addedURLs := make(map[string]bool)
	for bb, _ := range respons {
		for _, aa := range respons[bb].Array() {
			ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + aa.String() + "\"" + "}"
			url := gjson.Parse(ResponseJia).Get("hostname").String()
			DomainsIP.Domains = append(DomainsIP.Domains, url)
			// 检查是否已存在相同的 URL
			if !addedURLs[url] {
				// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
				ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
				addedURLs[url] = true
			}
		}

	}
	mu.Lock()
	//命令输出展示
	color.RGBStyleFromString("205,155,29")
	color.RGBStyleFromString("205,155,29").Println("\nRacent 查询子域名")
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

func Racent(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Fullhunt 威胁平台查询\n")
	urls := fmt.Sprintf("https://face.racent.com/tool/query_ctlog?token=%s&keyWord=%s", options.ENConfig.Cookies.Racent, domain)

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
		gologger.Errorf("Racent API 链接访问失败尝试切换代理\n")
		return ""
	}
	if gjson.Get(string(resp.Body()), "data.count").Int() == 0 {
		gologger.Labelf("Racent 未发现域名 %s\n", domain)
		return ""
	}
	res, ensOutMap := GetEnInfo(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Fullhunt", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
