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

var wg sync.WaitGroup

type FuncInfo struct {
	Func func(string, *Utils.LongOptions, *outputfile.DomainsIP) // 函数签名应匹配你的函数
}

func Email(domain string, enOptions *Utils.LongOptions, Domainip *outputfile.DomainsIP) {
	//color.RGBStyleFromString("244,211,49").Println("\n--------------------探测邮箱--------------------")

	// 创建一个函数切片，包含要执行的函数及其参数
	funcInfos := []FuncInfo{
		{baidu.Baidu},
		{brave.Brave},
		{duckduckgo.Duckduckgo},
		{yahoo.YahooEmail},
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

	if enOptions.LongConfig.Cookies.Github != "" {
		wg.Add(1)
		go func() {

			Github.GithubEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.LongConfig.Email.Emailhunter != "" {
		wg.Add(1)
		go func() {

			hunter.HunterEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	if enOptions.LongConfig.Email.EmailIntelx != "" {
		wg.Add(1)
		go func() {

			Intelx.IntelxEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}
	if enOptions.LongConfig.Email.TombaKey != "" {
		wg.Add(1)
		go func() {

			Tomba.TombaEmail(domain, enOptions, Domainip)

			wg.Done()
		}()
	}

	wg.Wait()

}
