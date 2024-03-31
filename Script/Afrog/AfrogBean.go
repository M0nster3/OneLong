package Afrog

type EnsGo struct {
	Name    string
	Field   []string //获取的字段名称 看JSON
	KeyWord []string //关键词
}

func GetENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"Afrog": {
			Name:    "漏扫结果",
			Field:   []string{"Url", "info", "infoname"},
			KeyWord: []string{"Url", "危害程度", "漏洞名称"},
		},
	}
	for k, _ := range ensInfoMap {
		ensInfoMap[k].KeyWord = append(ensInfoMap[k].KeyWord, "数据关联")
		ensInfoMap[k].Field = append(ensInfoMap[k].Field, "inFrom")
	}
	return ensInfoMap

}
