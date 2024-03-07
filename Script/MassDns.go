package Script

import (
	"OneLong/Utils"
	"OneLong/Utils/gologger"
	"fmt"
	"github.com/projectdiscovery/shuffledns/pkg/runner"
	"os"
	"path/filepath"
	"strings"
)

func Massdns(domain string, ENOptions *Utils.ENOptions) {
	tempOutputFile := Utils.GetTempPathFileName()
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
		fmt.Printf(domains)
	}
}
