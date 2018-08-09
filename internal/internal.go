package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// 获取程序跟目录
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}
	return file
}

// 时间转int64
func DateStrToint64(str, format string) int64 {
	theTime, err := time.ParseInLocation(format, str, time.Local)
	var unixTime int64
	if err == nil {
		unixTime = theTime.Unix()
	} else {
		unixTime = time.Now().Unix()
	}
	return unixTime
}

// 正则解析 `(?P<abc>Hello)(.*)(?P<cba>Go).`
func RegStrToMap(regStr string, str string) (map[string]string, error) {
	// 正则对象
	reg := regexp.MustCompile(regStr)
	// 所有定义的下标
	names := reg.SubexpNames()
	// 返回数据map
	strMap := make(map[string]string)
	// 取数据
	i := 0
	for _, v := range names {
		if v == "" {
			continue
		}
		i++
		// 获取单个值
		strOne := reg.ReplaceAllString(str, "$"+v)
		// 当有一个未匹配到会出现所有都是原字符串
		if strOne == str {
			continue
		}
		strMap[v] = strOne
	}
	if len(strMap) < i {
		return strMap, errors.New("未匹配全部属性")
	}
	return strMap, nil
}

// 格式化value，mongodb插入数据使用
func FormatMongoValue(instr interface{}, key, tableName string) (interface{}, int32) {
	// js, _ := json.Marshal(AppnameList)
	// log.Println(string(js))
	// 不好的方法转成字符串
	str := fmt.Sprint(instr)
	// log.Println(str)
	// log.Println(key)
	// log.Println(AppnameList[tableName].Fields[key].Type)
	// 保存返回值
	var myValue interface{}
	// 如果字段配置不存在则返回字符串，这样程序可以稳定运行，刚开始没有设置，后期设置需要删除数据表（或库）
	if AppnameList[tableName] == nil {
		return str, 0
	}
	if AppnameList[tableName].Fields[key] == nil {
		return str, 0
	}
	// 根据设置转换数据类型
	switch AppnameList[tableName].Fields[key].Type {
	case "string":
		myValue = str
		break
	case "int":
		if valueInt, err := strconv.Atoi(str); err == nil {
			myValue = valueInt
		} else {
			myValue = int(0)
		}
		break
	case "int64":
		if valueInt, err := strconv.ParseInt(str, 10, 64); err == nil {
			myValue = valueInt
		} else {
			myValue = int64(0)
		}
		break
	case "float64":
		if valueFloat, err := strconv.ParseFloat(str, 64); err == nil {
			myValue = valueFloat // int(valueFloat*100) / 100
			// log.Println(valueFloat)
		} else {
			myValue = float64(0.0)
		}
		break
	default:
		myValue = str
	}
	return myValue, AppnameList[tableName].Fields[key].Index
}

// 验证下标是否合法json
func VerifyMappedJsonKeys(data map[string]interface{}, keys map[string]string) bool {
	if len(data) != len(keys) {
		return false
	}
	// 判断所有下标值不是nil
	returnVal := true
	for _, v := range keys {
		if _, ok := data[v]; ok == false {
			returnVal = false
		}
	}
	return returnVal
}

// 验证下标是否合法regular
func VerifyMappedRegularKeys(data map[string]string, keys map[string]string) bool {
	if len(data) != len(keys) {
		return false
	}
	// 判断所有下标值不是nil
	returnVal := true
	for _, v := range keys {
		if _, ok := data[v]; ok == false {
			returnVal = false
		}
	}
	return returnVal
}

// try 防止崩溃
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

//将int32数组转换成string
func IntArrayToString(ids []int32) (str string) {
	for _, v := range ids {
		str += "," + strconv.Itoa(int(v))
	}
	return strings.TrimLeft(str, ",")
}

// 字符串模版替换--将一个模版字符串中的变量用map中的值替换
func TplAnalysisToString(tpl string, data *map[string]interface{}) (str string, err error) {
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}
	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)
	err = t.Execute(bw, data)
	bw.Flush()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// 获取当前可用数据源
func GetAlarmDataSource() (string, error) {
	// 读取是否配置了告警数据源
	dataSource, _ := CFG.String("statisd", "data_source")
	if dataSource != "" {
		return dataSource, nil
	}
	isinf, _ := CFG.Bool("influxdb", "enable")
	if isinf {
		dataSource = "influxdb"
	}
	ismgo, _ := CFG.Bool("mongodb", "enable")
	if ismgo {
		dataSource = "mongodb"
	}
	isela, _ := CFG.Bool("elasticsearch", "enable")
	if isela {
		dataSource = "elasticsearch"
	}
	if dataSource == "" {
		return "influxdb", errors.New("数据源配置错误")
	}
	return dataSource, nil
}
