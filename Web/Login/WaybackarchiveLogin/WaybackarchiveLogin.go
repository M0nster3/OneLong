package WaybackarchiveLogin

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	//"strconv"
	//"strings"
	"time"
)

func WaybackarchiveLogin(domain string, options *Utils.LongOptions, DomainsIP *outputfile.DomainsIP, login int) string {
	//gologger.Infof("waybackarchive 历史快照查询\n")
	var wg sync.WaitGroup
	dir := filepath.Join(Utils.GetPathDir(), "Script/Dict/Login.txt")
	file, err := os.Open(dir)
	if err != nil {
		gologger.Errorf("无法打开文件后台目录文件%s\n", dir)
		return ""
	}
	defer file.Close()

	// 使用哈希集合存储文本中的内容
	contentSet := make(map[string]bool)

	// 创建 Scanner 对象
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		contentSet[line] = true
	}

	urls := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original&collapse=urlkey", domain)
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":          {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
		"Accept-Encoding": {"gzip"},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(3 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
	if err != nil {
		for attempt := 0; attempt < 4; attempt++ {
			time.Sleep(5 * time.Second) // 在重试前等待
			resp, err = client.R().Get(urls)
			if err != nil || resp == nil || resp.RawResponse == nil {
				time.Sleep(5 * time.Second) // 在重试前等待
				continue
			}
			// 如果得到有效响应，处理响应
			if resp.Body() != nil {
				// 处理响应的逻辑
				break
			}
		}
	}

	if err != nil && login == 1 {
		gologger.Errorf("waybackarchive 历史快照链接访问失败尝试切换代理\n")
		return ""
	}
	if err != nil || resp.Body() == nil {
		return ""
	}
	//if resp.Body() == nil {
	//	//gologger.Labelf("waybackarchive 历史快照未发现域名 %s\n", domain)
	//	return ""
	//}

	loginurls := strings.Split(string(resp.Body()), "\n")
	//last, err := os.ReadFile(filepath.Join(Utils.GetPathDir(), "Script/Dict/ExcludeLogin.txt"))
	//if err != nil {
	//	gologger.Errorf("Alienvault API 读取 ExcludeLogin 文件失败\n")
	//	return ""
	//}
	//last := "svg,css,eot,ttf,woff,jpg,png,jpeg,js,woff2,htm,gif,html,xml,swf"
	var mu sync.Mutex // 用于保护 addedURLs
	//addedURLs := sync.Map{}
	LoginUrls := make(map[string]bool)
	WenUrls := make(map[string]bool)
OuterLoop:
	for _, loginurl := range loginurls {
		loginurl := loginurl
		host := regexp.MustCompile(`(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`)
		url := host.FindAllStringSubmatch(strings.TrimSpace(loginurl), -1)
		if url != nil && len(url) > 0 {
			for content := range contentSet {
				if strings.Contains(url[0][0], content) {
					if !LoginUrls[url[0][0]] {
						gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
						DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
						LoginUrls[url[0][0]] = true
					}
					continue OuterLoop
				}
			}

			wg.Add(1)
			go func() {
				for content := range contentSet {
					lowurl := strings.ToLower(loginurl)
					index := strings.Index(lowurl, content)
					if index != -1 {
						// 截取到之前的内容和 content本身
						//result := loginurl[:index+len(content)]
						Hurl := loginurl
						protocol := ""
						if strings.HasPrefix(Hurl, "http://") {
							protocol = "http://"
							Hurl = Hurl[len("http://"):]
						} else if strings.HasPrefix(Hurl, "https://") {
							protocol = "https://"
							Hurl = Hurl[len("https://"):]
						}

						parts := strings.Split(Hurl, "/")

						if len(parts) > 1 {
							Hurl = protocol + parts[0] + "/" + parts[1]
						} else {
							Hurl = protocol + Hurl
						}
						mu.Lock()
						if !LoginUrls[Hurl] {
							wen := strings.Split(loginurl, "?")
							if len(wen) != 0 {
								if !WenUrls[wen[0]] {
									gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
									DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
									//LoginUrls[result] = true
									LoginUrls[Hurl] = true
									WenUrls[wen[0]] = true
								}
							} else {
								gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
								DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
								//LoginUrls[result] = true
								LoginUrls[Hurl] = true
							}

						}
						mu.Unlock()
					}
					//if strings.Contains(loginurl, content) {
					//	//wen := strings.Split(loginurl, "?")
					//	//if len(wen) > 1 {
					//	//	lastThree := strings.Split(wen[0], ".")
					//	//	lastWen := strings.Split(lastThree[len(lastThree)-1], "?")
					//	//	if !strings.Contains(string(last), strings.ToLower(lastThree[len(lastThree)-1])) && !strings.Contains(string(last), lastWen[0]) {
					//	//		mu.Lock()
					//	//		if _, ok := addedURLs.LoadOrStore(wen[0], true); !ok {
					//	//			index := strings.Index(loginurl, content)
					//	//			if index != -1 {
					//	//				// 截取到之前的内容和 content本身
					//	//				result := loginurl[:index+len(content)]
					//	//				if !LoginUrls[result] {
					//	//					gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
					//	//					DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
					//	//					LoginUrls[result] = true
					//	//				}
					//	//			}
					//	//
					//	//		}
					//	//		mu.Unlock()
					//	//
					//	//	}
					//	//} else {
					//	//	lastThree := strings.Split(loginurl, ".")
					//	//	lastWen := strings.Split(lastThree[len(lastThree)-1], "?")
					//	//	if !strings.Contains(string(last), strings.ToLower(lastThree[len(lastThree)-1])) && !strings.Contains(string(last), lastWen[0]) {
					//	//		gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
					//	//		//fmt.Println("AlienvaultLogin 匹配到链接:", loginurl.String())
					//	//		DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
					//	//	}
					//	//}
					//
					//}
				}
				wg.Done()
			}()
		}

	}
	wg.Wait()

	return "Success"
}
