package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//Float64frombytes ...
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

//Float64bytes ...
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

//CreateMd5 创建md5
func CreateMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//RandomRemov 随机取后 删除被选中的元素
func RandomRemov(strings []string, length int) ([]string, []string, error) {
	var res []string //需要返回对结果
	if len(strings) <= 0 {
		return nil, nil, errors.New("the length of the parameter strings should not be less than 0")
	}

	if length <= 0 || len(strings) <= length {
		return nil, nil, errors.New("the size of the parameter length illegal")
	}
	var temp []string //保存副本
	for _, v := range strings {
		temp = append(temp, v)
	}
	//打乱顺序
	for i := len(strings) - 1; i > 0; i-- {
		rand.Seed(time.Now().Unix())
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	for i := 0; i < length; i++ {
		res = append(res, strings[i])
	}
	//删除选中的元素
	for _, n := range res {
		for k, m := range temp {
			if n == m {
				temp = append(temp[:k], temp[k+1:]...)
			}
		}
	}
	return res, temp, nil
}

//Random 对数组进行一个获取随机的组合
func Random(strings []string, length int) ([]string, error) {
	var res []string //需要返回对结果
	if len(strings) <= 0 {
		return nil, errors.New("the length of the parameter strings should not be less than 0")
	}

	if length <= 0 || len(strings) <= length {
		return nil, errors.New("the size of the parameter length illegal")
	}

	for i := len(strings) - 1; i > 0; i-- {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	for i := 0; i < length; i++ {
		res = append(res, strings[i])
	}
	return res, nil
}

//GetDbByHost 根据域名获取数据库名
func GetDbByHost(hostName string) string {
	dbname := strings.Replace(hostName, `.`, `_`, -1)
	return dbname
}

//ReadDirList 读目录下所有的文件名
func ReadDirList(dir string, fixName string) ([]string, error) {
	var relateKey []string
	var r *regexp.Regexp
	if fixName != "" {
		r, _ = regexp.Compile(fixName)
	}
	dirList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, v := range dirList {
		if fixName != "" {
			if !r.MatchString(v.Name()) {
				continue
			}
		}
		relateKey = append(relateKey, dir+"/"+v.Name())
	}
	return relateKey, nil
}

//ParseLongkey 解析长尾词描述行
func ParseLongkey(keyLine string) (string, int) {
	//看手相(21586|0)|长尾词数:21586,百度指数:1227,百度pc检索量:8,百度移动检索量:90,竞价激烈程度:2300
	temp := strings.Split(keyLine, `|`)
	var k1 string
	var k2 string
	if len(temp) > 3 {
		k1 = temp[1]
		k2 = temp[2]
	} else {
		k1 = temp[0]
		k2 = temp[1]
	}
	pos1 := strings.Index(k1, "(")
	key := k1[:pos1]
	pos2 := strings.Index(k2, ")")
	weight, err := strconv.Atoi(k2[:pos2])
	if err != nil {
		return "", 0
	}
	return key, weight
}

//ReadFileLines 读文件内容，一行为一个元素写入slice，并返回
func ReadFileLines(path string) ([]string, error) {
	var items []string
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("读文件出错")
	}
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, line)
	}
	return items, nil
}

//ReadFileLines 读文件内容，一行为一个元素写入slice，并返回
func ReadDataLines(data string) ([]string, error) {
	var items []string
	reader := bytes.NewReader([]byte(data))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, line)
	}
	return items, nil
}

//ReadFileList 获取目录下所有指定后缀的文件，并返回列表
func ReadFileList(dir string, fixName string) ([]string, error) {
	return nil, nil
}

//ExitFileDir 是否存在该目录或文件
func ExitFileDir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//ExeCmd ...
func ExeCmd(param string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", param)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return string(res), errors.Wrap(err, "execmd fail")
	}
	return string(res), nil
}
func ExeCommand(cmdString string) (string, error) {
	cmd := exec.Command(cmdString)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	res := fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	return res, err
}
func ExeCommandShell1(cmdString string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}
	outStr := string(stdout.Bytes())
	return outStr, nil
}

func ExeCommandShell(cmdString string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}
	outStr := string(stdout.Bytes())
	//res := fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	//log.Info("cmd res:", res)
	return outStr, nil
}

//SubString ...
func SubString(src string, sub string) string {
	index := strings.LastIndex(src, sub)
	if index == -1 {
		return src
	}
	return src[:index] + sub
}

