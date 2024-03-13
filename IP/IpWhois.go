package IP

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"time"
)

func IpWhois(domain string, ip string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	urls := fmt.Sprintf("http://api.webscan.cc/?action=query&ip=%s", ip)
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

	client.Header.Set("Content-Type", "application/json")
	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, _ := clientR.Send()
	for {
		if resp.RawResponse == nil {
			resp, _ = clientR.Send()
			time.Sleep(2 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	Domains := gjson.Parse(string(resp.Body())).Array()
	var add int
	for _, appdomain := range Domains {
		resdoamin := gjson.Get(appdomain.String(), "domain").String()
		if strings.Contains(resdoamin, domain) {
			DomainsIP.Domains = append(DomainsIP.Domains, resdoamin)

		} else {
			add += 1
		}
		if add == 5 {
			return
		}

	}

}
