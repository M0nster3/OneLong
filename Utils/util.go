package Utils

import (
	"OneLong/Utils/gologger"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RangeRand 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

// GetRandomString2 生成指定长度的随机字符串
func RangeString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// 获得配置文件的绝对路径
func GetPathDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// GetTempPathFileName 获取一个临时文件名
func GetTempPathFileName() (pathFileName string) {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.tmp", RangeString(16)))
}

// PathExists 判断文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// SetStr 数据去重
// target 输入数据
func SetStr(target []string) []string {
	setMap := make(map[string]int)
	var result []string
	for _, v := range target {
		if v != "" {
			if _, ok := setMap[v]; !ok {
				setMap[v] = 0
				result = append(result, v)
			}
		}
	}
	return result
}

func IsInList(target string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// CheckPid 检查pid是哪家单位
func CheckPid(pid string) (res string) {
	if len(pid) == 32 {
		res = "qcc"
	} else if len(pid) == 14 {
		res = "aqc"
	} else if len(pid) == 8 || len(pid) == 7 || len(pid) == 6 || len(pid) == 9 || len(pid) == 10 {
		res = "tyc"
	} else if len(pid) == 33 || len(pid) == 34 {
		if pid[0] == 'p' {
			gologger.Errorf("无法查询法人信息\n")
			res = ""
		} else {
			res = "xlb"
		}
	} else {
		gologger.Errorf("pid长度%d不正确，pid: %s", len(pid), pid)
		return ""
	}
	return res
}

func FormatInvest(scale string) float64 {
	if scale == "-" || scale == "" || scale == " " {
		return -1
	} else {
		scale = strings.ReplaceAll(scale, "%", "")
	}

	num, err := strconv.ParseFloat(scale, 64)
	if err != nil {
		gologger.Errorf("转换失败：%s\n", err)
		return -1
	}
	return num
}
func DName(str string) (srt string) { // 获得文件名
	str = strings.ReplaceAll(str, "(", "（")
	str = strings.ReplaceAll(str, ")", "）")
	str = strings.ReplaceAll(str, "<em>", "")
	str = strings.ReplaceAll(str, "</em>", "")
	return str
}

//func VerifyEmailFormat(email string) bool {
//	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
//	reg := regexp.MustCompile(pattern)
//	return reg.MatchString(email)
//}
