package alienvaultLogin

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

	//"strconv"
	//"strings"
	"time"
)

var wg sync.WaitGroup

func AlienvaultLogin(domain string, options *Utils.ENOptions, DomainsIP *outputfile.DomainsIP) {

	dir := filepath.Join(Utils.GetPathDir(), "Script/Dict/Login.txt")
	file, err := os.Open(dir)
	if err != nil {
		gologger.Errorf("无法打开文件后台目录文件%s\n", dir)
		return
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
	urls := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list?limit=100&page=1", domain)
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

	client.Header.Set("Content-Type", "application/json")
	client.Header.Del("Cookie")

	//强制延时1s
	time.Sleep(3 * time.Second)
	//加入随机延迟
	time.Sleep(time.Duration(options.GetDelayRTime()) * time.Second)
	clientR := client.R()

	clientR.URL = urls
	resp, err := clientR.Get(urls)

	for add := 1; add < 4; add += 1 {
		if resp.RawResponse == nil {
			resp, _ = clientR.Get(urls)
			time.Sleep(4 * time.Second)
		} else if resp.Body() != nil {
			break
		}
	}

	if err != nil {
		gologger.Errorf("Alienvault API 链接访问失败尝试切换代理\n")
		return
	}

	count := gjson.GetBytes(resp.Body(), "full_size").Int()
	if count == 0 {
		//gologger.Labelf("Alienvault 未发现后台域名 %s\n", domain)
		return
	}
	intcount := int(count)
	for add := 1; add <= intcount/100+1; add += 1 {
		urls = fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list?limit=100&page=%d", domain, add)
		resp, err = clientR.Get(urls)
		loginurls := gjson.GetBytes(resp.Body(), "url_list.#.url").Array()

		for _, loginurl := range loginurls {
			wg.Add(1)
			loginurl := loginurl
			go func() {
				for content := range contentSet {
					if strings.Contains(loginurl.String(), content) {
						gologger.Infof("AlienvaultLogin 匹配到链接:%s\n", loginurl.String())
						//fmt.Println("AlienvaultLogin 匹配到链接:", loginurl.String())
						DomainsIP.LoginUrl = append(DomainsIP.LoginUrl, loginurl.String())
					}
				}
				wg.Done()
			}()

		}
		wg.Wait()

	}

}
