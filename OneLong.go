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
	Gogogo.StartScan(&enOptions)
	//var Config Port.Config
	//Config.Target = ""
	//Config.Rate = 5000
	//Config.Port = "1-65535"
	//Port.DoMasscanPlusNmap(Config)

	//
	//Port.Port()
	//Email.Email(enOptions.Domain, &enOptions, &Domainip)

}
