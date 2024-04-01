package Email

import (
	"OneLong/Email/Github"
	"OneLong/Email/Intelx"
	"OneLong/Email/Tomba"
	"OneLong/Email/baidu"
	"OneLong/Email/brave"
	"OneLong/Email/duckduckgo"
	"OneLong/Email/hunter"
	"OneLong/Email/yahoo"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"sync"
)

func Email(domain string, enOptions *Utils.ENOptions, Domainip *outputfile.DomainsIP) {
	//color.RGBStyleFromString("244,211,49").Println("\n--------------------探测邮箱--------------------")
	var wg sync.WaitGroup
	if domain != "" {
		wg.Add(1)
		go func() {

			baidu.Baidu(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if domain != "" {
		wg.Add(1)
		go func() {

			brave.Brave(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if domain != "" {
		wg.Add(1)
		go func() {

			duckduckgo.Duckduckgo(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Cookies.Github != "" {
		wg.Add(1)
		go func() {

			Github.GithubEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Email.Emailhunter != "" {
		wg.Add(1)
		go func() {

			hunter.HunterEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.ENConfig.Email.EmailIntelx != "" {
		wg.Add(1)
		go func() {

			Intelx.IntelxEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.ENConfig.Email.TombaKey != "" {
		wg.Add(1)
		go func() {

			Tomba.TombaEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if domain != "" {
		wg.Add(1)
		go func() {

			yahoo.YahooEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	wg.Wait()

}
