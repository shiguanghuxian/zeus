package redis

// 最新的指标
import (
	"errors"
	"fmt"
	"log"
	"time"

	"53it.net/zues/internal"
)

// redis存放一个指标数据
type NewestData struct {
	Kpiid    string      `json:"kpiid"`
	Value    interface{} `json:"value"`
	Datetime int64       `json:"datetime"`
	Instance string      `json:"instance"`
}

// 把最新的数据放到redis
func SaveNewestData(rawDataAll []*map[string]interface{}, tableName string) error {
	// 循环插入数据
	for _, v := range rawDataAll {
		return SaveRawDataOne(v, tableName)
	}
	return nil
}

// 保存单条数据
func SaveRawDataOne(rawDataFinger *map[string]interface{}, tableName string) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	rawData := *rawDataFinger
	// 主下标
	key := tableName + ":" + fmt.Sprint(rawData["group"]) + ":" + fmt.Sprint(rawData["hostname"]) + ":" + fmt.Sprint(rawData["ip"])
	// 格式化字段数据类型－排除datetime，kpiid，instance
	for kk, vv := range rawData {
		if kk == "datetime" {
			continue
		}
		if vv != nil {
			rawData[kk], _ = internal.FormatMongoValue(vv, kk, tableName)
		} else {
			delete(rawData, kk)
		}
	}
	// 插入数据
	resp := dbRedis().Cmd("HMSET", key, rawData)
	// 添加情况
	if resp.Err != nil {
		internal.LogFile.E("redis添加最新数据失败:" + resp.Err.Error())
		return resp.Err
	}
	return nil
	// 保存设备信息
	// return SaveDeviceInfo(rawData)
}

// 读取redis最新数据
func GetOneNewestData(rawData map[string]string, tableName, field string) (string, error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	if field == "0" || tableName == "" {
		return "", errors.New("tableName和field不能为空")
	}
	// key
	key := tableName + ":" + rawData["group"] + ":" + rawData["hostname"] + ":" + rawData["ip"]
	// log.Println(key)
	// log.Println(field)
	// 取数据
	valBytes, err := dbRedis().Cmd("HGET", key, field).Bytes()
	if err != nil {
		internal.LogFile.E("redis读取最新数据错误:" + err.Error())
		return "", err
	}
	return string(valBytes), nil
}

// 获取一个设备所有指标
func GetAllNewestData(rawData map[string]string, tableName string) (valmap map[string]string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	// key
	key := tableName + ":" + rawData["group"] + ":" + rawData["hostname"] + ":" + rawData["ip"]
	// log.Println(key)
	valmap, err = dbRedis().Cmd("HGETALL", key).Map()
	if err != nil {
		return nil, err
	}
	return valmap, err
}

// SaveDeviceInfo 存储设备列表
func SaveDeviceInfo(rawData map[string]interface{}) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	// 记录设备id的key
	rkey := fmt.Sprintf("deviceids:%s:%s:%s", rawData["group"], rawData["hostname"], rawData["ip"])
	// 记录设备信息
	key := "devicelist:"
	if rawData["id"] != nil {
		key += fmt.Sprint(rawData["id"]) + ":"
		log.Println(rkey, "----", fmt.Sprint(rawData["id"]))
		SetKeyVal(rkey, fmt.Sprint(rawData["id"]))
	} else {
		key += "0:"
		SetKeyVal(rkey, "0")
	}
	key += fmt.Sprint(rawData["group"]) + ":" + fmt.Sprint(rawData["device_type"]) + ":" + fmt.Sprint(rawData["hostname"]) + ":" + fmt.Sprint(rawData["ip"])
	resp := dbRedis().Cmd("HMSET", key, rawData)
	// 添加情况
	if resp.Err != nil {
		internal.LogFile.E("redis保存设备到列表失败:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}

// 根据key删除数据
func DelDeviceInfo(rawData map[string]interface{}) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	key := "devicelist:"
	if rawData["id"] != nil {
		key += fmt.Sprint(rawData["id"]) + ":"
	} else {
		key += "0:"
	}
	key += fmt.Sprint(rawData["group"]) + ":" + fmt.Sprint(rawData["device_type"]) + ":" + fmt.Sprint(rawData["hostname"]) + ":" + fmt.Sprint(rawData["ip"])

	resp := dbRedis().Cmd("DEl", key)
	// 添加情况
	if resp.Err != nil {
		internal.LogFile.E("redis删除设备失败:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}

// 模糊查询keys列表
func GetKeysList(prefix string) (keyList []string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	keyList, err = dbRedis().Cmd("KEYS", prefix+"*").List()
	if err != nil {
		internal.LogFile.E("redis keys列表:" + err.Error())
		return keyList, err
	}
	return keyList, nil
}
