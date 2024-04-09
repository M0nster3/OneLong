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
	color.RGBStyleFromString("237,64,35").Println("工具仅用于信息收集，请勿用于非法用途\n开发人员不承担任何责任，也不对任何滥用或损坏负责.\n")
	color.RGBStyleFromString("244,211,49").Println("使用方式: \n\tOneLong -n 企业名称\n\tOneLong -d target.com\n\tOneLong -f file.txt\n")
}

func Flag(Info *LongOptions) {
	Banner()
	flag.BoolVar(&Info.Low, "low", false, "只进行被动扫描，不进行子域名以及端口爆破，不进行漏洞扫描")
	flag.BoolVar(&Info.NoBao, "nbao", false, "不进行爆破子域名")
	flag.BoolVar(&Info.NoPoc, "npoc", false, "不进行漏洞扫描")
	flag.BoolVar(&Info.NoPort, "nport", false, "不进行爆破端口")
	flag.StringVar(&Info.Port, "p", "--top-ports 1000", "端口爆破默认为Top1000，还可修改为100或者10")
	flag.StringVar(&Info.InputFile, "f", "", "批量扫描")
	flag.StringVar(&Info.KeyWord, "n", "", "企业关键词 eg 百度")
	flag.StringVar(&Info.Domain, "d", "", "域名")
	flag.StringVar(&Info.Output, "o", "", "结果输出的文件夹位置(可选)")
	flag.Float64Var(&Info.InvestNum, "invest", 70, "投资比例 ")
	flag.IntVar(&Info.Deep, "deep", 5, "递归搜索n层公司")
	flag.BoolVar(&Info.IsSearchBranch, "is-branch", false, "深度查询分支机构信息（数量巨大），默认不查询")
	flag.IntVar(&Info.DelayTime, "delay", 0, "填写最大延迟时间（秒）将会在1-n间随机延迟")
	flag.StringVar(&Info.Proxy, "proxy", "", "设置代理例如:-proxy=http://127.0.0.1:7897")
	flag.IntVar(&Info.TimeOut, "timeout", 1, "每个请求默认1（分钟）超时")
	flag.Parse()
}
