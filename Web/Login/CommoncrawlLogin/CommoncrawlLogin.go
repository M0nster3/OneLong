// Package commoncrawl logic
package CommoncrawlLogin

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func CommoncrawlLogin(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) string {
	//gologger.Infof("Commoncrawl  API 查询 \n")
	//gologger.Labelf("只实现普通Api 如果是企业修改Api接口 免费的每月250次\n")
	urls := "https://index.commoncrawl.org/collinfo.json"
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
		//client.SetProxy("192.168.203.111:1111")
	}
	client.Header = http.Header{
		"User-Agent": {Utils.RandUA()},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
	}

	client.Header.Set("Content-Type", "application/json")
	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(1 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)
	for add := 1; add < 4; add += 1 {
		if resp.RawResponse == nil || resp.StatusCode() == 503 {
			resp, _ = clientR.Get(urls)
			time.Sleep(1 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}
	if err != nil {
		gologger.Errorf("Commoncrawl API 链接访问失败尝试切换代理\n")
		return ""
	}
	buff := gjson.Parse(string(resp.Body())).Array()

	//addedURLs := make(map[string]bool)
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
	for aa, item := range buff {
		if aa == 4 {
			break
		}
		// 从当前条目获取域名
		cdx := item.Get("cdx-api").String()
		url := fmt.Sprintf("%s?url=*.%s/*&output=json&fl=url&page=0", cdx, domain)
		clientss := resty.New()
		clienta := clientss.R()
		clienta.Header = http.Header{
			"User-Agent": {Utils.RandUA()},
			"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
		}

		clienta.Header.Set("Content-Type", "application/json")
		clienta.Header.Del("Cookie")
		clienta.URL = url
		time.Sleep(1 * time.Second)
		respa, err := clienta.Get(url)
		if strings.Contains(string(respa.Body()), "No Captures found ") {

			//gologger.Labelf("Commoncrawl 后台查询未查询到域名 %s\n", domain)
			return ""

		}
		for add := 1; add < 20; add += 1 {
			if err != nil || respa.StatusCode() == 503 {
				clients := resty.New()
				clientaa := clients.R()
				clientaa.Header = http.Header{
					"User-Agent": {Utils.RandUA()},
					"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
					//"X-Key":      {options.ENConfig.Cookies.Binaryedge},
				}

				clientaa.Header.Set("Content-Type", "application/json")
				respa, _ = clientaa.Get(url)
				time.Sleep(10 * time.Second)
			} else if resp.Body() != nil {
				break
			}
		}
		if respa.StatusCode() == 503 {
			fmt.Printf("503")
			return ""
		}
		if respa.StatusCode() == 404 {
			return ""
		}

		//hostname := `(?:[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?\.)+` + regexp.QuoteMeta(domain)
		// 编译正则表达式
		//re := regexp.MustCompile(hostname)

		loginurls := strings.Split(string(respa.Body()), "\n")

		for _, loginurl := range loginurls {
			wg.Add(1)
			loginurl := loginurl
			go func() {
				addurl := gjson.Get(loginurl, "url").String()
				for content := range contentSet {
					if strings.Contains(addurl, content) {
						gologger.Infof("CommoncrawlLogin 匹配到后台链接:%s\n", addurl)
						//fmt.Println("CommoncrawlLogin 匹配到后台链接:", addurl)
						DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, addurl)
					}
				}
				wg.Done()
			}()

		}
		wg.Wait()

	}

	return "Success"
}
