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

func IpWhois(domain string, ip string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) {
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
	time.Sleep(3 * time.Second)

	clientR := client.R()

	clientR.URL = urls
	resp, _ := clientR.Get(urls)
	for i := 0; i < 4; i++ {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(3 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	Domains := gjson.Parse(string(resp.Body())).Array()
	//if len(Domains) == 0 {
	//	return
	//}
	var add int
	var resdoaminc string
	//if ip == "208.113.222.142" {
	//	print(111)
	//}
	for _, appdomain := range Domains {
		if len(Domains) == 0 {
			break
		}
		resdoamin := gjson.Get(appdomain.String(), "domain").String()
		resdoaminb := strings.Split(resdoamin, ".")
		if len(resdoaminb) > 2 {
			resdoaminc = resdoaminb[len(resdoaminb)-2] + "." + resdoaminb[len(resdoaminb)-1]
		} else {
			resdoaminc = resdoamin
		}
		if strings.Contains(domain, resdoaminc) {
			DomainsIP.IPA = append(DomainsIP.IPA, ip)
			DomainsIP.Domains = append(DomainsIP.Domains, resdoamin)

		} else if add == 15 {
			break
		} else {
			add += 1
		}

	}

}
