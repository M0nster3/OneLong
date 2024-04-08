package Port

type EnsGo struct {
	Name    string
	Field   []string //获取的字段名称 看JSON
	KeyWord []string //关键词
}

func GetENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Port": {
			Name:    "Port端口",
			Field:   []string{"Address", "PORT", "SERVICE", "VERSION"},
			KeyWord: []string{"IP", "端口", "协议", "Banner"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].KeyWord = append(ensInfoMap[k].KeyWord, "数据关联")
		ensInfoMap[k].Field = append(ensInfoMap[k].Field, "inFrom")
	}
	return ensInfoMap

}