//TimeEncode ...
func TimeEncode(second int) string {
	var byTime = []int{24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"天", "小时", "分钟", "秒"}
	res := ""
	for i := 0; i < len(byTime); i++ {
		if second < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(second / byTime[i]))
		second = second % byTime[i]
		if temp > 0 {
			tempStr := strconv.FormatFloat(temp, 'f', -1, 64)
			res = res + tempStr + unit[i]
		}
	}
	return res
}

//GbkToUtf8 转码
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//GetSegSentence ....
func GetSegSentence(data string) []string {
	stop := []string{",", "，", "。", ".", "!", "！", "?", "？", "；"}
	//items := []rune(data)
	items := []rune(data)
	var res []string
	//i := 0
	stopNew := make([]interface{}, len(stop))
	for i, j := range stop {
		stopNew[i] = j
	}
	m := 0
	for k, v := range items {
		if SliceContain(stopNew, string(v)) {
			res = append(res, string(items[m:k+1]))
			m = k + 1
		}
	}
	lable := make([]string, 0)
	for _, v := range res {
		temp := []rune(v)
		lenTemp := utf8.RuneCountInString(v)
		lable = append(lable, string(temp[lenTemp-1]))
	}
	return res
}
func EnsureDir(name string) error {
	if !ExitFileDir(name) {
		err := os.MkdirAll(name, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
func SliceContain(str []interface{}, dep interface{}) bool {
	if len(str) <= 0 {
		return false
	}

	for _, v := range str {
		switch v.(type) {
		case string:
			str1, ok := v.(string)
			if !ok {
				return false
			}
			str2, ok := dep.(string)
			if !ok {
				return false
			}
			if str1 == str2 {
				return true
			}
		case int:
			str1, ok := v.(int)
			if !ok {
				return false
			}
			str2, ok := dep.(int)
			if !ok {
				return false
			}
			if str1 == str2 {
				return true
			}
		}
	}
	return false
}

//StringContain ...
func StringContain(str []string, dep string) bool {
	if len(str) <= 0 {
		return false
	}
	for _, v := range str {
		if v == dep {
			return true
		}
	}
	return false
}

//HanAbstract 提取中文
func HanAbstract(str string) string {

	r := []rune(str)
	//strSlice := []string{}
	cnstr := ""
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnstr = cnstr + string(r[i])
			//strSlice = append(strSlice, cnstr)
		}
	}
	//if 0 == len(strSlice) {
	//无中文，需要跳过，后面再找规律

	//}
	//fmt.Println("原字符串:", str, "    提取出的中文字符串:", cnstr)
	return cnstr

}
func AbstractDig(s string) (int, error) {
	var valid = regexp.MustCompile("[0-9]")
	matchs := valid.FindAllStringSubmatch(s, -1)
	var res string
	for _, v := range matchs {
		res = res + v[0]
	}
	value, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return value, nil

}
func FormatXml(data string) string {
	//data := string(content)
	count := len(data)
	index := 0
	header := 0
	tail := 0
	temple := ""
	start := 1
	res := ""
	//log.Info("count:", count)
	for index < count-1 {
		if string(data[index]) == `<` && start == 1 {
			header = index
			start = 0
		} else {
			if string(data[index]) == `>` && string(data[index+1]) == `<` {
				tail = index
				temple = data[header : tail+1]
				start = 1
				res = res + temple + "\r\n"
				//fmt.Println("temp:", temple)
			}
		}
		index++
	}
	res = res + strings.TrimSpace(data[header:])
	return res
}
func MustLoadLocation(name string) (*time.Location, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func GetLocalTime(t string) (time.Time, error) {
	var tt time.Time
	loc, err := MustLoadLocation("Asia/Shanghai")
	if err != nil {
		return tt, err
	}
	tt, err = time.ParseInLocation("2006-01-02 15:04:05", t, loc)
	if err != nil {
		return tt, err
	}
	return tt, nil
}

///America/New_York
func GetZoneTime(t string, zone string) (time.Time, error) {
	var tt time.Time
	loc, err := MustLoadLocation(zone)
	if err != nil {
		return tt, err
	}
	tt, err = time.ParseInLocation("2006-01-02 15:04:05", t, loc)
	if err != nil {
		return tt, err
	}
	return tt, nil
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}
