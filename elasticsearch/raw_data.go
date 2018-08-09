package elasticsearch

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"53it.net/zql"
	"53it.net/zues/internal"
)

// var mutex1 sync.Mutex

// 添加一条
func AddRawDataOne(rdFinger *map[string]interface{}, tableName string) (err error) {
	idleConn, err := getClient()
	defer func() {
		releaseClient(idleConn)
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("elastic致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	rd := *rdFinger
	// 处理字段
	for kk, vv := range rd {
		if kk == "date" || kk == "kpiid" {
			continue
		}
		// 保留一个时间字段
		if kk == "datetime" {
			// 时间戳转时间对象
			if dateInt64, ok := vv.(int64); ok {
				//				log.Println("ok", time.Unix(dateInt64, 0).Local())
				rd["date"] = time.Unix(dateInt64, 0).Local()
				rd["datetime"] = dateInt64
			} else {
				rd["date"] = time.Now().Local()
				rd["datetime"] = time.Now().Unix()
				//				log.Println("no", time.Now().Local())
			}
			continue
		}

		if vv != nil {
			rd[kk], _ = internal.FormatMongoValue(vv, kk, tableName)
		} else {
			delete(rd, kk)
		}
	}
	// 执行插入
	_, err = idleConn.C.Index().
		Index(getIndexName()).
		Type(tableName).
		BodyJson(rd).
		Refresh(true).
		Do()
	if err != nil {
		internal.LogFile.W("elasticsearch数据插入失败:" + err.Error())
	}
	// log.Println(time.Now().Format("2006-01-02 15:04:05"))
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
	data := *dataFinger
	idleConn, err := getClient()
	defer func() {
		releaseClient(idleConn)
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	// 处理字段
	dataInsert := make(map[string]interface{}, 0)
	for kk, vv := range data {
		// 保留一个时间字段
		if kk == "datetime" {
			vvvv, err := strconv.ParseInt(fmt.Sprint(vv), 10, 64)
			if err == nil {
				dataInsert["date"] = time.Unix(vvvv, 0).Local() //.Format("2006-01-02T15:04:05")
			} else {
				dataInsert["date"] = time.Now()
				internal.LogFile.E("elasticsearch数据插入失败,时间转换错误:" + err.Error())
			}
		}
		dataInsert[kk], _ = internal.FormatMongoValue(vv, kk, tableName)
	}
	// 执行插入
	_, err = idleConn.C.Index().
		Index(getIndexName()).
		Type(tableName).
		BodyJson(dataInsert).
		Refresh(true).
		Do()
	if err != nil {
		internal.LogFile.W("elasticsearch数据插入失败:" + err.Error())
	}
	return err
}

// zql查询数据
func GetZqlList(zqlStr string) ([]map[string]interface{}, error) {
	idleConn, err := getClient()
	releaseClient(idleConn)
	zqlObj, err := zql.New("", zqlStr)
	if err != nil {
		return nil, err
	}
	objList, err := zqlObj.GetElasticQuery(idleConn.C, getIndexName(), true)
	// log.Println(zqlObj.GetElasticQueryStr()) // 输出查询字符串
	if err != nil {
		return nil, err
	}
	return objList, nil
}
