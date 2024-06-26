package Login

type EnsGo struct {
	name    string
	field   []string //获取的字段名称 看JSON
	keyWord []string //关键词
}

func getENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Login": {
			name:    "后台地址",
			field:   []string{"hostname", "title"},
			keyWord: []string{"后台地址", "Title"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].keyWord = append(ensInfoMap[k].keyWord, "数据关联")
		ensInfoMap[k].field = append(ensInfoMap[k].field, "inFrom")
	}
	return ensInfoMap

}
