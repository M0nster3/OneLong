package main

import (
	"OneLong/Utils"
	"OneLong/Utils/Gogogo"
	outputfile "OneLong/Utils/OutPutfile"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	//var Domainip outputfile.DomainsIP
	Utils.ConfigParse(&enOptions)
	//var Domainip outputfile.DomainsIP
	if enOptions.KeyWord != "" {
		Gogogo.CompanyRunJob(&enOptions)
	} else {
		Gogogo.DomainRunJob(&enOptions)
	}
	//Domains.Domains("nthu.edu.tw", &enOptions, &Domainip)
	//mobile.pinduoduo.com
	//mobile.pinduoduo.com
	//WaybackarchiveLogin.WaybackarchiveLogin("mobile.pinduoduo.com", &enOptions, &Domainip)

	//Intelx.IntelxEmail("nthu.edu.tw", &enOptions, &Domainip)
	//Tomba.TombaEmail("nthu.edu.tw", &enOptions, &Domainip)
	//Urlscan.Urlscan("nthu.edu.tw", &enOptions, &Domainip)
	//baidu.Baidu("nthu.edu.tw", &enOptions, &Domainip)
	//Email.Email("nthu.edu.tw", &enOptions, &Domainip) //nthu.edu.tw
	////HttpZhiwen.Status("nthu.edu.tw", &enOptions, &Domainip)
	outputfile.OutPutExcelByMergeEnInfo(&enOptions)
	//yahoo.YahooEmail("nthu.edu.tw", &enOptions, &Domainip)
	//Github.GithubLogin("nthu.edu.tw", &enOptions, &Domainip)
	//Port.Port()

}
