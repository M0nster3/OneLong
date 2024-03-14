package main

import (
	"OneLong/Api/Crtsh"
	"OneLong/Api/Quake"
	"OneLong/Api/Robtex"
	"OneLong/Api/ZoomEye"
	"OneLong/Api/alienvault"
	"OneLong/Api/anubis"
	"OneLong/Api/bevigil"
	"OneLong/Api/binaryedge"
	"OneLong/Api/certspotter"
	"OneLong/Api/chaos"
	"OneLong/Api/commoncrawl"
	"OneLong/Api/digitorus"
	"OneLong/Api/dnsdumpster"
	"OneLong/Api/dnsrepo"
	"OneLong/Api/fullhunt"
	"OneLong/Api/hackertarget"
	"OneLong/Api/leakix"
	"OneLong/Api/netlas"
	"OneLong/Api/rapiddns"
	"OneLong/Api/shodan"
	"OneLong/Api/sitedossier"
	"OneLong/Api/virustotal"
	"OneLong/Api/waybackarchive"
	"OneLong/Api/whoisxmlapi"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/HttpZhiwen"
	"sync"
)

func main() {
	var enOptions Utils.ENOptions
	var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	//////如果不是API模式就直接运行了
	doamin := "mmh.org.tw"
	var wg sync.WaitGroup
	wg.Add(26)
	go func() {
		alienvault.Alienvault(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		anubis.Anubis(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		binaryedge.Binaryedge(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		digitorus.Digitorus(doamin, &enOptions, &Domainip) //Digitorus  API 证书查询域名

		wg.Done()
	}()
	go func() {
		dnsdumpster.Dnsdumpster(doamin, &enOptions, &Domainip)

		wg.Done()
	}()
	go func() {
		dnsrepo.Dnsrepo(doamin, &enOptions, &Domainip)

		wg.Done()
	}()
	go func() {
		fullhunt.Fullhunt(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		//hunter.Hunter(doamin, &enOptions, &Domainip)

		wg.Done()
	}()
	go func() {
		bevigil.Bevigil(doamin, &enOptions, &Domainip)

		wg.Done()
	}()
	go func() {
		whoisxmlapi.Whoisxmlapi(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		waybackarchive.Waybackarchive(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		virustotal.Virustotal(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		sitedossier.Sitedossier(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		shodan.Shodan(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		Robtex.Robtex(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		ZoomEye.ZoomEye(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		//Censys.Censys(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		chaos.Chaos(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		commoncrawl.Commoncrawl(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		Crtsh.Crtsh(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		hackertarget.Hackertarget(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		leakix.Leakix(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		netlas.Netlas(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		Quake.Quake(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		rapiddns.Rapiddns(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	go func() {
		certspotter.Certspotter(doamin, &enOptions, &Domainip)
		wg.Done()
	}()
	wg.Wait()

	HttpZhiwen.Status(doamin, &enOptions, &Domainip) //这里的domain只起到比对
	outputfile.OutPutExcelByMergeEnInfo(&enOptions)

	//Ehole.Ehole("http://i.3311csci.com", &enOptions, &Domainip)
	//Gogogo.RunJob(&enOptions)
	//IP.IpWhois("39.106.155.178", &enOptions, &Domainip)
	//Script.Massdns("baidu.com",&enOptions, &Domainip)

	//CDN.IPIPPP("baidu.com")

}
