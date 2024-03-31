package Afrog

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zan8in/afrog/v3"
	"os"
	"path/filepath"
)

// 用于保护 addedURLs
func GetEnInfoAfrog(response string, DomainsIP *outputfile.DomainsIP) (*Utils.EnInfos, map[string]*outputfile.ENSMap) {
	ensInfos := &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensInfos.SType = "Afrog"
	ensOutMap := make(map[string]*outputfile.ENSMap)
	for k, v := range GetENMap() {
		ensOutMap[k] = &outputfile.ENSMap{Name: v.Name, Field: v.Field, KeyWord: v.KeyWord}
	}
	var targeturl []string
	var targetinfoname []string
	var targetinfoseg []string
	res := gjson.Get(response, "#.fulltarget").Array()
	infoname := gjson.Get(response, "#.pocinfo.infoname").Array()
	infoseg := gjson.Get(response, "#.pocinfo.infoseg").Array()
	for a, r := range res {
		infon := infoname[a].String()
		infoe := infoseg[a].String()
		targeturl = append(targeturl, r.String())
		targetinfoname = append(targetinfoname, infon)
		targetinfoseg = append(targetinfoseg, infoe)
	}

	//Result := gjson.GetMany(response, "passive_dns.#.address", "passive_dns.#.hostname")
	//ensInfos.Infoss = make(map[string][]map[string]string)
	//获取公司信息
	//ensInfos.Infos["passive_dns"] = append(ensInfos.Infos["passive_dns"], gjson.Parse(Result[0].String()))
	//Field:   []string{"Url", "info", "infoname"},
	for aa, _ := range targeturl {
		ResponseJia := fmt.Sprintf("{\"Url\": \"%s\", \"info\": \"%s\", \"infoname\": \"%s\"}", targeturl[aa], targetinfoseg[aa], targetinfoname[aa])
		ensInfos.Infos["Afrog"] = append(ensInfos.Infos["Afrog"], gjson.Parse(ResponseJia))
	}

	return ensInfos, ensOutMap

}

func Afrog(options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	var SetProxy string
	if options.Proxy != "" {
		SetProxy = options.Proxy
	}

	tempOutputFile := Utils.GetTempJsonPathFileName()
	tempInputFile := Utils.GetTempPathFileName()
	file, _ := os.OpenFile(tempInputFile, os.O_CREATE|os.O_WRONLY, 0644)
	for _, url := range DomainsIP.DomainA {
		_, err := file.WriteString(url)
		if err != nil {
			return
		}
	}
	defer file.Close()
	defer os.Remove(tempOutputFile)
	defer os.Remove(tempInputFile)
	if err := afrog.NewScanner([]string{}, afrog.Scanner{
		//Severity:  "High",
		//Search:      "shiro",
		TargetsFile: tempInputFile,
		//Json:        tempOutputFile,
		Smart: true, //智能控制并发
		Json:  tempOutputFile,
		//Json:       "/Users/Monster/Downloads/afrog-main/examples/batch_scan/3.json",
		Proxy:      SetProxy,
		AppendPoc:  []string{filepath.Join(Utils.GetPathDir(), "Script/AfrogAddPoc")},
		ConfigFile: filepath.Join(Utils.GetPathDir(), "config.yaml"),
	}); err != nil {
		fmt.Println(err.Error())
	}

	resfile, err := os.ReadFile(tempOutputFile)
	if err != nil {
		return
	}

	res, ensOutMap := GetEnInfoAfrog(string(resfile), DomainsIP)

	outputfile.MergeOutPut(res, ensOutMap, "Afrog", options)
}
