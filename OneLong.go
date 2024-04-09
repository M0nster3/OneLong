package main

import (
	"OneLong/Utils"
	"OneLong/Utils/Gogogo"
)

func main() {
	var enOptions Utils.LongOptions
	//var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	//Domains.Github("baidu.com", &enOptions, &Domainip)
	Gogogo.StartScan(&enOptions)

}
