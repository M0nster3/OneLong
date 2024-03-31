package outputfile

import (
	"OneLong/Utils"
	"OneLong/Utils/gologger"
	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type ENSMap struct {
	Name    string
	Field   []string
	JField  []string
	KeyWord []string
	Only    string
}
type DomainsIP struct {
	Domains     []string
	DomainA     []string //DomainA 指纹识别里面需要的
	IP          []string
	IPA         []string //ipa 指纹识别里面需要的
	TitleBUff   []string
	Size        []string
	Status_code []string
	Zhiwen      []string
	LoginUrl    []string
	LoginUrlA   []string
	LoginTitle  []string
	A           []string
}

var EnsInfosList = make(map[string][][]interface{})
var ENSMapList = make(map[string]*ENSMap)

var ENSMapLN = map[string]*ENSMap{
	"enterprise_info": {
		Name:    "企业信息",
		JField:  []string{"name", "legal_person", "status", "phone", "email", "registered_capital", "incorporation_date", "address", "scope", "reg_code", "pid"},
		KeyWord: []string{"企业名称", "法人代表", "经营状态", "电话", "邮箱", "注册资本", "成立日期", "注册地址", "经营范围", "统一社会信用代码", "PID"},
	},
	"icp": {
		Name:    "ICP信息",
		Only:    "domain",
		JField:  []string{"wesbite_name", "website", "domain", "icp", "company_name"},
		KeyWord: []string{"网站名称", "网址", "域名", "网站备案/许可证号", "公司名称"},
	},
	"wx_app": {
		Name:    "微信小程序",
		JField:  []string{"name", "category", "logo", "qrcode", "read_num"},
		KeyWord: []string{"名称", "分类", "头像", "二维码", "阅读量"},
	},
	"wechat": {
		Name:    "微信公众号",
		JField:  []string{"name", "wechat_id", "description", "qrcode", "avatar"},
		KeyWord: []string{"名称", "ID", "简介", "二维码", "头像"},
	},
	"weibo": {
		Name:    "微博",
		JField:  []string{"name", "profile_url", "description", "avatar"},
		KeyWord: []string{"微博昵称", "链接", "简介", "头像"},
	},
	"supplier": {
		Name:    "供应商",
		JField:  []string{"name", "scale", "amount", "report_time", "data_source", "relation", "pid"},
		KeyWord: []string{"名称", "金额占比", "金额", "报告期/公开时间", "数据来源", "关联关系", "PID"},
	},
	"job": {
		Name:    "招聘",
		JField:  []string{"name", "education", "location", "publish_time", "salary"},
		KeyWord: []string{"招聘职位", "学历", "办公地点", "发布日期", "薪资"},
	},
	"invest": {
		Name:    "投资",
		JField:  []string{"name", "legal_person", "status", "scale", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "投资比例", "PID"},
	},
	"branch": {
		Name:    "分支机构",
		JField:  []string{"name", "legal_person", "status", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "PID"},
	},
	"holds": {
		Name:    "控股企业",
		JField:  []string{"name", "legal_person", "status", "scale", "level", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "投资比例", "持股层级", "PID"},
	},
	"app": {
		Name:    "应用",
		JField:  []string{"name", "category", "version", "update_at", "description", "logo", "bundle_id", "link", "market"},
		KeyWord: []string{"名称", "分类", "当前版本", "更新时间", "简介", "logo", "Bundle ID", "链接", "market"},
	},
	"copyright": {
		Name:    "软件著作权",
		JField:  []string{"name", "short_name", "category", "reg_num", "pub_type"},
		KeyWord: []string{"软件全称", "软件简称", "分类", "登记号", "权利取得方式"},
	},
	"partner": {
		Name:    "股东信息",
		JField:  []string{"name", "scale", "reg_cap", "pid"},
		KeyWord: []string{"股东名称", "持股比例", "认缴出资金额", "PID"},
	},
}
var EnJsonList = make(map[string][]map[string]interface{})

// MergeOutPut 数据合并到MAP
func MergeOutPut(ensInfos *Utils.EnInfos, ensMap map[string]*ENSMap, info string, options *Utils.ENOptions) map[string][][]interface{} {
	if options.Output != "!" {
		//gologger.Infof("%s【%s】信息合并\n", info, ensInfos.Name)
		for k, s := range ensInfos.Infos {
			ENSMapList[k] = ensMap[k]
			var data [][]interface{}
			if k == "icp" {
				for _, aa := range s {
					domain := gjson.Get(aa.String(), "domain").String()
					options.ICP = append(options.ICP, domain)
				}
			}
			for _, y := range s {
				results := gjson.GetMany(y.Raw, ensMap[k].Field...)
				var str []interface{}
				for _, t := range results {
					str = append(str, t.String())
				}
				//if !options.IsApiMode {
				//	str = append(str, info+"【"+ensInfos.Name+"】")
				//} else {
				//	str = append(str, info)
				//}
				str = append(str, info)
				data = append(data, str)
			}
			EnsInfosList[k] = append(EnsInfosList[k], data...)
		}
	}
	return EnsInfosList

}

// OutPutExcelByMergeEnInfo 根据合并信息导出表格
func OutPutExcelByMergeEnInfo(options *Utils.ENOptions) {
	if options.Output != "!" {
		tmp := options.Output
		_, err := os.Stat(tmp)
		if err != nil {
			gologger.Infof("【%s】目录不存在，自动创建\n", tmp)
			err := os.Mkdir(tmp, os.ModePerm)
			if err != nil {
				gologger.Fatalf("缺少%s文件夹，并且创建失败！", tmp)
			}
		}
		// Save spreadsheet by the given path.
		fileName := ""
		if len([]rune(options.CompanyName)) > 20 {
			fileName = options.KeyWord
		} else if options.KeyWord != "" {
			fileName = options.CompanyName
		} else {
			options.Domain = strings.ReplaceAll(options.Domain, "http://", "")
			options.Domain = strings.ReplaceAll(options.Domain, "https://", "")
			options.Domain = strings.ReplaceAll(options.Domain, "/", "")
			fileName = options.Domain
		}
		savaPath := tmp + "/" + fileName + "--" + time.Now().Format("2006-01-02") + "--" + strconv.FormatInt(time.Now().Unix(), 10)

		savaPath += ".xlsx"
		f := excelize.NewFile()
		if options.Domain == "" {
			gologger.Infof("【%s】导出中\n", options.CompanyName)
		} else {

			gologger.Infof("【%s】导出中\n", options.Domain)
		}

		for k, s := range EnsInfosList {
			if _, ok := ENSMapList[k]; ok {
				gologger.Infof("正在导出%s\n", ENSMapList[k].Name)
				headers := ENSMapList[k].KeyWord
				headers = append(headers, "查询信息")
				data := s
				f, _ = Utils.ExportExcel(ENSMapList[k].Name, headers, data, f)
			} else {
				gologger.Errorf("导出错误信息 %s\n", k)
			}
		}

		f.DeleteSheet("Sheet1")

		if err := f.SaveAs(savaPath); err != nil {
			gologger.Fatalf("导出失败：%s", err)
		}

		gologger.Infof("导出成功路径： %s\n", savaPath)
		EnJsonList = make(map[string][]map[string]interface{})
		EnsInfosList = make(map[string][][]interface{})
		ENSMapList = make(map[string]*ENSMap)
	}
}

// OutPutExcelByEnInfo 直接导出单独表格信息
func OutPutExcelByEnInfo(ensInfos *Utils.EnInfos, ensMap map[string]*ENSMap, options *Utils.ENOptions) {
	if options.Output != "!" {
		if ensInfos.Name == "" {
			ensInfos.Name = options.KeyWord
		}
		if ensInfos.Name != "" {
			//初始化导出目录
			tmp := options.Output
			_, err := os.Stat(tmp)
			if err != nil {
				gologger.Infof("【%s】目录不存在，自动创建\n", tmp)
				err := os.Mkdir(tmp, os.ModePerm)
				if err != nil {
					gologger.Fatalf("缺少%s文件夹，并且创建失败！", tmp)
				}
			}
			// 修复导出文件名过长的问题
			fileName := ""
			if len([]rune(ensInfos.Name)) > 20 {
				fileName = options.KeyWord
			} else {
				fileName = ensInfos.Name
			}
			savaPath := tmp + "/" + fileName + "--" + time.Now().Format("2006-01-02") + "--" + strconv.FormatInt(time.Now().Unix(), 10)

			savaPath += ".xlsx"
			// 导出表格信息
			f := excelize.NewFile()
			gologger.Infof("【%s】导出中\n", ensInfos.Name)
			for k, s := range ensInfos.Infos {
				gologger.Infof("正在导出%s\n", ensMap[k].Name)
				headers := ensMap[k].KeyWord
				var data [][]interface{}
				for _, y := range s {
					var str []interface{}
					results := gjson.GetMany(y.Raw, ensMap[k].Field...)
					for _, t := range results {
						str = append(str, t.String())
					}
					data = append(data, str)
				}
				f, _ = Utils.ExportExcel(ensMap[k].Name, headers, data, f)
			}
			f.DeleteSheet("Sheet1")
			if err := f.SaveAs(savaPath); err != nil {
				gologger.Fatalf("表格导出失败：%s", err)
			}

			gologger.Infof("导出成功路径： %s\n", savaPath)
		} else {
			gologger.Errorf("无法导出，公司名不存在")
		}
	}
}
