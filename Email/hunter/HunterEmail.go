package hunter

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
	respons := gjson.Get(response, "Email").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Baidu"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
	addedURLs := make(map[string]bool)
	for aa, _ := range respons {
		hostname := gjson.Get(respons[aa].String(), "Email").String()
		if !addedURLs[hostname] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Email"] = append(ensInfos.Infos["Email"], gjson.Parse(respons[aa].String()))
			addedURLs[hostname] = true
		}

	}

	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.Getfield, ensInfos, options)
	return ensInfos, ensOutMap

}

func clearresponse(results string) string {

	replacements := []string{
		"<em>", "</em>", // 替换 <em> 和 </em>
		"<b>", "</b>", // 替换 <b> 和 </b>
		"%3a",                   // 替换 %3a
		"<strong>", "</strong>", // 替换 <strong> 和 </strong>
		"<wbr>", "</wbr>", // 替换 <wbr> 和 </wbr>
	}
	replacements2 := []string{
		"<", ">", ":", "=", ";", "&", "%3A", "%3D", "%3C", "%2f", "/", "\\", // 其他需要替换的字符
	}

	// 执行替换
	for _, search := range replacements {
		results = strings.ReplaceAll(results, search, "")
	}
	for _, search := range replacements2 {
		results = strings.ReplaceAll(results, search, " ")
	}
	return results

}
func HunterEmail(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {

	var respnsehe string

	urls := fmt.Sprintf("https://api.hunter.io/v2/domain-search?domain=%s&api_key=%s&limit=10", domain, options.ENConfig.Email.Emailhunter)

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"application/vnd.github.v3.text-match+json"},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)

	for attempt := 0; attempt < 4; attempt++ {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if resp.RawResponse == nil || err != nil {
		gologger.Errorf("Github 链接无法访问尝试切换代理 \n")
		return ""
	}
	//if gjson.Get(string(resp.Body()), "total_count").Int() == 0 {
	//	gologger.Labelf("github 未发现域名 %s\n", domain)
	//	return ""
	//} else if len(gjson.Get(string(resp.Body()), "items").Array()) == 0 {
	//	return ""
	//}
	respnsehe = clearresponse(string(resp.Body()))
	Email := `[a-zA-Z0-9.\-_+#~!$&',;=:]+@` + `[a-zA-Z0-9.-]*` + strings.ReplaceAll(domain, "www.", "")
	re := regexp.MustCompile(Email)

	Emails := re.FindAllStringSubmatch(strings.TrimSpace(respnsehe), -1)

	result1 := "{\"Email\":["
	for add := 0; add < len(Emails); add++ {
		result1 += "{" + "\"Email\"" + ":" + "\"" + Emails[add][0] + "\"" + "}" + ","

	}
	result1 = result1 + "]}"

	//for _, aa := range matches {
	//	fmt.Print("111111\n")
	//	fmt.Print(aa)
	//}
	res, ensOutMap := GetEnInfo(result1, DomainsIP)
	//
	outputfile.MergeOutPut(res, ensOutMap, "Hunter", options)
	return "Success"
}
