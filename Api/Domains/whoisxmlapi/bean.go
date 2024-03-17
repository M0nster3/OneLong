package whoisxmlapi

type EnsGo struct {
	name    string
	field   []string //获取的字段名称 看JSON
	keyWord []string //关键词
}

func getENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Urls": {
			name:    "备案查询子域",
			field:   []string{"address", "hostname"},
			keyWord: []string{"IP", "域名"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].keyWord = append(ensInfoMap[k].keyWord, "数据关联  ")
		ensInfoMap[k].field = append(ensInfoMap[k].field, "inFrom")
	}
	return ensInfoMap

}
