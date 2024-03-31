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
	"strings"
	"sync"

	//"strconv"
	//"strings"
	"time"
)

func WaybackarchiveLogin(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
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
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept":     {"text/html,application/json,application/xhtml+xml, image/jxr, */*"},
	}

	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
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
	if err != nil {
		gologger.Errorf("waybackarchive 历史快照链接访问失败尝试切换代理\n")
		return ""
	}
	if resp.Body() == nil {
		//gologger.Labelf("waybackarchive 历史快照未发现域名 %s\n", domain)
		return ""
	}

	loginurls := strings.Split(string(resp.Body()), "\n")
	last := "svg,css,eot,ttf,woff,jpg,png,jpeg,js,woff2,htm,gif,html"
	var mu sync.Mutex // 用于保护 addedURLs
	addedURLs := sync.Map{}
	for _, loginurl := range loginurls {
		wg.Add(1)
		loginurl := loginurl

		go func() {
			for content := range contentSet {
				if strings.Contains(loginurl, content) {
					wen := strings.Split(loginurl, "?")
					if len(wen) > 1 {
						lastThree := strings.Split(loginurl, ".")
						lastWen := strings.Split(lastThree[len(lastThree)-1], "?")
						if !strings.Contains(last, lastThree[len(lastThree)-1]) && !strings.Contains(last, lastWen[0]) {
							mu.Lock()
							if _, ok := addedURLs.LoadOrStore(wen[0], true); !ok {
								gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
								DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
							}
							mu.Unlock()

						}
					} else {
						lastThree := strings.Split(loginurl, ".")
						lastWen := strings.Split(lastThree[len(lastThree)-1], "?")
						if !strings.Contains(last, lastThree[len(lastThree)-1]) && !strings.Contains(last, lastWen[0]) {
							gologger.Infof("waybackarchive 匹配到链接:%s\n", loginurl)
							//fmt.Println("AlienvaultLogin 匹配到链接:", loginurl.String())
							DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl)
						}
					}
				}
			}
			wg.Done()
		}()

	}
	wg.Wait()
	//outputfile.OutPutExcelByMergeEnInfo(options)
	//
	//Result := gjson.GetMany(string(resp.Body()), "passive_dns.#.address", "passive_dns.#.hostname")
	//AlienvaultResult[0] = append(AlienvaultResult[0], Result[0].String())
	//AlienvaultResult[1] = append(AlienvaultResult[1], Result[1].String())
	//
	//fmt.Printf(Result[0].String())
	return "Success"
}
