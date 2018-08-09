package mongo

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"53it.net/zql"
	"53it.net/zues/internal"
)

// 添加多条数据
func AddRawDataAll(rd []*map[string]interface{}, tableName string) error {
	var err error
	for _, v := range rd {
		AddRawDataOne(v, tableName)
	}
	return err
}

// 添加一条
func AddRawDataOne(rdFinger *map[string]interface{}, tableName string) (err error) {
	// mutex.Lock()
	// 数据连接和选库
	sessionOne, selectDb := dbmgo()
	defer func() {
		sessionOne.Close()
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		// mutex.Unlock()
		return
	}()
	rd := *rdFinger
	collection := selectDb.C(getMongodbTableName(tableName)) // 使用配置table表名也就是appname
	// 格式化字段数据类型－排除datetime，kpiid
	for kk, vv := range rd {
		if kk == "kpiid" || kk == "date" {
			continue
		}
		if kk == "datetime" {
			// 时间戳转时间对象
			if dateInt64, ok := vv.(int64); ok {
				// log.Println("ok", time.Unix(dateInt64, 0).Local())
				rd["date"] = time.Unix(dateInt64, 0).Local()
				rd["datetime"] = dateInt64
			} else {
				rd["date"] = time.Now().Local()
				rd["datetime"] = time.Now().Unix()
				// log.Println("no", time.Now().Local())
			}
			continue
		}
		if vv != nil {
			rd[kk], _ = internal.FormatMongoValue(vv, kk, tableName)
		} else {
			delete(rd, kk)
		}
	}
	err = collection.Insert(rd)
	if err != nil {
		internal.LogFile.W("mongodb数据插入失败:" + err.Error())
	}
	return err
}

// 使用zql插入
func InsertIntoZQL(zqlStr string) error {
	objZQL, err := zql.New("", zqlStr)
	if err != nil {
		return err
	}
	vals, tname := objZQL.GetInsertIntoData()
	return InsertInto(vals, tname)
}

// 插入数据
func InsertInto(dataFinger *map[string]interface{}, tableName string) (err error) {
	// mutex.Lock()
	// 数据连接和选库
	sessionOne, selectDb := dbmgo()
	defer func() {
		sessionOne.Close()
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		// mutex.Unlock()
		return
	}()
	data := *dataFinger
	collection := selectDb.C(getMongodbTableName(tableName))
	dataInsert := make(map[string]interface{}, 0)
	// 格式化字段数据类型－排除datetime，kpiid
	for kk, vv := range data {
		if kk == "datetime" {
			vvvv, err := strconv.ParseInt(fmt.Sprint(vv), 10, 64)
			if err == nil {
				dataInsert["date"] = time.Unix(vvvv, 0).Local() //.Format("2006-01-02 15:04:05")
			} else {
				dataInsert["date"] = time.Now()
				internal.LogFile.E("mongodb数据插入失败,时间转换错误:" + err.Error())
			}
		}
		dataInsert[kk], _ = internal.FormatMongoValue(vv, kk, tableName)
	}
	err = collection.Insert(dataInsert)
	if err != nil {
		internal.LogFile.W("mongodb数据插入失败:" + err.Error())
	}
	return err
}

// 查询测试 zql
func GetZqlList(zqlStr string, args ...time.Time) (list []map[string]interface{}, err error) {
	// 数据连接和选库
	sessionOne, selectDb := dbmgo()
	defer func() {
		sessionOne.Close()
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	zqlObj, err := zql.New("", zqlStr)
	if err != nil {
		return nil, err
	}
	// list := make([]map[string]interface{}, 0)
	if len(args) > 0 {
		err = zqlObj.GetMongoQuery(selectDb, getMongodbSubTableName(args[0]), &list)
	} else {
		err = zqlObj.GetMongoQuery(selectDb, getMongodbSubTableName(), &list)
	}
	return list, err
}
