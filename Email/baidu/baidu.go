package baidu

import (
	Email2 "OneLong/Email/yahoo"
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

	addedURLs := make(map[string]bool)
	for aa, _ := range respons {
		hostname := gjson.Get(respons[aa].String(), "Email").String()
		if !addedURLs[hostname] {
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Email"] = append(ensInfos.Infos["Email"], gjson.Parse(respons[aa].String()))
			addedURLs[hostname] = true
		}

	}

	return ensInfos, ensOutMap

}

func parseUrl(domain string) []string {

	var urls []string
	for num := 0; num < 500; num += 10 {
		url := "https://www.baidu.com/s?wd=%40" + domain + "&pn=xx&oq=" + domain
		url = strings.ReplaceAll(url, "xx", fmt.Sprintf("%d", num))
		urls = append(urls, url)
	}
	return urls

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

func Baidu(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	//gologger.Infof("Alienvault\n")

	urlss := Email2.ParseUrl(domain)
	var respnsehe string
	for _, urls := range urlss {

		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
		if options.Proxy != "" {
			client.SetProxy(options.Proxy)
		}
		client.Header = http.Header{
			"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"Accept-Encoding": {"gzip, deflate, br, zstd"},
			"Accept-Language": {"zh-CN,zh;q=0.9"},
			"Cache-Control":   {"max-age=0"},
			"Connection":      {"keep-alive"},
			//"Cookie":                    {"BIDUPSID=E5CBE91A1AFA0DE768B8836161C4B5DD; PSTM=1711609247; BAIDUID=E5CBE91A1AFA0DE7D8785DAA4084CF17:FG=1; BD_UPN=12314753; H_PS_PSSID=40170_40080_40368_40379_40415_40445_40464_40458_40481_40317_39661_40487_40511_40514_40398_60041_60027_60034_60047; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; kleck=86bf3a1f689c07d154ee2ce3f9366ebc; BA_HECTOR=0l2l818lak2k24a48g2k858hs14dms1j0a7je1t; BAIDUID_BFESS=E5CBE91A1AFA0DE7D8785DAA4084CF17:FG=1; ZFY=sF0pysXgRmUHpErRNjf9aHTIUiPPBAQqlz:BaQbYC8Oc:C; delPer=0; BD_CK_SAM=1; PSINO=7; H_PS_645EC=671doRul3CR5KeWAZpKiD%2FvWfSm12fpyFlBeTAQCileoBXECcVJwChS9tKs"},
			"Host":                      {"www.baidu.com"},
			"Sec-Ch-Ua":                 {"\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Google Chrome\";v=\"122\""},
			"Sec-Ch-Ua-Mobile":          {"?0"},
			"Sec-Ch-Ua-Platform":        {"\"Windows\""},
			"Sec-Fetch-Dest":            {"document"},
			"Sec-Fetch-Mode":            {"navigate"},
			"Sec-Fetch-Site":            {"none"},
			"Sec-Fetch-User":            {"?1"},
			"Upgrade-Insecure-Requests": {"1"},
			"User-Agent":                {Utils.RandUA()},
		}

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
			} else if strings.Contains(string(resp.Body()), "网络不给力，请稍后重试") {
				break
			}
		}
		if err != nil || strings.Contains(string(resp.Body()), "网络不给力，请稍后重试") {
			gologger.Errorf("Baidu 爬虫失败尝试切换代理\n")
			break
		}
		respnsehe += string(resp.Body())

	}

	respnsehe = clearresponse(respnsehe)
	if strings.Contains(respnsehe, "title 百度安全验证  title ") {
		gologger.Errorf("Baidu需要进行验证\n")
		return
	}
	Email := `[a-zA-Z0-9.\-_+#~!$&',;=:]+@` + `[a-zA-Z0-9.-]*` + strings.ReplaceAll(domain, "www.", "")

	re := regexp.MustCompile(Email)

	Emails := re.FindAllStringSubmatch(strings.TrimSpace(respnsehe), -1)

	result1 := "{\"Email\":["
	for add := 0; add < len(Emails); add++ {
		result1 += "{" + "\"Email\"" + ":" + "\"" + Emails[add][0] + "\"" + "}" + ","

	}
	result1 = result1 + "]}"

	res, ensOutMap := GetEnInfo(result1, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Baidu", options)

}
