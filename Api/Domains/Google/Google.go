package Google

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	color.RGBStyleFromString("199,21,133").Println("\nGoogle.go 查询子域名")
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
func countCharacters(arr []string) map[string]int {
	charCount := make(map[string]int)

	for _, char := range arr {
		charCount[char]++
	}

	return charCount
}
func Google(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Google 搜索域名 \n")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	urls := "https://www.googleapis.com/customsearch/v1"

	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	}
	client.Header.Set("Content-Type", "application/json")

	//强制延时1s

	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)

	var buff []string
	for start := 1; start < 120; start += 10 {
		queryParams := map[string]string{
			"key":    options.ENConfig.Cookies.GoogleApi,
			"cx":     options.ENConfig.Cookies.GoogleID,
			"q":      "site:." + domain,
			"fields": "items/link",
			"start":  strconv.Itoa(start),
			"num":    "10",
		}
		clientR := client.R()
		for key, value := range queryParams {
			clientR = clientR.SetQueryParam(key, value)
		}
		response, err := clientR.Get(urls)
		for attempt := 0; attempt < 4; attempt++ {
			if response.RawResponse == nil {
				response, _ = clientR.Get(urls)
				time.Sleep(1 * time.Second)
			} else if response.Body() != nil {
				break
			}
		}
		if err != nil {
			gologger.Errorf("Github 链接无法访问尝试切换代理 \n")
			return ""
		}
		time.Sleep(1 * time.Second)
		llink := gjson.Get(string(response.Body()), "items.#.link").Array()
		hostname := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
		// 编译正则表达式
		re := regexp.MustCompile(hostname)
		for _, aa := range llink {
			matches := re.FindAllStringSubmatch(strings.TrimSpace(aa.String()), -1)
			buff = append(buff, matches[0][0])

		}
		// 查找匹配的内容

		if response.Size() == 3 {
			break
		} else if start > 100 {
			resultaa := countCharacters(buff)
			statement := " "
			for char, count := range resultaa {
				if count > 1 {
					statement = statement + "-site:" + char + " "
				}
			}
			for startB := 1; startB < 120; startB += 10 {
				queryParams := map[string]string{
					"key":    options.ENConfig.Cookies.GoogleApi,
					"cx":     options.ENConfig.Cookies.GoogleID,
					"q":      "site:." + domain + statement,
					"fields": "items/link",
					"start":  strconv.Itoa(startB),
					"num":    "10",
				}
				clientR := client.R()
				for key, value := range queryParams {
					clientR = clientR.SetQueryParam(key, value)
				}
				response, _ = clientR.Get(urls)
				time.Sleep(1 * time.Second)
				startB += 10
				llink := gjson.Get(string(response.Body()), "items.#.link").Array()
				// 编译正则表达式
				re := regexp.MustCompile(hostname)
				for _, aa := range llink {
					matches := re.FindAllStringSubmatch(strings.TrimSpace(aa.String()), -1)
					buff = append(buff, matches[0][0])

				}
				// 查找匹配的内容

				if startB > 100 || response.Size() == 3 {
					break
				}
			}
			break
		}
	}

	// 查找具有特定 class 的元素并获取其内容
	//var Hostname []string
	buff = Utils.SetStr(buff)
	var result string
	result = "{\"passive_dns\":["
	for i := 0; i < len(buff); i++ {
		result += "{\"hostname\"" + ":" + "\"" + buff[i] + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, buff[i])
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Google", options)

	return "Success"
}
