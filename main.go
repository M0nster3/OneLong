package main

import (
	"OneLong/Script"
	"OneLong/Utils"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	////如果不是API模式就直接运行了
	//doamin := "sadadadadad.com"
	//alienvault.Alienvault(doamin, &enOptions)
	////anubis.Anubis(doamin, &enOptions)
	////binaryedge.Binaryedge(doamin, &enOptions)
	////digitorus.Digitorus(doamin, &enOptions) //Digitorus  API 证书查询域名
	////dnsdumpster.Dnsdumpster(doamin, &enOptions)
	////dnsrepo.Dnsrepo(doamin, &enOptions)
	////fullhunt.Fullhunt(doamin, &enOptions)
	////hunter.Hunter(doamin, &enOptions)
	////bevigil.Bevigil(doamin, &enOptions)
	////whoisxmlapi.Whoisxmlapi(doamin, &enOptions)
	////waybackarchive.Waybackarchive(doamin, &enOptions)
	////virustotal.Virustotal(doamin, &enOptions)
	////sitedossier.Sitedossier(doamin, &enOptions)
	////shodan.Shodan(doamin, &enOptions)
	////Robtex.Robtex(doamin, &enOptions)
	////ZoomEye.ZoomEye(doamin, &enOptions)
	////Censys.Censys(doamin, &enOptions)
	////chaos.Chaos(doamin, &enOptions)
	////commoncrawl.Commoncrawl(doamin, &enOptions)
	////Crtsh.Crtsh(doamin, &enOptions)
	////hackertarget.Hackertarget(doamin, &enOptions)
	////leakix.Leakix(doamin, &enOptions)
	////netlas.Netlas(doamin, &enOptions)
	////Quake.Quake(doamin, &enOptions)
	////rapiddns.Rapiddns(doamin, &enOptions)
	//certspotter.Certspotter(doamin, &enOptions)
	//outputfile.OutPutExcelByMergeEnInfo(&enOptions)
	////Gogogo.RunJob(&enOptions)
	Script.Massdns("baidu.com", &enOptions)

}
