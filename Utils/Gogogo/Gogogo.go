package Gogogo

import (
	"OneLong/Api/App/aldzs"
	"OneLong/Api/App/qimai"
	"OneLong/Api/Company/Aiqicha"
	"OneLong/Api/Company/tianyancha"
	"OneLong/Api/Domains"
	"OneLong/Email"
	"OneLong/Script/Afrog"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/HttpZhiwen"
	"OneLong/Web/Login"
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

// FileScan 可批量导入文件查询
func StartScan(options *Utils.LongOptions) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)

	// 在另一个 goroutine 中等待信号
	go func() {
		// 等待信号
		<-sig

		outputfile.OutPutExcelByMergeEnInfo(options)

		os.Exit(0) // 可以执行一些清理工作后退出程序
	}()
	if options.InputFile != "" {
		res := Utils.ReadFile(options.InputFile)

		time.Sleep(5 * time.Second)
		res = Utils.SetStr(res)
		color.RGBStyleFromString("237,64,35").Printf(fmt.Sprintf("--------------------批量查询%d条域名资产--------------------\n", len(res)))

		for id, v := range res {
			color.RGBStyleFromString("237,64,35").Printf(fmt.Sprintf("当前批量查询第 %d 条域名,还剩 %d 条\n", id+1, len(res)-id-1))
			if v == "" {
				continue
			}
			//color.RGBStyleFromString("244,211,49").Printf(fmt.Sprintf("--------------------【第%d条】 %s 查询中--------------------\n", k+1, v))
			if !strings.Contains(v, ".") {
				options.CompanyID = ""
				options.KeyWord = v
				CompanyRunJob(options)
			} else {
				options.Domain = v
				DomainRunJob(options)
			}
		}

		outputfile.OutPutExcelByMergeEnInfo(options)

	} else {
		if !strings.Contains(options.Domain, ".") {
			CompanyRunJob(options)
		} else {
			DomainRunJob(options)
		}

		outputfile.OutPutExcelByMergeEnInfo(options)

	}
}

