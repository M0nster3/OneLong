package CNAME

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Structure for each provider stored in providers.json file
type ProviderData struct {
	Name     string   `json:"name"`
	Cname    []string `json:"cname"`
	Response []string `json:"response"`
}

var Providers []ProviderData

var Targets []string

var (
	HostsList  string
	Threads    = 20
	All        bool
	Verbose    = true
	ForceHTTPS bool
	Timeout    int
	OutputFile string
)

func InitializeProviders() {
	raw, err := ioutil.ReadFile(filepath.Join(Utils.GetPathDir(), "Script/CName/cname.json"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &Providers)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func ReadFile(file string) (lines []string, err error) {
	fileHandle, err := os.Open(file)
	if err != nil {
		return lines, err
	}

	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines, nil
}

func Get(url string, timeout int, https bool) (resp gorequest.Response, body string, errs []error) {
	if https == true {
		url = fmt.Sprintf("https://%s/", url)
	} else {
		url = fmt.Sprintf("http://%s/", url)
	}

	resp, body, errs = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Timeout(time.Duration(timeout)*time.Second).Get(url).
		Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0").
		End()

	return resp, body, errs
}

func CNAMEExists(key string) bool {
	for _, provider := range Providers {
		for _, cname := range provider.Cname {
			if strings.Contains(key, cname) {
				return true
			}
		}
	}

	return false
}

func Check(target string, TargetCNAME string) {
	_, body, errs := Get(target, Timeout, ForceHTTPS)
	if len(errs) <= 0 {
		if TargetCNAME == "ALL" {
			for _, provider := range Providers {
				for _, response := range provider.Response {
					if strings.Contains(body, response) == true {
						fmt.Printf("\n[\033[31;1;4m%s\033[0m] Takeover Possible At %s ", provider.Name, target)
						return
					}
				}
			}
		} else {
			// This is a less false positives way
			for _, provider := range Providers {
				for _, cname := range provider.Cname {
					if strings.Contains(TargetCNAME, cname) {
						for _, response := range provider.Response {
							if strings.Contains(body, response) == true {
								if provider.Name == "cloudfront" {
									_, body2, _ := Get(target, 120, true)
									if strings.Contains(body2, response) == true {
										fmt.Printf("\n[\033[31;1;4m%s\033[0m] Takeover Possible At : %s", provider.Name, target)
									}
								} else {
									fmt.Printf("\n[\033[31;1;4m%s\033[0m] Takeover Possible At %s with CNAME %s", provider.Name, target, TargetCNAME)
								}
							}
							return
						}
					}
				}
			}
		}
	} else {
		if Verbose == true {
			log.Printf("[ERROR] Get: %s => %v", target, errs)
		}
	}

	return
}

func Checker(target string) {
	TargetCNAME, err := net.LookupCNAME(target)
	if err != nil {
		return
	} else {
		if All != true && CNAMEExists(TargetCNAME) == true {
			if Verbose == true {
				log.Printf("[SELECTED] %s => %s", target, TargetCNAME)
			}
			Check(target, TargetCNAME)
		} else if All == true {
			if Verbose == true {
				log.Printf("[ALL] %s ", target)
			}
			Check(target, "ALL")
		}
	}
}

func Cname(DomainsIP *outputfile.DomainsIP) {

	InitializeProviders()

	Targets = DomainsIP.Domains

	hosts := make(chan string, Threads)
	processGroup := new(sync.WaitGroup)
	processGroup.Add(Threads)

	for i := 0; i < Threads; i++ {
		go func() {
			for {
				host := <-hosts
				if host == "" {
					break
				}

				Checker(host)
			}

			processGroup.Done()
		}()
	}

	for _, Host := range Targets {
		hosts <- Host
	}

	close(hosts)
	processGroup.Wait()

	fmt.Printf("\n[~] Enjoy your hunt !\n")
}
