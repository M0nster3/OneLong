package Port

import (
	"OneLong/Utils"
	"bytes"
	"fmt"
	gonmap "github.com/lair-framework/go-nmap"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// NewNmap 创建nmap对象
func NewNmap(config Config) *Nmap {
	config.CmdBin = "nmap"
	if runtime.GOOS == "windows" {
		config.CmdBin = "nmap.exe"
	}
	return &Nmap{Config: config}
}

type Service struct {
	nmapServiceData   map[string]string
	customServiceData map[string]string
}

// loadNmapService 加载nmap的service定义
func (s *Service) loadNmapService() {
	content, err := os.ReadFile(filepath.Join(Utils.GetPathDir(), "Script/Nmap/nmap-services.txt"))
	if err != nil {
		//logging.RuntimeLog.Info(err)
		//logging.CLILog.Info(err)
	} else {
		for _, line := range strings.Split(string(content), "\n") {
			txt := strings.TrimSpace(line)
			if txt == "" || strings.HasPrefix(txt, "#") {
				continue
			}
			servicesListArray := strings.Split(txt, "\t")
			if len(servicesListArray) < 3 {
				continue
			}
			s.nmapServiceData[servicesListArray[1]] = servicesListArray[0]
		}
	}
}

// loadCustomService 加载自定义service定义文件
func (s *Service) loadCustomService() {
	content, err := os.ReadFile(filepath.Join(Utils.GetPathDir(), "Script/Port/Nmap/services-custom.txt"))
	if err != nil {
		//logging.RuntimeLog.Info(err)
		//logging.CLILog.Info(err)
	} else {
		for _, line := range strings.Split(string(content), "\n") {
			txt := strings.TrimSpace(line)
			if txt == "" || strings.HasPrefix(txt, "#") {
				continue
			}
			servicesListArray := strings.Split(txt, " ")
			if len(servicesListArray) < 2 {
				continue
			}

			s.customServiceData[strings.TrimSpace(servicesListArray[0])] = strings.TrimSpace(servicesListArray[1])
		}
	}
}

// NewService 创建service对象
func NewService() Service {
	s := Service{
		nmapServiceData:   make(map[string]string),
		customServiceData: make(map[string]string),
	}
	s.loadNmapService()
	s.loadCustomService()

	return s
}

// getResultIPPortMap 提取扫描结果的ip和port
func getResultIPPortMap(result map[string]*IPResult) (ipPortMap map[string]string) {
	ipPortMap = make(map[string]string)
	for ip, r := range result {
		var ports []string
		for p, _ := range r.Ports {
			ports = append(ports, fmt.Sprintf("%d", p))
		}
		if len(ports) > 0 {
			ipPortMap[ip] = strings.Join(ports, ",")
		}
	}
	return
}

var fpNmapThreadNumber = make(map[string]int)

const (
	HighPerformance   = "High"
	NormalPerformance = "Normal"
)

// WorkerPerformanceMode worker默认的性能模式为Normal
var WorkerPerformanceMode = NormalPerformance

type Nmap struct {
	Config Config
	Result Result
}

func CheckIPV4(ip string) bool {
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
	r, _ := regexp.Compile(ipReg)

	return r.MatchString(ip)

}

// CheckIPV4Subnet 检查是否是ipv4地址段
func CheckIPV4Subnet(ip string) bool {
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])/\d{1,2}$`
	r, _ := regexp.Compile(ipReg)

	return r.MatchString(ip)
}

// CheckIPV6 检查是否是ipv6地址
func CheckIPV6(ip string) bool {
	ipv6Regex := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	match, _ := regexp.MatchString(ipv6Regex, ip)

	return match
}

func CheckIPV6Subnet(ip string) bool {
	ipv6Regex := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))/\d{1,3}$`
	match, _ := regexp.MatchString(ipv6Regex, ip)

	return match
}

