package Ehole

import (
	"OneLong/Script/Ehole/module/finger"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"strings"
)

func Ehole(domain []string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {
	var Buff []string
	for _, aa := range domain {
		if !strings.Contains(aa, "http://") && !strings.Contains(aa, "https://") {
			Buff = append(Buff, "https://"+aa)
		} else {
			Buff = append(Buff, aa)
		}
	}
	s := finger.NewScan(Buff, 100, options.Proxy)
	s.StartScan(DomainsIP)
}
