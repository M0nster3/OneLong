package HttpZhiwen

import (
	"OneLong/IP"
	"OneLong/Script/Ehole"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/CDN"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
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
	Ehole.Ehole(DomainsIP.Domains, options, DomainsIP)

	result := "["
	var add int
	if len(DomainsIP.DomainA) < len(DomainsIP.IP) {
		for add = 0; add < len(DomainsIP.DomainA); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"
		}
		for ii := add; ii < len(DomainsIP.IP); ii++ {
			result += "{\"address\"" + ":" + "\"" + DomainsIP.IP[ii] + "\"" + "},"
		}

	} else {
		for add = 0; add < len(DomainsIP.IP); add++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[add] + "\"" + "," + "\"address\"" + ":" + "\"" + DomainsIP.IP[add] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[add] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[add] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[add] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[add] + "\"" + "},"

		}
		for ii := add; ii < len(DomainsIP.DomainA); ii++ {
			result += "{\"hostname\"" + ":" + "\"" + DomainsIP.DomainA[ii] + "\"" + "," + "\"Size\"" + ":" + "\"" + DomainsIP.Size[ii] + "\"" + "," + "\"status_code\"" + ":" + "\"" + DomainsIP.Status_code[ii] + "\"" + "," + "\"title\"" + ":" + "\"" + DomainsIP.TitleBUff[ii] + "\"" + "," + "\"Zhiwen\"" + ":" + "\"" + DomainsIP.Zhiwen[ii] + "\"" + "},"
		}
	}
	result = result + "]"
	res, ensOutMap := GetEnInfo(result)

	outputfile.MergeOutPut(res, ensOutMap, "Ehole", options)
}