// ParseContentResult 解析nmap的XML文件
func (nmap *Nmap) ParseContentResult(content []byte) (result Result) {
	result.IPResult = make(map[string]*IPResult)
	nmapRunner, err := gonmap.Parse(content)
	if err != nil {
		//logging.RuntimeLog.Error(err)
		return
	}
	s := NewService()
	for _, host := range nmapRunner.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		var ip string
		for _, addr := range host.Addresses {
			ip = addr.Addr
			break
		}
		if ip == "" {
			continue
		}
		if !result.HasIP(ip) {
			result.SetIP(ip)
		}
		for _, port := range host.Ports {
			if port.State.State == "open" && port.Protocol == "tcp" {
				if !result.HasPort(ip, port.PortId) {
					result.SetPort(ip, port.PortId)
				}
				service := port.Service.Name
				if service == "" {
					service = s.FindService(port.PortId, ip)
				}
				result.SetPortAttr(ip, port.PortId, PortAttrResult{
					Source:  "portscan",
					Tag:     "service",
					Content: service,
				})
				banner := strings.Join([]string{port.Service.Product, port.Service.Version}, " ")
				if strings.TrimSpace(banner) != "" {
					result.SetPortAttr(ip, port.PortId, PortAttrResult{
						Source:  "portscan",
						Tag:     "banner",
						Content: banner,
					})
				}
			}
		}
	}
	return
}

// parseResult 解析nmap结果
func (nmap *Nmap) parseResult(outputTempFile string) {
	content, err := os.ReadFile(outputTempFile)
	if err != nil {
		//logging.RuntimeLog.Error(err)
		return
	}
	result := nmap.ParseContentResult(content)
	for ip, ipr := range result.IPResult {
		nmap.Result.IPResult[ip] = ipr
	}
}

// RunNmap 调用并执行一次nmap
func (nmap *Nmap) RunNmap(targets []string, ipv6 bool) {
	inputTargetFile := Utils.GetTempPathFileName()
	resultTempFile := Utils.GetTempPathFileName()
	defer os.Remove(inputTargetFile)
	defer os.Remove(resultTempFile)

	err := os.WriteFile(inputTargetFile, []byte(strings.Join(targets, "\n")), 0666)
	if err != nil {
		//logging.RuntimeLog.Error(err.Error())
		return
	}

	var cmdArgs []string
	cmdArgs = append(
		cmdArgs,
		nmap.Config.Tech, "-T4", "--open", "-n", "--randomize-hosts",
		"--min-rate", strconv.Itoa(nmap.Config.Rate), "-oX", resultTempFile, "-iL", inputTargetFile,
	)
	if ipv6 {
		cmdArgs = append(cmdArgs, "-6")
	}
	if !nmap.Config.IsPing {
		cmdArgs = append(cmdArgs, "-Pn")
	}
	if strings.HasPrefix(nmap.Config.Port, "--top-ports") {
		cmdArgs = append(cmdArgs, "--top-ports")
		cmdArgs = append(cmdArgs, strings.Split(nmap.Config.Port, " ")[1])
	} else {
		cmdArgs = append(cmdArgs, "-p", nmap.Config.Port)
	}
	if nmap.Config.ExcludeTarget != "" {
		cmdArgs = append(cmdArgs, "--exclude", nmap.Config.ExcludeTarget)
	}
	cmd := exec.Command(nmap.Config.CmdBin, cmdArgs...)
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		//logging.RuntimeLog.Error(err, stderr)
		//logging.CLILog.Error(err, stderr)
		return
	}
	nmap.parseResult(resultTempFile)
}

// Do 执行nmap
func (nmap *Nmap) nmapDo() {
	nmap.Result.IPResult = make(map[string]*IPResult)

	var targetIpV4, targetIpV6 []string
	for _, target := range strings.Split(nmap.Config.Target, ",") {
		t := strings.TrimSpace(target)
		if CheckIPV4(t) || CheckIPV4Subnet(t) {
			targetIpV4 = append(targetIpV4, t)
		} else if CheckIPV6(t) || CheckIPV6Subnet(t) {
			targetIpV6 = append(targetIpV6, t)
		}
	}
	if len(targetIpV4) > 0 {
		nmap.RunNmap(targetIpV4, false)
	}
	if len(targetIpV6) > 0 {
		nmap.RunNmap(targetIpV6, true)
	}
	FilterIPHasTooMuchPort(&nmap.Result, false)
}
