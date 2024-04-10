package main

import (
	"OneLong/Utils"
	"OneLong/Utils/Gogogo"
)

func main() {
	var enOptions Utils.LongOptions
	//var Domainip outputfile.DomainsIP
	//var config Port.Config
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	Gogogo.StartScan(&enOptions)
	//Port.DoMasscanPlusNmap(config, &enOptions, &Domainip)

}
