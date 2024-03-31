package Utils

import (
	"OneLong/Utils/gologger"
	"flag"
	"github.com/gookit/color"
)

const banner = `


 ██████╗ ███╗   ██╗███████╗██╗      ██████╗ ███╗   ██╗ ██████╗ 
██╔═══██╗████╗  ██║██╔════╝██║     ██╔═══██╗████╗  ██║██╔════╝ 
██║   ██║██╔██╗ ██║█████╗  ██║     ██║   ██║██╔██╗ ██║██║  ███╗
██║   ██║██║╚██╗██║██╔══╝  ██║     ██║   ██║██║╚██╗██║██║   ██║
╚██████╔╝██║ ╚████║███████╗███████╗╚██████╔╝██║ ╚████║╚██████╔╝
 ╚═════╝ ╚═╝  ╚═══╝╚══════╝╚══════╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝
`

func Banner() {
	gologger.Printf("%s\n\n", banner)
	gologger.Printf("\t\thttps://github.com/M0nster3/OneLong\n\n")
	//gologger.Labelf("工具仅用于信息收集，请勿用于非法用途\n")
	color.RGBStyleFromString("237,64,35").Println("工具仅用于信息收集，请勿用于非法用途\n开发人员不承担任何责任，也不对任何滥用或损坏负责.\n")
	color.RGBStyleFromString("244,211,49").Println("使用方式: \n\tOneLong -n 企业名称\n\tOneLong -d target.com\n")
}

func Flag(Info *ENOptions) {
	Banner()
	flag.BoolVar(&Info.NoBao, "nb", false, "不进行爆破子域名")
	flag.BoolVar(&Info.NoPoc, "np", false, "不进行漏洞扫描")
	flag.StringVar(&Info.KeyWord, "n", "", "关键词 eg 小米")
	flag.StringVar(&Info.Domain, "d", "", "域名")
	//flag.StringVar(&Info.CompanyID, "i", "", "公司PID")
	//flag.StringVar(&Info.InputFile, "f", "", "批量查询，文本按行分隔")
	//flag.StringVar(&Info.ScanType, "type", "aqc", "API类型 eg qcc")
	flag.StringVar(&Info.Output, "o", "", "结果输出的文件夹位置(可选)")
	//flag.BoolVar(&Info.IsMergeOut, "is-merge", false, "合并导出")
	//flag.BoolVar(&Info.IsJsonOutput, "json", false, "json导出")
	//查询参数指定
	flag.Float64Var(&Info.InvestNum, "invest", 70, "投资比例 ")
	//flag.StringVar(&Info.GetFlags, "field", "", "获取字段信息 eg icp")
	flag.IntVar(&Info.Deep, "deep", 5, "递归搜索n层公司")
	flag.BoolVar(&Info.IsHold, "hold", true, "是否查询控股公司，默认查询")
	flag.BoolVar(&Info.IsSupplier, "supplier", true, "是否查询供应商信息,默认查询")
	flag.BoolVar(&Info.IsGetBranch, "branch", true, "查询分支机构（分公司）信息 默认查询")
	flag.BoolVar(&Info.IsSearchBranch, "is-branch", false, "深度查询分支机构信息（数量巨大），默认不查询")
	//web api
	//flag.BoolVar(&Info.IsApiMode, "api", false, "是否API模式")
	//flag.StringVar(&Info.ClientMode, "client", "", "客户端模式通道 eg: task")
	//flag.BoolVar(&Info.IsDebug, "debug", false, "是否显示debug详细信息")
	//flag.BoolVar(&Info.IsShow, "is-show", true, "是否展示信息输出")
	//其他设定
	//flag.BoolVar(&Info.IsInvestRd, "uncertain-invest", false, "包括未公示投资公司（无法确定占股比例）")
	//flag.BoolVar(&Info.IsGroup, "is-group", false, "查询关键词为集团")
	//flag.BoolVar(&Info.ISKeyPid, "is-pid", false, "批量查询文件是否为公司PID")
	flag.IntVar(&Info.DelayTime, "delay", 0, "填写最大延迟时间（秒）将会在1-n间随机延迟")
	flag.StringVar(&Info.Proxy, "proxy", "", "设置代理例如:-proxy=http://127.0.0.1:7897")
	flag.IntVar(&Info.TimeOut, "timeout", 1, "每个请求默认1（分钟）超时")
	//flag.BoolVar(&Info.IsMerge, "no-merge", false, "批量查询【取消】合并导出")
	//flag.BoolVar(&Info.Version, "v", false, "版本信息")
	//flag.BoolVar(&Info.IsEmailPro, "email", false, "获取email信息")
	flag.Parse()
}
