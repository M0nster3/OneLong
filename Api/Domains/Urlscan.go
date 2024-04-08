package Domains

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

// 用于保护 addedURLs
func GetEnInfoUrlscan(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "urlscan"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	for aa, _ := range respons {
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(respons[aa].String()))
	}

	//命令输出展示

	var data [][]string
	var keyword []string
	for _, y := range GetENMap() {
		for _, ss := range y.KeyWord {
			if ss == "数据关联" {
				continue
			}
			keyword = append(keyword, ss)
		}

		for _, res := range ensInfos.Infos["Urls"] {
			results := gjson.GetMany(res.Raw, y.Field...)
			var str []string
			for _, s := range results {
				str = append(str, s.String())
			}
			data = append(data, str)
		}

	}

	Utils.DomainTableShow(keyword, data, "Urlscan")

	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func Urlscan(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {

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
	time.Sleep(3 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
	for add := 1; add < 4; add += 1 {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(3 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if err != nil {
		gologger.Errorf("Urlscan链接访问失败尝试切换代理\n")
		return ""
	}
	if len(gjson.Get(string(resp.Body()), "results").Array()) == 0 {
		//gologger.Labelf("Urlscan未发现域名 %s\n", domain)
		return ""
	}
	hostname := gjson.Get(string(resp.Body()), "results.#.task.domain").Array()
	address := gjson.Get(string(resp.Body()), "results.#.page.ip").Array()
	//HeiUrl := "https://www.lianwoapp.com，https://www.momojc.cn，https://www.omnex.com.cn，https://omnex.com.cn，https://www.admach.net，https://i.zkaq.cn，http://csfdk.com，https://jnbeiyou.com，https://www.wwhxkj.com，google.comhttp://www.hellolunarly.com/,https://doc.nasdt.cn,http://fyi6.uuidapi.anz.online.cfgsg.comhttp://vhn4.uuidapi.anz.online.cfgsg.comhttp://socialmediamanagementuk.comhttp://www.billurcevre.comhttps://www.artisanbenefitauctioneer.comhttps://jiance.icloudshield.comhttp://kekeyue.comhttps://t.cohttp://www.shicaituku.comhttps://www.bdjiayu.comhttp://telegramlh.com,chqzyy.com"
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

	res, ensOutMap := GetEnInfoUrlscan(result1, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Urlscan", options)

	return "Success"
}