// CompanyRunJob 运行项目
func CompanyRunJob(options *Utils.LongOptions) {

	color.RGBStyleFromString("244,211,49").Println("\n--------------------查询企业信息--------------------")
	var Domainip outputfile.DomainsIP
	if options.Proxy != "" {
		gologger.Infof("代理地址: %s\n", options.Proxy)
	}
	gologger.Infof("关键词:【%s|%s】数据源：%s 数据字段：%s\n", options.KeyWord, options.CompanyID, options.GetType, options.GetField)
	var wg sync.WaitGroup

	//爱企查
	if Utils.IsInList("aqc", options.GetType) {
		if options.CompanyID == "" || (options.CompanyID != "" && Utils.CheckPid(options.CompanyID) == "aqc") {
			wg.Add(1)
			go func() {
				//查询企业信息
				res, ensOutMap := Aiqicha.GetEnInfoByPid(options)
				if options.IsMergeOut {
					//合并导出
					outputfile.MergeOutPut(res, ensOutMap, "爱企查", options)
				} else {
					//单独导出
					outputfile.OutPutExcelByEnInfo(res, ensOutMap, options)
				}
				//hook.BiuScan(res, options)
				wg.Done()
			}()
		}
	}
	//
	////天眼查
	if Utils.IsInList("tyc", options.GetType) {
		if options.CompanyID == "" || (options.CompanyID != "" && Utils.CheckPid(options.CompanyID) == "tyc") {
			wg.Add(1)
			if options.LongConfig.Cookies.Tianyancha == "" || options.LongConfig.Cookies.Tycid == "" {
				gologger.Fatalf("【TYC】MUST LOGIN 请在配置文件补充天眼查COOKIE和tycId\n")
			}
			go func() {
				defer func() {
					if x := recover(); x != nil {
						gologger.Errorf("[TYC] ERROR: %v\n", x)
						wg.Done()
					}
				}()
				res, ensOutMap := tianyancha.GetEnInfoByPid(options)
				if options.IsMergeOut {
					outputfile.MergeOutPut(res, ensOutMap, "天眼查", options)
				} else {
					outputfile.OutPutExcelByEnInfo(res, ensOutMap, options)
				}
				//hook.BiuScan(res, options)
				wg.Done()
			}()
		}
	}

	// 微信小程序查询
	if Utils.IsInList("aldzs", options.GetType) {
		wg.Add(1)
		res, ensOutMap := aldzs.GetInfoByKeyword(options)
		if options.IsMergeOut {
			outputfile.MergeOutPut(res, ensOutMap, "阿拉丁指数", options)
		} else {
			outputfile.OutPutExcelByEnInfo(res, ensOutMap, options)
		}
		wg.Done()
	}
	// 七麦数据
	if Utils.IsInList("qimai", options.GetType) {
		wg.Add(1)
		go func() {
			res, ensOutMap := qimai.GetInfoByKeyword(options)
			outputfile.MergeOutPut(res, ensOutMap, "七麦数据", options)
			//outputfile.OutPutExcelByEnInfo(res, ensOutMap, options)
			wg.Done()
		}()
	}
	wg.Wait()

	options.ICP = Utils.SetStr(options.ICP)
	if len(options.ICP) == 0 {
		color.RGBStyleFromString("237,64,35").Println(fmt.Sprintf("当前 %s 字段没有相关备案域名", options.KeyWord))
		outputfile.OutPutExcelByMergeEnInfo(options)
		os.Exit(1)
	}
	color.RGBStyleFromString("244,211,49").Println("\n--------------------查询子域名--------------------")
	for _, domain := range options.ICP {
		Domains.Domains(domain, options, &Domainip)
	}

	var domain string
	for aa, _ := range options.ICP {
		sp := strings.Split(options.ICP[aa], ".")
		domain = domain + sp[len(sp)-2] + "." + sp[len(sp)-1] + " "
	}
	//Domains.Domains(domain, options, &Domainip)
	//color.RGBStyleFromString("244,211,49").Println("\n--------------------整合域名、IP、指纹--------------------")
	HttpZhiwen.Status(domain, options, &Domainip)
	color.RGBStyleFromString("244,211,49").Println("\n--------------------探测网站后台--------------------")
	Login.Login(Domainip.DomainA, options, &Domainip)
	// 如果不是API模式，而且不是批量文件形式查询 不是API 就合并导出到表格里面
	if !options.NoPoc {
		color.RGBStyleFromString("244,211,49").Println("\n--------------------漏洞扫描--------------------")
		Afrog.Afrog(options, &Domainip)
	}
	color.RGBStyleFromString("244,211,49").Println("\n--------------------探测邮箱--------------------")
	Email.Email(options.Domain, options, &Domainip)

	//if options.IsMergeOut && options.InputFile == "" {
	//	outputfile.OutPutExcelByMergeEnInfo(options)
	//}
}
func DomainRunJob(options *Utils.LongOptions) {

	options.Domain = strings.ReplaceAll(options.Domain, "http://", "")
	options.Domain = strings.ReplaceAll(options.Domain, "https://", "")
	var Domainip outputfile.DomainsIP
	if options.Proxy != "" {
		gologger.Infof("代理地址: %s\n", options.Proxy)
	}
	gologger.Infof("当前查询的域名 %s", options.Domain)
	//color.RGBStyleFromString("237,64,35").Println("查询子域名\n")
	re := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)
	matches := re.FindAllStringSubmatch(options.Domain, -1)
	if matches != nil {
		color.RGBStyleFromString("244,211,49").Println("\n当前不支持IP查询")
		return
	}

	color.RGBStyleFromString("244,211,49").Println("\n--------------------查询子域名--------------------")
	Domains.Domains(options.Domain, options, &Domainip)
	HttpZhiwen.Status(options.Domain, options, &Domainip)

	color.RGBStyleFromString("244,211,49").Println("\n--------------------探测邮箱--------------------")
	Email.Email(options.Domain, options, &Domainip)

	if !options.NoPoc {
		color.RGBStyleFromString("244,211,49").Println("\n--------------------漏洞扫描--------------------")
		Afrog.Afrog(options, &Domainip)
	}
	color.RGBStyleFromString("244,211,49").Println("\n--------------------探测网站后台--------------------")
	Login.Login(Domainip.DomainA, options, &Domainip)

}
