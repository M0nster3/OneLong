package Domains

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"github.com/projectdiscovery/chaos-client/pkg/chaos"
	"github.com/tidwall/gjson"
	"strings"
)

// 用于保护 addedURLs
func GetEnInfoChaos(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Chaos"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}

	for aa, _ := range respons {
		ensInfos.Infos["Urls"] = append(ensInfos.Infos["Urls"], gjson.Parse(respons[aa].String()))
	}

	//命令输出展示

	var data [][]string
	var keyword []string
	for _, y := range GetENMap() {
		for _, ss := range y.KeyWord {
			if ss == "数据关联" {
				continue
			}
			keyword = append(keyword, ss)
		}

		for _, res := range ensInfos.Infos["Urls"] {
			results := gjson.GetMany(res.Raw, y.Field...)
			var str []string
			for _, s := range results {
				str = append(str, s.String())
			}
			data = append(data, str)
		}

	}

	Utils.DomainTableShow(keyword, data, "chaos")

	return ensInfos, ensOutMap

}

func Chaos(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Chaos API 域名查询 \n")
	var Hostname []string
	chaosClient := chaos.New(options.LongConfig.Cookies.Chaos)
	for result := range chaosClient.GetSubdomains(&chaos.SubdomainsRequest{
		Domain: domain,
	}) {
		if result.Error != nil {
			break
		}
		res := strings.ReplaceAll(result.Subdomain, "*.", "")
		Hostname = append(Hostname, res+"."+domain)

	}
	if len(Hostname) == 0 {
		//gologger.Labelf("Chaos API 未发现到域名 %s\n", domain)
		return ""
	}
	var result string
	result = "{\"passive_dns\":["
	var add int
	for add = 0; add < len(Hostname); add++ {
		result += "{\"hostname\"" + ":" + "\"" + Hostname[add] + "\"" + "},"
		DomainsIP.Domains = append(DomainsIP.Domains, Hostname[add])
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfoChaos(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Chaos Api查询", options)

	return "Success"
}
