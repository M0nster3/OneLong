package Domains

import (
	"OneLong/Script"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"github.com/gookit/color"
	"sync"
)

var wg sync.WaitGroup

type FuncInfo struct {
	Func func(string, *Utils.LongOptions, *outputfile.DomainsIP) string // 函数签名应匹配你的函数
}

func Domains(domain string, enOptions *Utils.LongOptions, Domainip *outputfile.DomainsIP) {

	// 创建一个函数切片，包含要执行的函数及其参数
	funcInfos := []FuncInfo{
		{Alienvault},
		{Urlscan},
		{IP138},
		{Anubis},
		{Digitorus},
		{Dnsdumpster},
		{Dnsrepo},
		{Waybackarchive},
		{Crtsh},
		{Netlas},
		{Rapiddns},
		{Certspotter},
		{Hackertarget},
	}
	// 为每个非空域名启动一个 goroutine
	if domain != "" {
		for _, info := range funcInfos {
			wg.Add(1)
			go func(fInfo FuncInfo) {
				fInfo.Func(domain, enOptions, Domainip)
				wg.Done()
			}(info)
		}
	}
	if enOptions.LongConfig.Cookies.Binaryedge != "" {

		wg.Add(1)
		go func() {
			Binaryedge(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Fullhunt != "" {
		wg.Add(1)
		go func() {
			Fullhunt(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.FofaKey != "" && enOptions.LongConfig.Cookies.FofaEmail != "" {

		wg.Add(1)
		go func() {
			Fofa(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.Github != "" {
		wg.Add(1)
		go func() {
			Github(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Hunter != "" {
		wg.Add(1)
		go func() {
			Hunter(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.Bevigil != "" {
		wg.Add(1)
		go func() {
			Bevigil(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Racent != "" {
		wg.Add(1)
		go func() {
			Racent(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Whoisxmlapi != "" {
		wg.Add(1)
		go func() {
			Whoisxmlapi(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Virustotal != "" {

		wg.Add(1)
		go func() {
			Virustotal(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Shodan != "" {

		wg.Add(1)
		go func() {
			Shodan(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Zoomeye != "" {

		wg.Add(1)
		go func() {
			ZoomEye(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.CensysToken != "" && enOptions.LongConfig.Cookies.CensysSecret != "" {
		wg.Add(1)
		go func() {

			Censys(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.Chaos != "" {

		wg.Add(1)
		go func() {
			Chaos(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Quake != "" {
		wg.Add(1)
		go func() {
			Quake(domain, enOptions, Domainip)
			wg.Done()
		}()
	}

	if enOptions.LongConfig.Cookies.Securitytrails != "" {

		wg.Add(1)
		go func() {
			Securitytrails(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	if enOptions.LongConfig.Cookies.GoogleApi != "" && enOptions.LongConfig.Cookies.GoogleID != "" {

		wg.Add(1)
		go func() {
			Google(domain, enOptions, Domainip)
			wg.Done()
		}()
	}
	Commoncrawl(domain, enOptions, Domainip)
	Sitedossier(domain, enOptions, Domainip)
	Leakix(domain, enOptions, Domainip)
	Robtex(domain, enOptions, Domainip)
	wg.Wait()
	if !enOptions.NoBao {
		color.RGBStyleFromString("205,155,29").Println("\n--------------------Massdns爆破子域名--------------------")
		Script.Massdns(domain, enOptions, Domainip)
	}

}
