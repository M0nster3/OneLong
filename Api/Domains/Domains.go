package Domains

import (
	"OneLong/Api/Domains/Censys"
	"OneLong/Api/Domains/Crtsh"
	"OneLong/Api/Domains/Fofa"
	"OneLong/Api/Domains/Github"
	"OneLong/Api/Domains/Google"
	"OneLong/Api/Domains/IP138"
	"OneLong/Api/Domains/Quake"
	"OneLong/Api/Domains/Racent"
	"OneLong/Api/Domains/Robtex"
	"OneLong/Api/Domains/Urlscan"
	"OneLong/Api/Domains/ZoomEye"
	"OneLong/Api/Domains/alienvault"
	"OneLong/Api/Domains/anubis"
	"OneLong/Api/Domains/bevigil"
	"OneLong/Api/Domains/binaryedge"
	"OneLong/Api/Domains/certspotter"
	"OneLong/Api/Domains/chaos"
	"OneLong/Api/Domains/commoncrawl"
	"OneLong/Api/Domains/digitorus"
	"OneLong/Api/Domains/dnsdumpster"
	"OneLong/Api/Domains/dnsrepo"
	"OneLong/Api/Domains/fullhunt"
	"OneLong/Api/Domains/hackertarget"
	"OneLong/Api/Domains/hunter"
	"OneLong/Api/Domains/leakix"
	"OneLong/Api/Domains/netlas"
	"OneLong/Api/Domains/rapiddns"
	"OneLong/Api/Domains/securitytrails"
	"OneLong/Api/Domains/shodan"
	"OneLong/Api/Domains/sitedossier"
	"OneLong/Api/Domains/virustotal"
	"OneLong/Api/Domains/waybackarchive"
	"OneLong/Api/Domains/whoisxmlapi"
	"OneLong/Script"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"github.com/gookit/color"
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
	if domain != "" {
		wg.Add(1)
		go func() {

			Urlscan.Urlscan(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			IP138.IP138(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			anubis.Anubis(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			digitorus.Digitorus(domain, enOptions, Domainip) //Digitorus  API 证书查询域名

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			dnsdumpster.Dnsdumpster(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			dnsrepo.Dnsrepo(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if domain != "" {
		wg.Add(1)
		go func() {

			waybackarchive.Waybackarchive(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if domain != "" {
		wg.Add(1)
		go func() {

			Crtsh.Crtsh(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			netlas.Netlas(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			rapiddns.Rapiddns(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if domain != "" {
		wg.Add(1)
		go func() {

			certspotter.Certspotter(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			hackertarget.Hackertarget(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.Binaryedge != "" {

		wg.Add(1)
		go func() {
			binaryedge.Binaryedge(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Fullhunt != "" {
		wg.Add(1)
		go func() {

			fullhunt.Fullhunt(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.FofaKey != "" && enOptions.ENConfig.Cookies.FofaEmail != "" { //---------------------------------

		wg.Add(1)
		go func() {

			Fofa.Fofa(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.Github != "" { //-----------------------------
		wg.Add(1)
		go func() {

			Github.Github(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Hunter != "" {
		wg.Add(1)
		go func() {

			hunter.Hunter(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.Bevigil != "" {
		wg.Add(1)
		go func() {

			bevigil.Bevigil(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	////////////////////////
	if enOptions.ENConfig.Cookies.Racent != "" { //----------------------------

		wg.Add(1)
		go func() {

			Racent.Racent(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Whoisxmlapi != "" {

		wg.Add(1)
		go func() {

			whoisxmlapi.Whoisxmlapi(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Virustotal != "" {

		wg.Add(1)
		go func() {

			virustotal.Virustotal(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Shodan != "" {

		wg.Add(1)
		go func() {

			shodan.Shodan(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Zoomeye != "" {

		wg.Add(1)
		go func() {

			ZoomEye.ZoomEye(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.CensysToken != "" && enOptions.ENConfig.Cookies.CensysSecret != "" {
		wg.Add(1)
		go func() {

			Censys.Censys(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.Chaos != "" {

		wg.Add(1)
		go func() {

			chaos.Chaos(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Quake != "" {
		wg.Add(1)
		go func() {

			Quake.Quake(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Securitytrails != "" { //-----------------------------------------------------

		wg.Add(1)
		go func() {

			securitytrails.Securitytrails(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Cookies.GoogleApi != "" && enOptions.ENConfig.Cookies.GoogleID != "" { //------------------------------------------------------

		wg.Add(1)
		go func() {

			Google.Google(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	commoncrawl.Commoncrawl(domain, enOptions, Domainip)
	sitedossier.Sitedossier(domain, enOptions, Domainip)
	leakix.Leakix(domain, enOptions, Domainip)
	Robtex.Robtex(domain, enOptions, Domainip)
	wg.Wait()
	if !enOptions.NoBao {
		color.RGBStyleFromString("205,155,29").Println("\n--------------------Massdns爆破子域名--------------------")
		Script.Massdns(domain, enOptions, Domainip)
	}

}
