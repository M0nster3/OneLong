package Utils

import (
	"OneLong/Utils/gologger"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

func ConfigParse(options *LongOptions) {
	// 配置文件检查
	if options.Low {
		options.NoPoc = true
		options.NoBao = true
		options.NoPort = true
	}
	options.IsHold = true
	options.IsSupplier = true
	options.IsGetBranch = true
	options.GetField = DefaultAllInfos
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
	conf := new(LongConfig)
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

	if options.KeyWord == "" && options.Domain == "" && options.InputFile == "" {
		gologger.Errorf("参数输入错误")
		os.Exit(0)
	}
	options.LongConfig = conf
	options.IsShow = false
	options.IsMergeOut = true
	options.Deep = 5
	options.GetType = []string{"aqc", "tyc", "aldzs", "qimai"}
	if options.LongConfig.Cookies.Aiqicha == "" {
		options.GetType = []string{"tyc", "aldzs", "qimai"}
	} else if options.LongConfig.Cookies.Tycid == "" || options.LongConfig.Cookies.Securitytrails == "" {
		options.GetType = []string{"aqc", "aldzs", "qimai"}
	}
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

	options.GetField = SetStr(options.GetField)

}
