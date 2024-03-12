package HttpZhiwen

import (
	"OneLong/IP"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/CDN"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Status(domaina string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {

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
	DomainsIP.IP = Utils.SetStr(DomainsIP.IP)
	DomainsIP.Domains = Utils.SetStr(DomainsIP.Domains)
	gologger.Infof("检测域名存活\n")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}

	var TitleBUff []string
	var hostnameip []string
	var status_code []string
	for _, domain := range DomainsIP.Domains {
		urls := domain
		client.Header = http.Header{
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
			"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		}
		clientR := client.R()
		clientR.URL = urls
		response, _ := clientR.Send()
		for {
			if response.RawResponse == nil {
				response, _ = clientR.Send()
				time.Sleep(2 * time.Second)
			} else if response.Body() != nil {
				break
			}
		}
		if response.StatusCode() != 502 {
			status_code = append(status_code, string(response.StatusCode()))
			// 解析HTML响应体
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(response.Body())))
			if err != nil {
				log.Fatalf("Failed to read document: %v", err)
			}

			// 查找<title>标签并获取其文本内容
			title := doc.Find("title").Text()
			fmt.Printf("Title: %s\n", title)
			TitleBUff = append(TitleBUff, title)
			ips, _ := net.LookupIP(domain)
			for _, ip := range ips {
				hostnameip = append(hostnameip, ip.String())
			}

		} else {
			hostnameip = append(hostnameip, "")
			status_code = append(status_code, "")
			TitleBUff = append(TitleBUff, "")
		}
	}
	result := "{["
	var add int
	if len(DomainsIP.Domains) < len(DomainsIP.IP) {
		for add = 0; add < len(DomainsIP.Domains); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + hostnameip[add] + "," + "\"status_code\"" + ":" + "\"" + status_code[add] + "," + "\"title\"" + ":" + "\"" + TitleBUff[add] + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "},"
		}
		for ii := add; ii < len(DomainsIP.IP); ii++ {
			result += "{\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + hostnameip[add] + "," + "\"status_code\"" + ":" + "\"" + status_code[add] + "," + "\"title\"" + ":" + "\"" + TitleBUff[add] + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "},"
		}

	} else {
		for add = 0; add < len(DomainsIP.IP); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + hostnameip[add] + "," + "\"status_code\"" + ":" + "\"" + status_code[add] + "," + "\"title\"" + ":" + "\"" + TitleBUff[add] + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "},"

		}
		for ii := add; ii < len(DomainsIP.Domains); ii++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.Domains[add] + "\"" + "," + "\"hostnameip\"" + ":" + "\"" + hostnameip[add] + "," + "\"status_code\"" + ":" + "\"" + status_code[add] + "," + "\"title\"" + ":" + "\"" + TitleBUff[add] + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "},"

		}
	}

}
