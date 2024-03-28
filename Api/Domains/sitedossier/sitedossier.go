package sitedossier

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	responselist := gjson.Parse(response).Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Sitedossier"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))oomEye
	addedURLs := make(map[string]bool)
	for aa, _ := range responselist {
		ResponseJia := "{" + "\"hostname\"" + ":" + "\"" + responselist[aa].String() + "\"" + "}"
		urls := gjson.Parse(ResponseJia).Get("hostname").String()

		// 检查是否已存在相同的 URL
		if !addedURLs[urls] {
			DomainsIP.Domains = append(DomainsIP.Domains, urls)
			// 如果不存在重复则将 URL 添加到 Infos["Urls"] 中，并在 map 中标记为已添加
			ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(ResponseJia))
			addedURLs[urls] = true
		}

	}
	mu.Lock()
	//命令输出展示
	color.RGBStyleFromString("199,21,133").Println("\nsitedossier 查询子域名")
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

func Sitedossier(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Sitedossier Api查询\n")

	urls := fmt.Sprintf("http://www.sitedossier.com/parentdomain/%s", domain)
	UrlsB := urls
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

	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()
	result := "["
	clientR.URL = urls
	//强制延时1s
	time.Sleep(1 * time.Second)
	resp, err := clientR.Get(urls)
	for adda := 1; adda < 4; adda += 1 {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if err != nil {
		gologger.Errorf("Sitedossier API 链接访问失败尝试切换代理\n")
		return ""
	}
	if strings.Contains(string(resp.Body()), "No data currently available") {
		gologger.Labelf("Sitedossier Api 未发现域名 %s\n", domain)
		return ""
	} else if strings.Contains(string(resp.Body()), "only; case does not matter") {
		gologger.Labelf("如果想查询更多域名进入 www.sitedossier.com 输入验证码\n")
		return ""
	}
	// 定义匹配数字的正则表达式
	re := regexp.MustCompile(`out of a total of (\d{1,3}(,\d{3})*(\.\d+)?)`)

	// 使用正则表达式查找匹配的内容
	matches := re.FindStringSubmatch(string(resp.Body()))

	// 如果找到了匹配项
	// 提取匹配到的数字，并去除逗号
	totalItemsStr := strings.ReplaceAll(matches[1], ",", "")
	// 将字符串转换为整数
	totalItems, _ := strconv.Atoi(totalItemsStr)

	for add := 0; add < totalItems; add += 100 {

		if strings.Contains(string(resp.Body()), "Show next") || strings.Contains(string(resp.Body()), "Show remaining") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
			if err != nil {
				panic(err)
			}

			// 使用Find方法选择所有的<li>标签
			doc.Find("li").Each(func(i int, s *goquery.Selection) {
				// 在每个<li>标签中，使用Find方法选择<a>标签，并获取其href属性值
				href, exists := s.Find("a").Attr("href")
				if exists {
					// 打印href属性值
					result = result + "\"" + strings.ReplaceAll(href, "/site/", "") + "\","

				}
			})

			urls = UrlsB + "/" + strconv.Itoa(add)

		} else if strings.Contains(string(resp.Body()), "No data currently available.") {
			return ""
		} else if strings.Contains(string(resp.Body()), "End of list") {
			break
		} else if strings.Contains(string(resp.Body()), "you may see this page again. Thank you") {
			gologger.Labelf("如果想查询更多域名进入 www.sitedossier.com 输入验证码\n")
			return ""
		}

	}
	result = result + "]"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Sitedossier", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
