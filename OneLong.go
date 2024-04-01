package main

import (
	"OneLong/Utils"
	"OneLong/Utils/Gogogo"
)

func main() {
	var enOptions Utils.ENOptions
	Utils.Flag(&enOptions)
	Utils.ConfigParse(&enOptions)
	if enOptions.KeyWord != "" {
		Gogogo.CompanyRunJob(&enOptions)
	} else {
		Gogogo.DomainRunJob(&enOptions)
	}

}
