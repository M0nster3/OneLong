package HttpZhiwen

import (
	"OneLong/IP"
	"OneLong/Script/Ehole"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/CDN"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func GetEnInfo(response string) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	respons := gjson.Parse(response).Array()

	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range getENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))

	for _, aa := range respons {
		ensInfos.Infos["Zhiwen"] = append(ensInfos.Infos["Zhiwen"], gjson.Parse(aa.String()))

	}
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}
func Status(domaina string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	DomainsIP.IP = Utils.SetStr(DomainsIP.IP)
	//Ip反差域名

	gologger.Infof("IP反差域名 \n")
	for _, ip := range DomainsIP.IP {
		if !CDN.CheckIP(ip) {
			wg.Add(1)
			go func() {
				IP.IpWhois(domaina, ip, options, DomainsIP)
				wg.Done()
			}()

		}
	}
	wg.Wait()
	DomainsIP.Domains = Utils.SetStr(DomainsIP.Domains)
	gologger.Infof("检测域名存活\n")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}

	for _, domain := range DomainsIP.Domains {
		urls := domain
		client.Header = http.Header{
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
			"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		}
		clientR := client.R()

		if !strings.Contains(urls, "http://") && !strings.Contains(urls, "https://") {
			urls = "http://" + urls
		}
		clientR.URL = urls
		response, err := clientR.Send()
		if err != nil {
			DomainsIP.Hostnameip = append(DomainsIP.Hostnameip, "")
			DomainsIP.Status_code = append(DomainsIP.Status_code, "")
			DomainsIP.TitleBUff = append(DomainsIP.TitleBUff, "")
			DomainsIP.Zhiwen = append(DomainsIP.Zhiwen, "")
			continue
		}
		add := 0
		for {
			if add == 2 {
				urls = strings.ReplaceAll(urls, "http://", "https://")
				clientR.URL = urls
				add += 1
				continue
			}
			if response.RawResponse == nil {
				response, _ = clientR.Send()
				time.Sleep(2 * time.Second)
				add += 1
			} else if response.Body() != nil {
				break
			} else if add == 4 {
				break
			}

		}
		if response.Body() != nil {
			DomainsIP.Status_code = append(DomainsIP.Status_code, strconv.Itoa(response.StatusCode()))
			// 解析HTML响应体
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(response.Body())))
			if err != nil {
				log.Fatalf("Failed to read document: %v", err)
			}

			// 查找<title>标签并获取其文本内容
			title := doc.Find("title").Text()
			fmt.Printf("Title: %s\n", title)
			DomainsIP.TitleBUff = append(DomainsIP.TitleBUff, title)
			ips, _ := net.LookupIP(domain)
			var rangip string
			for _, ip := range ips {
				if len(ips) > 1 {
					rangip = rangip + ip.String() + " , "
				} else {
					rangip = rangip + ip.String()
				}

			}
			DomainsIP.Hostnameip = append(DomainsIP.Hostnameip, rangip)
			Ehole.Ehole(urls, options, DomainsIP)

		} else {
			DomainsIP.Hostnameip = append(DomainsIP.Hostnameip, "")
			DomainsIP.Status_code = append(DomainsIP.Status_code, "")
			DomainsIP.TitleBUff = append(DomainsIP.TitleBUff, "")
			DomainsIP.Zhiwen = append(DomainsIP.Zhiwen, "")
		}

	}
	result := "["
	var add int
	if len(DomainsIP.Domains) < len(DomainsIP.IP) {
		for add = 0; add < len(DomainsIP.Domains); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + DomainsIP.Hostnameip[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"
		}
		for ii := add; ii < len(DomainsIP.IP); ii++ {
			result += "{\"address\"" + ":" + "\"" + DomainsIP.IP[ii] + "\"" + "},"
		}

	} else {
		for add = 0; add < len(DomainsIP.IP); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + DomainsIP.Hostnameip[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"

		}
		for ii := add; ii < len(DomainsIP.Domains); ii++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[ii] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + DomainsIP.Hostnameip[ii] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[ii] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[ii] + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[ii] + "\"" + "},"

		}
	}
	result = result + "]"
	res, ensOutMap := GetEnInfo(result)

	outputfile.MergeOutPut(res, ensOutMap, "Ehole", options)
}
