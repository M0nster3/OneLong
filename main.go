package main

import (
	"OneLong/Api/Domains/Racent"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	var Domainip outputfile.DomainsIP
	//Domains.Domains(domain, &enOptions)
	Racent.Racent("baidu.com", &enOptions, &Domainip)
	//Fofa.Fofa("baidu.com", &enOptions, &Domainip)
	//Github.Github("freebuf.com", &enOptions, &Domainip)
	//Gogogo.RunJob(&enOptions)

	//outputfile.OutPutExcelByMergeEnInfo(&enOptions)

	//Ehole.Ehole("http://i.3311csci.com", &enOptions, &Domainip)

	//IP.IpWhois("39.106.155.178", &enOptions, &Domainip)
	//Script.Massdns("baidu.com",&enOptions, &Domainip)

	//CDN.IPIPPP("baidu.com")

}
