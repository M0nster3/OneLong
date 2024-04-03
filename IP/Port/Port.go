package Port

import (
	"context"
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

func Port() {
	options := runner.Options{
		Host:     goflags.StringSlice{""},
		ScanType: "s",
		OnResult: func(hr *result.HostResult) {
			log.Println(hr.Host, hr.Ports)
		},

		TopPorts: "full",
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
