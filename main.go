package main

import (
	"OneLong/Api/hunter"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/HttpZhiwen"
)

func main() {
	var enOptions Utils.ENOptions
	var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	//////如果不是API模式就直接运行了
	doamin := "qianxin.com"
	//alienvault.Alienvault(doamin, &enOptions, &Domainip)
	//anubis.Anubis(doamin, &enOptions, &Domainip)
	//binaryedge.Binaryedge(doamin, &enOptions, &Domainip)
	//digitorus.Digitorus(doamin, &enOptions, &Domainip) //Digitorus  API 证书查询域名
	//dnsdumpster.Dnsdumpster(doamin, &enOptions, &Domainip)
	//dnsrepo.Dnsrepo(doamin, &enOptions, &Domainip)
	//fullhunt.Fullhunt(doamin, &enOptions, &Domainip)
	hunter.Hunter(doamin, &enOptions, &Domainip)
	//bevigil.Bevigil(doamin, &enOptions, &Domainip)
	//whoisxmlapi.Whoisxmlapi(doamin, &enOptions, &Domainip)
	//waybackarchive.Waybackarchive(doamin, &enOptions, &Domainip)
	//virustotal.Virustotal(doamin, &enOptions, &Domainip)
	//sitedossier.Sitedossier(doamin, &enOptions, &Domainip)
	//shodan.Shodan(doamin, &enOptions, &Domainip)
	//Robtex.Robtex(doamin, &enOptions, &Domainip)
	//ZoomEye.ZoomEye(doamin, &enOptions, &Domainip)
	//Censys.Censys(doamin, &enOptions, &Domainip)
	//chaos.Chaos(doamin, &enOptions, &Domainip)
	//commoncrawl.Commoncrawl(doamin, &enOptions, &Domainip)
	//Crtsh.Crtsh(doamin, &enOptions, &Domainip)
	//hackertarget.Hackertarget(doamin, &enOptions, &Domainip)
	//leakix.Leakix(doamin, &enOptions, &Domainip)
	//netlas.Netlas(doamin, &enOptions, &Domainip)
	//Quake.Quake(doamin, &enOptions, &Domainip)
	//rapiddns.Rapiddns(doamin, &enOptions, &Domainip)
	//certspotter.Certspotter(doamin, &enOptions, &Domainip)
	//Domainip.IP = Utils.SetStr(Domainip.IP)
	//Domainip.Domains = Utils.SetStr(Domainip.Domains)
	HttpZhiwen.Status(doamin, &enOptions, &Domainip) //这里的domain只起到比对
	outputfile.OutPutExcelByMergeEnInfo(&enOptions)

	//Ehole.Ehole("http://i.3311csci.com", &enOptions, &Domainip)
	//Gogogo.RunJob(&enOptions)
	//IP.IpWhois("39.106.155.178", &enOptions, &Domainip)
	//Script.Massdns("baidu.com",&enOptions, &Domainip)

	//CDN.IPIPPP("baidu.com")

}
