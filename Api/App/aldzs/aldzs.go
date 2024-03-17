package aldzs

import (
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getReq(searchType string, data map[string]string, options *Utils.ENOptions) gjson.Result {
	url := fmt.Sprintf("https://zhishuapi.aldwx.com/Main/action/%s", searchType)
	client := resty.New()
	client.SetTimeout(time.Duration(options.TimeOut) * time.Minute)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	client.Header = http.Header{
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"},
		"Accept":       {"text/html, application/xhtml+xml, image/jxr, */*"},
		"Content-Type": {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Referer":      {"https://www.aldzs.com"},
	}
	resp, err := client.R().SetFormData(data).Post(url)
	if err != nil {
		fmt.Println(err)
	}
	res := gjson.Parse(string(resp.Body()))
	if res.Get("code").String() != "200" {
		gologger.Errorf("【aldzs】似乎出了点问题 %s \n", res.Get("msg"))
	}
	return res.Get("data")
}

func GetInfoByKeyword(options *Utils.ENOptions) (ensInfos *Utils.EnInfos, ensOutMap map[string]*outputfile.ENSMap) {
	ensInfos = &Utils.EnInfos{}
	ensInfos.Infos = make(map[string][]gjson.Result)
	ensOutMap = make(map[string]*outputfile.ENSMap)

	keyword := options.KeyWord
	//拿到Token信息
	token := options.ENConfig.Cookies.Aldzs
	gologger.Infof("查询关键词 %s 的小程序\n", keyword)
	appList := getReq("Search/Search/search", map[string]string{
		"appName":    keyword,
		"page":       "1",
		"token":      token,
		"visit_type": "1",
	}, options).Array()
	if len(appList) == 0 {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NO", "ID", "小程序名称", "所属公司", "描述"})
	for k, v := range appList {
		table.Append([]string{
			strconv.Itoa(k),
			v.Get("id").String(),
			v.Get("name").String(),
			v.Get("company").String(),
			v.Get("desc").String(),
		})
	}
	table.Render()
	//默认取第一个进行查询
	var appKey string
	add := 0
	for bb, aa := range appList {
		if strings.Contains(aa.Get("company").String(), keyword) {
			gologger.Infof("查询 %s 开发的相关小程序 【默认取100个】\n", aa.Get("company"))
			appKey = aa.Get("appKey").String()
			sAppList := getReq("Miniapp/App/sameBodyAppList", map[string]string{
				"appKey": appKey,
				"page":   "1",
				"size":   "100",
				"token":  token,
			}, options).Array()
			//res[add] = append(res[add], sAppList...)
			ensInfos.Infos["wx_app"] = append(ensInfos.Infos["wx_app"], sAppList...)
			table = tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"NO", "ID", "小程序名称", "描述"})
			for k, v := range sAppList {
				table.Append([]string{
					strconv.Itoa(k),
					v.Get("id").String(),
					v.Get("name").String(),
					v.Get("desc").String(),
				})
			}
			table.Render()

			for k, v := range getENMap() {
				ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
			}
			add = 1
		} else if bb == len(appList)-1 && add == 0 {
			gologger.Infof("查询 %s 开发的相关小程序 【默认取100个】\n", appList[0].Get("company"))
			appKey = appList[0].Get("appKey").String()
			sAppList := getReq("Miniapp/App/sameBodyAppList", map[string]string{
				"appKey": appKey,
				"page":   "1",
				"size":   "100",
				"token":  token,
			}, options).Array()
			//res[add] = append(res[add], sAppList...)
			ensInfos.Infos["wx_app"] = sAppList
			table = tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"NO", "ID", "小程序名称", "描述"})
			for k, v := range sAppList {
				table.Append([]string{
					strconv.Itoa(k),
					v.Get("id").String(),
					v.Get("name").String(),
					v.Get("desc").String(),
				})
			}
			table.Render()

			for k, v := range getENMap() {
				ensOutMap[k] = &outputfile.ENSMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
			}
		} else {
			continue
		}
	}

	return ensInfos, ensOutMap
}

type EnsGo struct {
	name     string
	api      string
	fids     string
	params   map[string]string
	field    []string
	keyWord  []string
	typeInfo []string
}

func getENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"wx_app": {
			name:    "微信小程序",
			field:   []string{"name", "categoryTitle", "logo", "", ""},
			keyWord: []string{"名称", "分类", "头像", "二维码", "阅读量"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].keyWord = append(ensInfoMap[k].keyWord, "数据关联  ")
		ensInfoMap[k].field = append(ensInfoMap[k].field, "inFrom")
	}
	return ensInfoMap
}
