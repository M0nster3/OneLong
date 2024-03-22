package main

import (
	"OneLong/Api/Domains"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/HttpZhiwen"
	"OneLong/Web/Login"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	var Domainip outputfile.DomainsIP
	//Gogogo.RunJob(&enOptions)
	Domains.Domains("cch.org.tw", &enOptions, &Domainip)
	HttpZhiwen.Status("cch.org.tw", &enOptions, &Domainip)
	Login.Login(Domainip.LoginUrlA, &enOptions, &Domainip)
	//Domains.Domains("au.edu.tw", &enOptions, &Domainip)
	//UrlScan.Urlscan("freebuf.com", &enOptions, &Domainip)
	//CommoncrawlLogin.CommoncrawlLogin("cjcu.edu.tw", &enOptions, &Domainip)
	//alienvaultLogin.AlienvaultLogin("cjcu.edu.tw", &enOptions, &Domainip)
	//CommoncrawlLogin.ParseLoginurl(&enOptions, &Domainip)
	//IP138.IP138("qianxin.com", &enOptions, &Domainip)
	//NetCraft.NetCraft("qianxin.com", &enOptions, &Domainip)
	//DomainsResult.DomainsResult(domain, &enOptions)
	//Racent.Racent("baidu.com", &enOptions, &Domainip)
	//Domains.Domains("yuntech.edu.tw", &enOptions, &Domainip)
	//Utils.DomainsResult("baidu.com", &enOptions, &Domainip)
	//Fofa.Fofa("baidu.com", &enOptions, &Domainip)
	//Github.Github("freebuf.com", &enOptions, &Domainip)
	//Gogogo.RunJob(&enOptions)
	//CommoncrawlLogin.CommoncrawlLogin("freebuf.com", &enOptions, &Domainip)
	outputfile.OutPutExcelByMergeEnInfo(&enOptions)

	//Ehole.Ehole("http://i.3311csci.com", &enOptions, &Domainip)

	//IP.IpWhois("39.106.155.178", &enOptions, &Domainip)
	//Script.Massdns("baidu.com",&enOptions, &Domainip)

	//CDN.IPIPPP("baidu.com")

}
