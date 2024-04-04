package Github

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
	ensInfos.SType = "Github"
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
	Utils.DomainTableShow(keyword, data, "Github")
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
func GithubEmail(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {

	var respnsehe string
	for add := 1; add < 11; add += 1 {
		urls := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100&page=%d&sort=indexed&&order=asc", domain, add)

		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
		if options.Proxy != "" {
			client.SetProxy(options.Proxy)
		}
		Authorization := " token " + options.ENConfig.Cookies.Github
		client.Header = http.Header{
			"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
			"Accept":          {"application/vnd.github.v3.text-match+json"},
			"Accept-Encoding": {"gzip"},
			"Authorization":   {Authorization},
		}

		client.Header.Del("Cookie")

		//强制延时1s
		time.Sleep(3 * time.Second)
		//加入随机延迟
		time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
		clientR := client.R()

		clientR.URL = urls
		resp, err := clientR.Get(urls)

		for attempt := 0; attempt < 4; attempt++ {
			if resp.RawResponse == nil {
				resp, _ = clientR.Get(urls)
				time.Sleep(3 * time.Second)
			} else if strings.Contains(string(resp.Body()), "API rate limit exceeded for") {
				time.Sleep(5 * time.Second)
			} else if resp.Body() != nil {
				break
			}
		}
		if resp.RawResponse == nil || err != nil {
			gologger.Errorf("Github 链接无法访问尝试切换代理 \n")
			break
		}
		if gjson.Get(string(resp.Body()), "total_count").Int() == 0 && add == 1 {
			//gologger.Labelf("github 未发现域名 %s\n", urls)
			break
		} else if gjson.Get(string(resp.Body()), "total_count").Int() == 0 && add != 1 {
			break
		} else if len(gjson.Get(string(resp.Body()), "items").Array()) == 0 {
			break
		}
		respnsehe += string(resp.Body())

	}
	respnsehe = clearresponse(respnsehe)
	Email := `[a-zA-Z0-9.\-_+#~!$&',;=:]+@` + `[a-zA-Z0-9.-]*` + strings.ReplaceAll(domain, "www.", "")
	re := regexp.MustCompile(Email)

	Emails := re.FindAllStringSubmatch(strings.TrimSpace(respnsehe), -1)
	if len(Emails) > 0 {
		result1 := "{\"Email\":["
		for add := 0; add < len(Emails); add++ {
			result1 += "{" + "\"Email\"" + ":" + "\"" + Emails[add][0] + "\"" + "}" + ","

		}
		result1 = result1 + "]}"

		res, ensOutMap := GetEnInfo(result1, DomainsIP)
		//
		outputfile.MergeOutPut(res, ensOutMap, "Github", options)
	}

	return "Success"
}
