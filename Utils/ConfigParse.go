package Utils

import (
	"OneLong/Utils/gologger"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

func ConfigParse(options *ENOptions) {
	// 配置文件检查
	if ok, _ := PathExists(cfgYName); !ok {
		gologger.Infof("未发现配置文件，创建配置文件，请从新执行命令\n")
		f, errs := os.Create(cfgYName) //创建文件
		_, errs = io.WriteString(f, configYaml)
		if errs != nil {
			gologger.Fatalf("配置文件创建失败 %s\n", errs)
			os.Exit(0)
		}
		gologger.Infof("配置文件生成成功\n")
		os.Exit(0)
	}

	//加载配置信息~
	conf := new(ENConfig)
	yamlFile, err := ioutil.ReadFile(cfgYName)
	if err != nil {
		gologger.Fatalf("配置文件解析错误 #%v ", err)
	}
	if err := yaml.Unmarshal(yamlFile, conf); err != nil {
		gologger.Fatalf("【配置文件加载失败】: %v", err)
	}

	//初始化输出文件夹位置
	if options.Output == "" && conf.Utils.Output != "" {
		options.Output = conf.Utils.Output
	} else if options.Output == "" {
		options.Output = "outs"
	}

	if options.KeyWord == "" && options.Domain == "" {
		gologger.Errorf("参数输入错误")
		os.Exit(0)
	}
	options.IsShow = false
	options.IsMergeOut = true
	options.Deep = 5
	//if options.KeyWord == "" {
	//	options.ScanType = "aqc"
	//}
	////数据源判断 默认为爱企查
	//if options.ScanType == "" && len(options.GetType) == 0 {
	//	options.ScanType = "aqc"
	//}
	////如果是指定全部数据
	//if options.ScanType == "all" {
	//	options.GetType = []string{"aqc", "tyc"}
	//	options.IsMergeOut = true
	//} else if options.ScanType != "" {
	options.GetType = []string{"aqc", "tyc", "aldzs", "qimai"}
	//options.GetType = []string{"aqc"}
	//options.ScanType = "aqc,tyc,aldzs,qimai"
	//options.GetType = strings.Split(options.ScanType, ",")
	//}
	options.GetType = SetStr(options.GetType)
	var tmp []string
	for _, v := range options.GetType {
		if _, ok := ScanTypeKeys[v]; !ok {
			gologger.Errorf("没有这个%s查询方式\n支持列表\n%s", v, ScanTypeKeys)
		} else {
			tmp = append(tmp, v)
		}
	}
	options.GetType = tmp

	////判断是否添加墨子任务
	//if options.IsBiuCreate {
	//	if conf.Biu.Api == "" || conf.Biu.Key == "" {
	//		gologger.Fatalf("没有配置 墨子 API地址与Api key （请前往个人设置->安全设置中获取Api Key） \n")
	//	}
	//}

	//if len(conf.Biu.Tags) == 0 {
	//	conf.Biu.Tags = []string{"ENScan"}
	//}

	// 判断获取数据字段信息

	//options.GetField = SetStr(options.GetField)
	options.GetField = DefaultAllInfos
	//if options.GetFlags == "" && len(conf.Utils.Field) == 0 {
	//	if len(options.GetField) == 0 {
	//		options.GetField = DefaultInfos
	//	}
	//} else if options.GetFlags == "all" {
	//	options.GetField = DefaultAllInfos
	//} else {
	//	if len(conf.Utils.Field) > 0 {
	//		options.GetField = conf.Utils.Field
	//	}
	//	if options.GetFlags != "" {
	//		options.GetField = strings.Split(options.GetFlags, ",")
	//		if len(options.GetField) <= 0 {
	//			gologger.Fatalf("没有获取到字段信息 \n" + options.GetFlags)
	//		}
	//
	//	}
	//}
	//是否获取分支机构
	if options.IsGetBranch {
		options.GetField = append(options.GetField, "branch")
	}
	// 投资信息如果不等于0，那就收集股东信息和对外投资信息
	if options.InvestNum != 0 {
		options.GetField = append(options.GetField, "invest")
		options.GetField = append(options.GetField, "partner")
	}
	if options.IsHold {
		options.GetField = append(options.GetField, "holds")
	}
	if options.IsSupplier {
		options.GetField = append(options.GetField, "supplier")
	}
	options.GetField = SetStr(options.GetField)
	//// 判断是否在给定范围内，防止产生入库问题
	//if options.IsApiMode {
	//	//var tmps []string
	//	//for _, v := range options.GetField {
	//	//	if _, ok := outputfile.ENSMapLN[v]; ok {
	//	//		tmps = append(tmps, v)
	//	//	} else {
	//	//		gologger.Debugf("%s不在范围内\n", v)
	//	//	}
	//	//}
	//	//options.GetType = tmps
	//}

	//if options.IsMerge == true {
	//	gologger.Infof("====已强制取消合并导出！====\n")
	//	options.IsMergeOut = false
	//}

	options.GetField = SetStr(options.GetField)

	options.ENConfig = conf

}
