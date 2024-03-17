package Gogogo

import (
	"OneLong/Api/App/aldzs"
	"OneLong/Api/App/qimai"
	"OneLong/Api/Company/Aiqicha"
	"OneLong/Api/Company/tianyancha"
	"OneLong/Api/Domains"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"OneLong/Web/HttpZhiwen"
	"fmt"
	"strings"
	"sync"
)

// RunJob 运行项目 添加新参数记得去Config添加
func RunJob(options *Utils.ENOptions) {
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
			if options.ENConfig.Cookies.Tianyancha == "" || options.ENConfig.Cookies.Tycid == "" {
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
	for add, domain := range options.ICP {
		Domains.Domains(domain, options, &Domainip)
		fmt.Print(add)
	}

	var domain string
	for aa, _ := range options.ICP {
		sp := strings.Split(options.ICP[aa], ".")
		domain = domain + sp[len(sp)-2] + "." + sp[len(sp)-1] + " "
	}

	HttpZhiwen.Status(domain, options, &Domainip)
	// 如果不是API模式，而且不是批量文件形式查询 不是API 就合并导出到表格里面
	if options.IsMergeOut && options.InputFile == "" {
		outputfile.OutPutExcelByMergeEnInfo(options)
	}

}
