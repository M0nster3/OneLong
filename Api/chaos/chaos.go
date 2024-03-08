package chaos

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"github.com/projectdiscovery/chaos-client/pkg/chaos"
	"github.com/tidwall/gjson"
	"strings"
)

func GetEnInfo(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {

	//respons := gjson.Get(response, "events").Array()
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "[", "")
	//respons := gjson.Parse(response).Array()
	respons := gjson.Get(response, "passive_dns").Array()
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Chaos"
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
	//zuo := strings.ReplaceAll(response, "[", "")
	//you := strings.ReplaceAll(zuo, "]", "")

	//ensInfos.Infos["hostname"] = append(ensInfos.Infos["hostname"], gjson.Parse(Result[1].String()))
	//getCompanyInfoById(pid, 1, true, "", options.GetField, ensInfos, options)
	return ensInfos, ensOutMap

}

func Chaos(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	gologger.Infof("Chaos API 域名查询 \n")
	var Hostname []string
	chaosClient := chaos.New(options.ENConfig.Cookies.Chaos)
	for result := range chaosClient.GetSubdomains(&chaos.SubdomainsRequest{
		Domain: domain,
	}) {
		if result.Error != nil {
			break
		}
		res := strings.ReplaceAll(result.Subdomain, "*.", "")
		Hostname = append(Hostname, res+"."+domain)

		//results <- subscraping.Result{
		//	Source: s.Name(), Type: subscraping.Subdomain, Value: fmt.Sprintf("%s.%s", result.Subdomain, domain),
		//}
		//s.results++
	}
	if len(Hostname) == 0 {
		gologger.Labelf("Chaos API 未发现到域名\n")
		return ""
	}
	var result string
	result = "{\"passive_dns\":["
	var add int
	for add = 0; add < len(Hostname); add++ {
		result += "{\"hostname\"" + ":" + "\"" + Hostname[add] + "\"" + "},"
	}
	result = result + "]}"
	res, ensOutMap := GetEnInfo(result, DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Chaos Api查询", options)
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
