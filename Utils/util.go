package Utils

import (
	"OneLong/Utils/gologger"
	"bufio"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 ",
	"(KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 ",
	"(KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 ",
	"(KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:61.0) ",
	"Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (X11; Linux i586; rv:31.0) Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (Macintosh) AppleWebKit/Linux x86_64 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	"Mozilla/5.0 (Android) AppleWebKit/Android (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Linux) AppleWebKit/Win64; x64 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	"Mozilla/5.0 (Macintosh) AppleWebKit/Macintosh Intel Mac OS X 10_15_7 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Android) AppleWebKit/iPad (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (iPhone) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Mozilla/5.0 (iPad) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Mozilla/5.0 (iPhone) AppleWebKit/Android (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (iPhone) AppleWebKit/Win64; x64 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/Android (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/Linux x86_64 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Mozilla/5.0 (iPhone) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	"Mozilla/5.0 (iPad) AppleWebKit/Macintosh Intel Mac OS X 10_15_7 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	"Mozilla/5.0 (iPad) AppleWebKit/Android (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	"Mozilla/5.0 (Linux) AppleWebKit/Win64; x64 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Linux) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/Linux x86_64 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	"Mozilla/5.0 (Android) AppleWebKit/Macintosh Intel Mac OS X 10_15_7 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	"Mozilla/5.0 (Macintosh) AppleWebKit/iPhone (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
}

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

func RandUA() string {
	// 生成一个随机数作为索引
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(UserAgents))))
	if err != nil {
		panic(err)
	}
	randomUserAgent := UserAgents[randomIndex.Int64()]
	return randomUserAgent
}

// GetRandomString2 生成指定长度的随机字符串
func RangeString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func ReadFile(filename string) []string {
	var result []string
	if FileExists(filename) {
		f, err := os.Open(filename)
		if err != nil {
			gologger.Fatalf("Read fail: %v\n", err)
		}
		fileScanner := bufio.NewScanner(f)
		// read line by line
		for fileScanner.Scan() {
			result = append(result, fileScanner.Text())
		}
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			gologger.Fatalf("Error while reading file: %s\n", err)
		}
		_ = f.Close()
	}
	result = SetStr(result)
	return result
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

// GetTempPathFileName 获取一个临时Json文件名
func GetTempJsonPathFileName() (pathFileName string) {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.json", RangeString(16)))
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
		v = strings.Trim(v, ".")
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
