package Utils

import (
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ENOptions struct {
	NoBao     bool
	NoPoc     bool
	Domain    string
	KeyWord   string // Keyword of Search
	CompanyID string // Company ID
	GroupID   string // Company ID
	InputFile string // Scan Input File
	Output    string
	ScanType  string
	Proxy     string

	IsGetBranch    bool //分支机构（分公司）信息
	IsSearchBranch bool

	InvestNum    float64
	DelayTime    int
	DelayMaxTime int64
	TimeOut      int

	IsHold      bool //控股公司
	IsSupplier  bool //供应商信息
	IsShow      bool
	CompanyName string
	GetField    []string
	GetType     []string
	Deep        int

	IsMergeOut bool //合并导出
	IsMerge    bool //聚合
	ICP        []string
	ENConfig   *ENConfig
}

func (h *ENOptions) GetDelayRTime() int64 {
	if h.DelayTime != 0 {
		h.DelayMaxTime = int64(h.DelayTime)
	}
	if h.DelayMaxTime == 0 {
		return 0
	}
	return RangeRand(1, h.DelayMaxTime)
}

// ENConfig YML配置文件，更改时注意变更 cfgYV 版本
type ENConfig struct {
	//Version float64 `yaml:"version"`
	Utils struct {
		Output string `yaml:"output"`
	}
	Biu struct {
		Api      string   `yaml:"api"`
		Key      string   `yaml:"key"`
		Port     string   `yaml:"port"`
		IsPublic bool     `yaml:"is-public"`
		Tags     []string `yaml:"tags"`
	}
	//Api struct {
	//	Server  string `yaml:"server"`
	//	Mongodb string `yaml:"mongodb"`
	//	Redis   string `yaml:"redis"`
	//}
	Cookies struct {
		Aldzs          string `yaml:"aldzs"`
		Xlb            string `yaml:"xlb"`
		Aiqicha        string `yaml:"aiqicha"`
		Binaryedge     string `yaml:"binaryedge"`
		Tianyancha     string `yaml:"tianyancha"`
		Tycid          string `yaml:"tycid"`
		Qcc            string `yaml:"qcc"`
		QiMai          string `yaml:"qimai"`
		Fullhunt       string `yaml:"fullhunt"`
		Hunter         string `yaml:"hunter"`
		Bevigil        string `yaml:"bevigil"`
		CensysToken    string `yaml:"CensysToken"`
		CensysSecret   string `yaml:"CensysSecret"`
		Zoomeye        string `yaml:"zoomeye"`
		Whoisxmlapi    string `yaml:"whoisxmlapi"`
		Virustotal     string `yaml:"virustotal"`
		Shodan         string `yaml:"shodan"`
		Chaos          string `yaml:"chaos"`
		Leakix         string `yaml:"leakix"`
		Netlas         string `yaml:"netlas"`
		Quake          string `yaml:"quake"`
		Securitytrails string `yaml:"securitytrails"`
		GoogleID       string `yaml:"googleid"`
		GoogleApi      string `yaml:"googleapi"`
		FofaKey        string `yaml:"fofaKey"`
		FofaEmail      string `yaml:"fofaEmail"`
		Github         string `yaml:"githubtoken"`
		Racent         string `yaml:"racent"`
	}
	Massdns struct {
		Resolvers   string `yaml:"resolvers"`
		Wordlist    string `yaml:"wordlist"`
		MassdnsPath string `yaml:"massdnsPath"`
	}
	Email struct {
		Emailhunter string `yaml:"emailhunter"`
		EmailIntelx string `yaml:"intelxEmail"`
		TombaKey    string `yaml:"tombaKey"`
		TombaSecret string `yaml:"tombaSecret"`
	}
}

type EnInfos struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string
	Pid         string
	LegalPerson string
	OpenStatus  string
	Email       string
	Telephone   string
	SType       string
	RegCode     string
	BranchNum   int64
	InvestNum   int64
	InTime      time.Time
	PidS        map[string]string
	Infos       map[string][]gjson.Result
	EnInfos     map[string][]map[string]interface{}
	EnInfo      []map[string]interface{}
}

type DBEnInfos struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string
	RegCode     string
	InTime      time.Time
	InvestCount int
	InfoCount   map[string][]string
	Info        []map[string]interface{}
}

