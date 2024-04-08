package main

import (
	"OneLong/Script/Port"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
)

func main() {
	var enOptions Utils.LongOptions
	var Domainip outputfile.DomainsIP
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	//Domains.Github("baidu.com", &enOptions, &Domainip)
	//Gogogo.StartScan(&enOptions)
	var Config Port.Config
	Config.Target = "106.225.163.169"
	//Config.Rate = 5000
	//Config.Port = "1-65535"
	Port.DoMasscanPlusNmap(Config, &enOptions, &Domainip)

	//
	//Port.Port()
	//Email.Email(enOptions.Domain, &enOptions, &Domainip)

}
