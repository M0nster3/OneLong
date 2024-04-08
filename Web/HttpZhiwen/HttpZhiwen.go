package HttpZhiwen

import (
	"OneLong/IP"
	"OneLong/Script/Ehole"
	"OneLong/Script/Port"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/CDN"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"net"
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

	for _, aa := range respons {
		ensInfos.Infos["Zhiwen"] = append(ensInfos.Infos["Zhiwen"], gjson.Parse(aa.String()))

	}

	return ensInfos, ensOutMap

}
func Status(domaina string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) {
	DomainsIP.IP = Utils.SetStr(DomainsIP.IP)
	//Ip反差域名

	//gologger.Infof("IP反差域名 \n")
	for _, ip := range DomainsIP.IP {
		wg.Add(1)
		ip := ip
		go func() {
			if !CDN.CheckIP(ip) {
				IP.IpWhois(domaina, ip, options, DomainsIP)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	DomainsIP.IPA = Utils.SetStr(DomainsIP.IPA)
	DomainsIP.Domains = Utils.SetStr(DomainsIP.Domains)
	color.RGBStyleFromString("244,211,49").Println("\n--------------------爆破C段端口--------------------")

	for _, C := range DomainsIP.IPA {
		ip := net.ParseIP(C)
		if ip == nil {
			continue
		}
		cidr := fmt.Sprint("%s/24", ip.String())
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		DomainsIP.C = append(DomainsIP.C, ipnet.String())
	}

	var Config Port.Config
	for _, C := range DomainsIP.C {
		Config.Target = C
		Config.Rate = options.LongConfig.Port.Masscan.Rate
		Config.Port = options.LongConfig.Port.Masscan.Port
		Port.DoMasscanPlusNmap(Config, options)
	}

	color.RGBStyleFromString("244,211,49").Println("\n--------------------检测指纹以及域名存活--------------------")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}

	Ehole.Ehole(DomainsIP.Domains, options, DomainsIP)

	result := "["
	var add int
	if len(DomainsIP.DomainA) < len(DomainsIP.IPA) {
		for add = 0; add < len(DomainsIP.DomainA); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IPA[add] + "\"" + "," + "\"A\"" + ":" + "\"" + DomainsIP.A[add] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"
		}
		for ii := add; ii < len(DomainsIP.IPA); ii++ {
			result += "{\"address\"" + ":" + "\"" + DomainsIP.IPA[ii] + "\"" + "},"
		}

	} else {
		for add = 0; add < len(DomainsIP.IPA); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IPA[add] + "\"" + "," + "\"A\"" + ":" + "\"" + DomainsIP.A[add] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"

		}

		for ii := add; ii < len(DomainsIP.DomainA) && ii < len(DomainsIP.A) && ii < len(DomainsIP.Size) && ii < len(DomainsIP.Status_code) && ii < len(DomainsIP.TitleBUff) && ii < len(DomainsIP.Zhiwen); ii++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[ii] + "\"" + "," + "\"A\"" + ":" + "\"" + DomainsIP.A[add] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[ii] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[ii] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[ii] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[ii] + "\"" + "},"
		}
	}
	result = result + "]"
	res, ensOutMap := GetEnInfo(result)

	outputfile.MergeOutPut(res, ensOutMap, "Ehole", options)
}
