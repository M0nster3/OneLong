package main

import (
	"OneLong/IP/Port"
	"OneLong/Utils"
)

func main() {
	var enOptions Utils.ENOptions
	//var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	//Gogogo.StartScan(&enOptions)
	var Config Port.Config
	Config.Target = ""
	Config.Rate = 5000
	Config.Port = "1-65535"
	Port.DoMasscanPlusNmap(Config)

	//
	//Port.Port()
	//Email.Email(enOptions.Domain, &enOptions, &Domainip)

}
