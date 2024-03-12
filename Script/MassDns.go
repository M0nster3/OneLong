package Script

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"fmt"
	"github.com/projectdiscovery/shuffledns/pkg/runner"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func checkWildcardDNS(domain string) bool {
	// 生成一个很可能不存在的子域名
	randomSubdomain := "random-subdomain-that-probably-doesnt-exist." + domain

	// 尝试解析这个子域名
	_, err := net.LookupIP(randomSubdomain)
	if err != nil {
		// 解析错误，很可能是因为不存在泛解析

		return false
	}

	// 解析成功，打印解析到的IP地址

	return true
}

func Massdns(domain string, ENOptions *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	if checkWildcardDNS(domain) {
		gologger.Infof("[存在泛解析不进行爆破]\n")
		return
	} else {
		gologger.Infof("[不是泛解析將爆破域名]\n")
	}
	tempOutputFile := Utils.GetTempPathFileName()
	defer os.Remove(tempOutputFile)
	tempDir, err := os.MkdirTemp("", Utils.RangeString(8))
	if err != nil {
		gologger.Errorf("%s", err)
	}
	options := &runner.Options{
		Directory:          tempDir,
		Domain:             domain,
		SubdomainsList:     "",
		ResolversFile:      filepath.Join(Utils.GetPathDir(), "Script/MsassDns/", ENOptions.ENConfig.Massdns.Resolvers),
		Wordlist:           filepath.Join(Utils.GetPathDir(), "Script/MsassDns/", ENOptions.ENConfig.Massdns.Wordlist),
		MassdnsPath:        filepath.Join(Utils.GetPathDir(), "Script/MsassDns/", ENOptions.ENConfig.Massdns.MassdnsPath),
		Output:             tempOutputFile,
		Json:               false,
		Silent:             false,
		Version:            false,
		Retries:            5,
		Verbose:            true,
		NoColor:            true,
		Threads:            300,
		MassdnsRaw:         "",
		WildcardThreads:    25,
		StrictWildcard:     true,
		WildcardOutputFile: "",
		Stdin:              false,
	}
	massdnsRunner, err := runner.New(options)
	if err != nil {
		gologger.Errorf("Could not create runner: %s", err)
	}

	massdnsRunner.RunEnumeration()
	content, err := os.ReadFile(tempOutputFile)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(content), "\n") {
		domains := strings.TrimSpace(line)
		DomainsIP.Domains = append(DomainsIP.Domains, domains)
		fmt.Printf(domains)
	}
}
