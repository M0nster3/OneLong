package main

import (
	"OneLong/Email/duckduckgo"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	var Domainip outputfile.DomainsIP
	//Utils.ConfigParse(&enOptions)
	////var Domainip outputfile.DomainsIP
	//if enOptions.KeyWord != "" {
	//	Gogogo.CompanyRunJob(&enOptions)
	//} else {
	//	Gogogo.DomainRunJob(&enOptions)
	//}
	//yahoo.Yahoo("qianxin.com", &enOptions, &Domainip)
	duckduckgo.Duckduckgo("nthu.edu.tw", &enOptions, &Domainip)
	//Port.Port()

}