var ENSMapAQC = map[string]string{
	"webRecord":     "icp",
	"appinfo":       "app",
	"wechatoa":      "wechat",
	"enterprisejob": "job",
	"microblog":     "weibo",
	"hold":          "holds",
	"shareholders":  "partner",
}

// // DefaultAllInfos 默认收集信息列表
var DefaultAllInfos = []string{"icp", "weibo", "wechat", "app", "weibo", "job", "wx_app", "copyright"}

// var CanSearchAllInfos = []string{"enterprise_info", "icp", "weibo", "wechat", "app", "weibo", "job", "wx_app", "copyright", "supplier", "invest", "branch", "holds", "partner"}
var ScanTypeKeys = map[string]string{
	"aqc":     "爱企查",
	"qcc":     "企查查",
	"tyc":     "天眼查",
	"xlb":     "小蓝本",
	"all":     "全部查询",
	"aldzs":   "阿拉丁",
	"coolapk": "酷安市场",
	"qimai":   "七麦数据",
	"chinaz":  "站长之家",
}

//var ScanTypeKeyV = map[string]string{
//	"爱企查":  "aqc",
//	"企查查":  "qcc",
//	"天眼查":  "tyc",
//	"小蓝本":  "xlb",
//	"阿拉丁":  "aldzs",
//	"酷安市场": "coolapk",
//	"七麦数据": "qimai",
//	"站长之家": "chinaz",
//}

// // RequestTimeOut 请求超时设置
// var RequestTimeOut = 30 * time.Second
// var (
//
//	BuiltAt   string
//	GoVersion string
//	GitAuthor string
//	BuildSha  string
//	GitTag    string
//
// )
var cfgYName = GetPathDir() + "/config.yaml"
var configYaml = `
Utils:
  output: ""            # 导出文件位置
cookies:
  aiqicha: ''           # 爱企查   Cookie
  tianyancha: ''        # 天眼查   Cookie
  tycid: ''        		# 天眼查   CApi ID(capi.tianyancha.com)
  aldzs: ''             # 阿拉丁   TOKEN(see README)
  qimai: ''             # 七麦数据  Cookie
  binaryedge: ''		# binaryedge  Cookie 免费查询250次
  fullhunt: ''			# Fullhunt Cookie 威胁平台 每月免费100次
  hunter: ''			# Hunter Cookie 威胁平台 每日免费500个数据
  bevigil: ''           # Bevigil Cookie 威胁平台 每月免费50次
  CensysToken: ''       # Censys Token 威胁平台 每月免费250次
  CensysSecret: ''      # Censys Secret 威胁平台 每月免费250次
  zoomeye: ''			# ZooEye Cookie 每月1000条
  whoisxmlapi: ''		# whoisxmlapi Cookie 免费500次
  virustotal: ''        # virustotal Cookie 每分钟4次 每天500次
  shodan: ''			# shodan Cookie 
  chaos: ''				# chaos Key
  leakix: ''			# leakix Key
  netlas: ''			# Netlas key
  quake: ''				# Quake key
  securitytrails: ''	# securitytrails 需要企业邮箱，每个月50次
  googleid: ''			# google id 免费的API只能查询前100条结果,每天免费提供 100 次搜索查询
  googleapi: ''			# google Api 免费的API只能查询前100条结果,每天免费提供 100 次搜索查询
  fofaKey: ''			# Fofa key
  fofaEmail: ''			# Fofa Email
  githubtoken: ''		# Github Token 
  racent: ''			# racent Token
massdns:
  resolvers: 'resolvers.txt'			# resolvers 文件名称
  wordlist: 'names.txt'			# 子域名爆破文件名称
  massdnsPath: 'massdns.exe'		# Massdns工具名称
email:
  emailhunter: ''		# Email hunter Token
  intelxEmail: ''		# Email Intelx Token
  tombaKey: ''			# Email tombaKey
  tombaSecret: ''		# Email tombaSecret
#Afrog配置
reverse:
  alphalog:
    domain: ""
    api_url: ""
  ceye:
    api-key: ""
    domain: ""
  dnslogcn:
    domain: dnslog.cn
  eye:
    host: ""
    token: ""
    domain: ""
  jndi:
    jndi_address: ""
    ldap_port: ""
    api_port: ""
  xray:
    x_token: ""
    domain: ""
    api_url: http://x.x.x.x:8777

`
