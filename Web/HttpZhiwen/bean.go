package HttpZhiwen

type EnsGo struct {
	name    string
	field   []string //获取的字段名称 看JSON
	keyWord []string //关键词
}

func getENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Zhiwen": {
			name:    "整合域名IP",
			field:   []string{"address", "hostname", "Size", "status_code", "title", "Zhiwen"},
			keyWord: []string{"整合收集IP", "整合收集域名", "Size", "响应码", "标题", "指纹"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].keyWord = append(ensInfoMap[k].keyWord, "数据关联")
		ensInfoMap[k].field = append(ensInfoMap[k].field, "inFrom")
	}
	return ensInfoMap

}
