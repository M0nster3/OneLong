package Ehole

import (
	"OneLong/Script/Ehole/module/finger"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
)

func Ehole(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {

	s := finger.NewScan([]string{domain}, 100, options.Proxy)
	s.StartScan(DomainsIP)
}
