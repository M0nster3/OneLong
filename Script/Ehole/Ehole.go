package Ehole

import (
	"OneLong/Script/Ehole/module/finger"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"os"
	"path/filepath"
	"strings"
)

func Ehole(domain []string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	var Buff []string
	BlackUrl, _ := os.ReadFile(filepath.Join(Utils.GetPathDir(), "Script/Dict/BlackUrl.txt"))
	for _, aa := range domain {
		parts := strings.Split(aa, ".")
		dom := parts[len(parts)-2:]
		if !strings.Contains(string(BlackUrl), strings.Join(dom, ".")) {
			if !strings.Contains(aa, "http://") && !strings.Contains(aa, "https://") {
				Buff = append(Buff, "https://"+aa)
			} else {
				Buff = append(Buff, aa)
			}
		}

	}
	s := finger.NewScan(Buff, 100, options.Proxy)
	s.StartScan(DomainsIP)
}
