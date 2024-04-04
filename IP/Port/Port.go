package Port

import (
	"OneLong/Utils"
	"context"
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

func Port(host string, ENOptions *Utils.ENOptions) {
	var PortProxy string
	if ENOptions.Proxy != "" {
		PortProxy = "127.0.0.1:7890"
	}
	options := runner.Options{
		Host:     goflags.StringSlice{"64.190.196.28"},
		ScanType: "s",
		OnResult: func(hr *result.HostResult) {
			log.Println(hr.Host, hr.Ports)
		},
		Proxy:    PortProxy,
		TopPorts: "1000",
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer naabuRunner.Close()
	err = naabuRunner.RunEnumeration(context.Background())
	if err != nil {
		return
	}
}
