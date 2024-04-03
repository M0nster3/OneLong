package main

import (
	"OneLong/Utils"
	"OneLong/Utils/Gogogo"
)

func main() {
	var enOptions Utils.ENOptions
	//var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	Gogogo.StartScan(&enOptions)
	//Email.Email(enOptions.Domain, &enOptions, &Domainip)

}
