package baidu

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"

	"github.com/tidwall/gjson"
	"regexp"
	"strings"
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
	//命令输出展示

	var data [][]string
	var keyword []string
	for _, y := range getENMap() {
		for _, ss := range y.keyWord {
			if ss == "数据关联" {
				continue
			}
			keyword = append(keyword, ss)
		}

		for _, res := range ensInfos.Infos["Email"] {
			results := gjson.GetMany(res.Raw, y.field...)
			var str []string
			for _, s := range results {
				str = append(str, s.String())
			}
			data = append(data, str)
		}

	}

	Utils.DomainTableShow(keyword, data, "Baidu")
	return ensInfos, ensOutMap

}

func baiduparseUrl(domain string) []string {

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

func Baidu(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) {
	//gologger.Infof("Alienvault\n")

	urlss := baiduparseUrl(domain)
	var inputTargetHost goflags.StringSlice
	for _, url := range urlss {
		inputTargetHost = append(inputTargetHost, url)
	}

	// 更新 httpxoptions 中的 InputTargetHost

	var respnsehe string
	httpxoptions := runner.Options{
		Methods:         "GET",
		InputTargetHost: inputTargetHost,
		RateLimit:       1,
		Threads:         3,
		HTTPProxy:       options.Proxy,
		//InputFile: "./targetDomains.txt", // path to file containing the target domains list
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)

			}
			//fmt.Printf("%s %s %d\n", r.Input, r.Host, r.StatusCode)
			respnsehe += r.Raw
		},
	}

	if err := httpxoptions.ValidateOptions(); err != nil {
		//log.Fatal(err)
	}

	httpxRunner, err := runner.New(&httpxoptions)
	if err != nil {
		//log.Fatal(err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()
	//for _, urls := range urlss {

	//client := resty.New()
	//// 获取系统根证书列表
	////roots, err := x509.SystemCertPool()
	////if err != nil {
	////	// 处理错误
	////}
	//
	//client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	//if options.Proxy != "" {
	//	client.SetProxy(options.Proxy)
	//}
	//client.Header = http.Header{
	//	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	//	"Host":            {"www.baidu.com"},
	//	"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
	//	"Accept-Encoding": {"gzip"},
	//}
	//
	////强制延时1s
	//time.Sleep(3 * time.Second)
	////加入随机延迟
	//time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	//clientR := client.R()
	//
	//clientR.URL = urls
	//resp, err := clientR.Get(urls)
	//
	//for add := 1; add < 4; add += 1 {
	//	if resp.RawResponse == nil {
	//		resp, _ = clientR.Get(urls)
	//		time.Sleep(3 * time.Second)
	//	} else if resp.Body() != nil {
	//		break
	//	} else if strings.Contains(string(resp.Body()), "网络不给力，请稍后重试") {
	//		break
	//	}
	//}
	//if err != nil || strings.Contains(string(resp.Body()), "网络不给力，请稍后重试") {
	//	gologger.Errorf("Baidu 爬虫失败尝试切换代理\n")
	//	break
	//}
	//respnsehe += string(resp.Body())

	//}
	//
	respnsehe = clearresponse(respnsehe)
	//if strings.Contains(respnsehe, "title 百度安全验证  title ") {
	//	gologger.Errorf("Baidu需要进行验证\n")
	//	return
	//}
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
