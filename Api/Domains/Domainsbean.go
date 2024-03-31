package Domains

type EnsGo struct {
	Name    string
	Field   []string //获取的字段名称 看JSON
	KeyWord []string //关键词
}

func GetENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Urls": {
			Name:    "备案查询子域",
			Field:   []string{"address", "hostname"},
			KeyWord: []string{"IP", "域名"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].KeyWord = append(ensInfoMap[k].KeyWord, "数据关联")
		ensInfoMap[k].Field = append(ensInfoMap[k].Field, "inFrom")
	}
	return ensInfoMap

}
