package finger

import (
	"OneLong/Script/Ehole/module/queue"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Web/CDN"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Outrestul struct {
	Url        string `json:"url"`
	Cms        string `json:"cms"`
	Server     string `json:"server"`
	Statuscode int    `json:"statuscode"`
	Length     int    `json:"length"`
	Title      string `json:"title"`
}

type FinScan struct {
	UrlQueue    *queue.Queue
	Ch          chan []string
	Wg          sync.WaitGroup
	Thread      int
	Proxy       string
	AllResult   []Outrestul
	FocusResult []Outrestul
	Finpx       *Packjson
}

func NewScan(urls []string, thread int, proxy string) *FinScan {
	s := &FinScan{
		UrlQueue:    queue.NewQueue(),
		Ch:          make(chan []string, thread),
		Wg:          sync.WaitGroup{},
		Thread:      thread,
		Proxy:       proxy,
		AllResult:   []Outrestul{},
		FocusResult: []Outrestul{},
	}
	err := LoadWebfingerprint(Utils.GetPathDir() + "/Ehole.json")
	if err != nil {
		color.RGBStyleFromString("237,64,35").Println("[error] fingerprint file error!!!")
		os.Exit(1)
	}
	s.Finpx = GetWebfingerprint()
	for _, url := range urls {
		s.UrlQueue.Push([]string{url, "0"})
	}
	return s
}

func (s *FinScan) StartScan(DomainsIP *outputfile.DomainsIP) {
	for i := 0; i <= s.Thread; i++ {
		s.Wg.Add(1)
		go func() {
			defer s.Wg.Done()
			s.fingerScan(DomainsIP)
		}()
	}
	s.Wg.Wait()
	//color.RGBStyleFromString("244,211,49").Println("\n重点资产：")
	//for _, aas := range s.FocusResult {
	//	fmt.Printf(fmt.Sprintf("[ %s | ", aas.Url))
	//	color.RGBStyleFromString("237,64,35").Printf(fmt.Sprintf("%s", aas.Cms))
	//	fmt.Printf(fmt.Sprintf(" | %s | %d | %d | %s ]\n", aas.Server, aas.Statuscode, aas.Length, aas.Title))
	//}

}

func MapToJson(param map[string][]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func (s *FinScan) fingerScan(DomainsIP *outputfile.DomainsIP) {
	for s.UrlQueue.Len() != 0 {
		dataface := s.UrlQueue.Pop()
		switch dataface.(type) {
		case []string:
			url := dataface.([]string)
			var data *resps
			data, err := httprequest(url, s.Proxy)
			if err != nil {
				url[0] = strings.ReplaceAll(url[0], "https://", "http://")
				data, err = httprequest(url, s.Proxy)
				if err != nil {
					continue
				}
			}
			for _, jurl := range data.jsurl {
				if jurl != "" {
					s.UrlQueue.Push([]string{jurl, "1"})
				}
			}
			headers := MapToJson(data.header)
			var cms []string
			for _, finp := range s.Finpx.Fingerprint {
				if finp.Location == "body" {
					if finp.Method == "keyword" {
						if iskeyword(data.body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "faviconhash" {
						if data.favhash == finp.Keyword[0] {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(data.body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "header" {
					if finp.Method == "keyword" {
						if iskeyword(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "title" {
					if finp.Method == "keyword" {
						if iskeyword(data.title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(data.title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
			}
			cms = RemoveDuplicatesAndEmpty(cms)
			cmss := strings.Join(cms, ",")
			out := Outrestul{data.url, cmss, data.server, data.statuscode, data.length, data.title}
			s.AllResult = append(s.AllResult, out)
			var ip string
			if len(out.Cms) != 0 {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				color.RGBStyleFromString("237,64,35").Println(outstr)
				s.FocusResult = append(s.FocusResult, out)
				zhiwen := out.Cms + " , " + out.Server + "[可能存在漏洞]"
				reurl := strings.ReplaceAll(out.Url, "https://", "")
				reurl1 := strings.ReplaceAll(reurl, "http://", "")
				// 编译正则表达式
				hostname := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
				re := regexp.MustCompile(hostname)

				matches := re.FindAllStringSubmatch(strings.TrimSpace(reurl1), -1)
				ips, _ := net.LookupIP(matches[0][0])

				for _, aa := range ips {
					boolcdn, _, _ := CDN.CheckCName(aa.String())
					if boolcdn {
						ip = "CDN"
						break
					}

					if len(ips) == 1 {
						ip = aa.String()
					} else {
						ip = ip + aa.String() + " , "
					}

				}
				if out.Statuscode != 502 {
					DomainsIP.Zhiwen = append(DomainsIP.Zhiwen, zhiwen)
					DomainsIP.A = append(DomainsIP.A, ip)
					DomainsIP.DomainA = append(DomainsIP.DomainA, out.Url)
					DomainsIP.Status_code = append(DomainsIP.Status_code, strconv.Itoa(out.Statuscode))
					DomainsIP.TitleBUff = append(DomainsIP.TitleBUff, out.Title)
					DomainsIP.Size = append(DomainsIP.Size, strconv.Itoa(out.Length))
				}

			} else {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				fmt.Println(outstr)
				var zhiwen string
				if out.Cms != "" {
					zhiwen = out.Cms + " , " + out.Server
				} else {
					zhiwen = out.Server
				}
				reurl := strings.ReplaceAll(out.Url, "https://", "")
				reurl1 := strings.ReplaceAll(reurl, "http://", "")
				hostname := `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
				// 编译正则表达式
				re := regexp.MustCompile(hostname)

				matches := re.FindAllStringSubmatch(strings.TrimSpace(reurl1), -1)
				ips, _ := net.LookupIP(matches[0][0])

				for _, aa := range ips {
					boolcdn, _, _ := CDN.CheckCName(aa.String())
					if boolcdn {
						ip = "CDN"
						break
					}

					if len(ips) == 1 {
						ip = aa.String()
					} else {
						ip = ip + aa.String() + " , "
					}

				}

				if out.Statuscode != 502 {
					DomainsIP.A = append(DomainsIP.A, ip)
					DomainsIP.Zhiwen = append(DomainsIP.Zhiwen, zhiwen)
					DomainsIP.DomainA = append(DomainsIP.DomainA, out.Url)
					DomainsIP.Status_code = append(DomainsIP.Status_code, strconv.Itoa(out.Statuscode))
					DomainsIP.TitleBUff = append(DomainsIP.TitleBUff, out.Title)
					DomainsIP.Size = append(DomainsIP.Size, strconv.Itoa(out.Length))
				}
			}
		default:
			continue
		}
	}
}
