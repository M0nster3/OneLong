package Domains

import (
	"OneLong/Api/Domains/alienvault"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"sync"
)

func Domains(domain string, enOptions *Utils.ENOptions, Domainip *outputfile.DomainsIP) {

	var wg sync.WaitGroup
	if domain != "" {
		wg.Add(1)
		go func() {

			alienvault.Alienvault(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		IP138.IP138(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		anubis.Anubis(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		digitorus.Digitorus(domain, enOptions, Domainip) //Digitorus  API 证书查询域名
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		dnsdumpster.Dnsdumpster(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		dnsrepo.Dnsrepo(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		waybackarchive.Waybackarchive(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		Crtsh.Crtsh(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		netlas.Netlas(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		rapiddns.Rapiddns(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		certspotter.Certspotter(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if domain != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		hackertarget.Hackertarget(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.Binaryedge != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//		binaryedge.Binaryedge(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Fullhunt != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		fullhunt.Fullhunt(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.FofaKey != "" && enOptions.ENConfig.Cookies.FofaEmail != "" { //---------------------------------
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		Fofa.Fofa(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.Github != "" { //-----------------------------
	//	wg.Add(1)
	//	go func() {
	//
	//		Github.Github(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Hunter != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		hunter.Hunter(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.Bevigil != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		bevigil.Bevigil(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//////////////////////////
	//if enOptions.ENConfig.Cookies.Racent != "" { //----------------------------
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		Racent.Racent(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Whoisxmlapi != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		whoisxmlapi.Whoisxmlapi(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Virustotal != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		virustotal.Virustotal(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Shodan != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		shodan.Shodan(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Zoomeye != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		ZoomEye.ZoomEye(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.CensysToken != "" && enOptions.ENConfig.Cookies.CensysSecret != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		Censys.Censys(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.Chaos != "" {
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		chaos.Chaos(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Quake != "" {
	//	wg.Add(1)
	//	go func() {
	//
	//		Quake.Quake(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//
	//if enOptions.ENConfig.Cookies.Securitytrails != "" { //-----------------------------------------------------
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		securitytrails.Securitytrails(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//if enOptions.ENConfig.Cookies.GoogleApi != "" && enOptions.ENConfig.Cookies.GoogleID != "" { //------------------------------------------------------
	//
	//	wg.Add(1)
	//	go func() {
	//
	//		Google.Google(domain, enOptions, Domainip)
	//
	//		wg.Done()
	//	}()
	//}
	//commoncrawl.Commoncrawl(domain, enOptions, Domainip)
	//sitedossier.Sitedossier(domain, enOptions, Domainip)
	//leakix.Leakix(domain, enOptions, Domainip)
	//Robtex.Robtex(domain, enOptions, Domainip)
	wg.Wait()
	//Script.Massdns(domain, enOptions, Domainip)
	//HttpZhiwen.Status(domain, enOptions, Domainip) //这里的domain只起到比对

}
