package main

import (
	"OneLong/Api/Domains"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	var Domainip outputfile.DomainsIP
	//Gogogo.RunJob(&enOptions)
	//Domains.Domains("baidu.com", &enOptions, &Domainip)
	//IP138.IP138("qianxin.com", &enOptions, &Domainip)
	//NetCraft.NetCraft("qianxin.com", &enOptions, &Domainip)
	//DomainsResult.DomainsResult(domain, &enOptions)
	//Racent.Racent("baidu.com", &enOptions, &Domainip)
	Domains.Domains("google.com", &enOptions, &Domainip)
	//Utils.DomainsResult("baidu.com", &enOptions, &Domainip)
	//Fofa.Fofa("baidu.com", &enOptions, &Domainip)
	//Github.Github("freebuf.com", &enOptions, &Domainip)
	//Gogogo.RunJob(&enOptions)

	outputfile.OutPutExcelByMergeEnInfo(&enOptions)

	//Ehole.Ehole("http://i.3311csci.com", &enOptions, &Domainip)

	//IP.IpWhois("39.106.155.178", &enOptions, &Domainip)
	//Script.Massdns("baidu.com",&enOptions, &Domainip)

	//CDN.IPIPPP("baidu.com")

}
