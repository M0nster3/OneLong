// package main
//
// import (
//
//	"OneLong/IP/Port"
//	"OneLong/Utils"
//	"OneLong/Utils/Gogogo"
//
// )
//
//	func main() {
//		var enOptions Utils.ENOptions
//		//var Domainip outputfile.DomainsIP
//		Utils.Flag(&enOptions)
//		Utils.ConfigParse(&enOptions)
//		Gogogo.StartScan(&enOptions)
//		var Config Port.Config
//		Config.Target = "107.163.229.83/24"
//		Config.Rate = 2000
//		Config.Port = "1-65535"
//		Port.DoMasscanPlusNmap(Config)
//
//		//
//		//Port.Port()
//		//Email.Email(enOptions.Domain, &enOptions, &Domainip)
//
// }
package main

import (
	"fmt"
	"net"
)

func main() {
	ips := []string{"192.168.0.1", "10.0.0.1", "172.16.0.1", "ssss"}

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Printf("Invalid IP address: %s\n", ipStr)
			continue
		}

		// 提取C段
		cidr := fmt.Sprintf("%s/24", ip.String()) // 使用/24表示C段
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("Error parsing CIDR: %s\n", err)
			continue
		}

		// 输出C段
		fmt.Printf("IP: %s, C段: %s\n", ip.String(), ipnet.String())
	}
}
