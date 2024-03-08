package hunter

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"encoding/base64"
	"fmt"

	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	responsip := gjson.Get(response, "data.arr.#.ip").Array()
	responsdomain := gjson.Get(response, "data.arr.#.domain").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Hunter"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
	addedURLs := make(map[string]bool)
	for aa, _ := range responsdomain {

		ResponseJia := "{\"address\"" + ":" + "\"" + responsip[aa].String() + "\"" + "," + "\"hostname\"" + ":" + "\"" + responsdomain[aa].String() + "\"" + "},"
		url := gjson.Parse(ResponseJia).Get("hostname").String()
		if !addedURLs[url] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
			addedURLs[url] = true
		}

	}
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func Hunter(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	gologger.Infof("Hunter 威胁平台查询\n")
	base64a := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("domain=\"%s\"", domain)))
	urls := "https://hunter.qianxin.com/openApi/search?api-key=" + options.ENConfig.Cookies.Hunter + "&search=" + base64a + "&page=1&page_size=100&is_web=3"
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
	resp, _ := clientR.Send()
	if gjson.Get(string(resp.Body()), "total").Int() == 0 {
		gologger.Labelf("Hunter 威胁平台未发现域名\n")
		return ""
	}
	res, ensOutMap := GetEnInfo(string(resp.Body()), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Hunter", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
