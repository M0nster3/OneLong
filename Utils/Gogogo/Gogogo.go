package Gogogo

import (
	"OneLong/Api/Aiqicha"
	"OneLong/Api/tianyancha"
	"OneLong/Utils"
	outputfile "OneLong/Utils/OutPutfile"
	"OneLong/Utils/gologger"
	"sync"
)

// RunJob 运行项目 添加新参数记得去Config添加
func RunJob(options *Utils.ENOptions) {
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
				//defer func() {
				//	if x := recover(); x != nil {
				//		gologger.Errorf("[QCC] ERROR: %v", x)
				//		wg.Done()
				//	}
				//}()
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
						gologger.Errorf("[TYC] ERROR: %v", x)
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

	wg.Wait()
	//gologger.Infof("alienvault %s 域名，Address查询\n", options.KeyWord)
	//var OtherApiUrl []string
	//for k, s := range outputfile.EnsInfosList {
	//	if k == "icp" {
	//		for _, url := range s {
	//			OtherApiUrl = append(OtherApiUrl, url[2].(string))
	//		}
	//	} else {
	//		continue
	//	}
	//}
	//OtherApiUrl = Utils.SetStr(OtherApiUrl)
	//for _, url := range OtherApiUrl {
	//	res := alienvault.Alienvault(url, options)
	//	if res != "" {
	//		fmt.Printf(alienvault.AlienvaultResult[0][0] + "\n")
	//		fmt.Printf(alienvault.AlienvaultResult[1][0] + "\n")
	//	}
	//}

	// 如果不是API模式，而且不是批量文件形式查询 不是API 就合并导出到表格里面
	if options.IsMergeOut && options.InputFile == "" {
		outputfile.OutPutExcelByMergeEnInfo(options)
	}

}
