package influxdb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"

	"strconv"

	"53it.net/zql"
	"53it.net/zues/internal"
)

// 添加一条
func AddRawDataOne(rdFinger *map[string]interface{}, tableName string) (err error) {
	idleConn, err := dbInflux()
	defer func() {
		releaseClient(idleConn)
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("influxdb致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	rd := *rdFinger
	// 获取表名--需要拼接上设备信息，否则出现数据覆盖
	tName := getInfluxdbTableName(tableName, fmt.Sprint(rd["group"]), fmt.Sprint(rd["hostname"]), strings.Replace(fmt.Sprint(rd["ip"]), ".", "_", -1))
	// 格式化字段数据类型
	tags := make(map[string]string)        // 带索引
	fields := make(map[string]interface{}) // 不带索引
	pointTime := time.Now()
	for kk, vv := range rd {
		if vv == nil {
			continue
		}
		// 上传时间用上传的，否则使用当前时间
		if kk == "datetime" {
			if vvvv, ok := vv.(int64); ok {
				pointTime = time.Unix(vvvv, 0)
			}
		}
		if kk == "ip" || kk == "group" || kk == "hostname" {
			tags[kk] = fmt.Sprint(vv)
		} else {
			fmtVal, index := internal.FormatMongoValue(vv, kk, tableName)
			if index == 1 {
				tags[kk] = fmt.Sprint(fmtVal)
			} else {
				fields[kk] = fmtVal
			}
		}
		rd["datetime"] = pointTime.Unix() // 时间戳用数值格式
	}
	// 创建一个点
	influxBatchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  getInfluxdbName(),
		Precision: "ms",
	})
	if err != nil {
		internal.LogFile.E("influxdb选库失败:" + err.Error())
		log.Println("influxdb选库失败:" + err.Error())
	}
	//	log.Println(pointTime)
	//	log.Println(time.Now())
	pt, err := client.NewPoint(tName, tags, fields, pointTime)
	if err != nil {
		internal.LogFile.W("influxdb创建点错误:" + err.Error())
	}
	influxBatchPoints.AddPoint(pt)
	// 写点
	err = idleConn.C.Write(influxBatchPoints)
	if err != nil {
		if strings.Index(err.Error(), "database not found:") > 0 {
			mutex.Lock()
			q := client.Query{
				Command: fmt.Sprintf("CREATE DATABASE \"%s\"", influxDbName),
			}
			_, err = idleConn.C.Query(q)
			if err != nil {
				internal.LogFile.E("influxdb 创建数据库:" + err.Error())
				log.Println("influxdb 创建数据库:" + err.Error())
				os.Exit(2)
			}
			mutex.Unlock()
			err = idleConn.C.Write(influxBatchPoints) // 再次插入
		}
	}
	if err != nil {
		internal.LogFile.W("influxdb数据插入失败:" + err.Error())
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
	idleConn, err := dbInflux()
	defer func() {
		releaseClient(idleConn)
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New(fmt.Sprint(errInfo))
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
		}
		return
	}()
	data := *dataFinger
	// 如果设备主机名、ip、原始分组信息都有则拼接表名
	if data["group"] != nil && data["hostname"] != nil && data["ip"] != nil {
		tableName = getInfluxdbTableName(tableName, fmt.Sprint(data["group"]), fmt.Sprint(data["hostname"]), strings.Replace(fmt.Sprint(data["ip"]), ".", "_", -1))
	}
	// 格式化字段数据类型
	tags := make(map[string]string, 0)        // 带索引
	fields := make(map[string]interface{}, 0) // 不带索引
	pointTime := time.Now()
	for kk, vv := range data {
		if vv == nil {
			continue
		}
		// 上传时间用上传的，否则使用当前时间
		if kk == "datetime" {
			vvvv, err := strconv.ParseInt(fmt.Sprint(vv), 10, 64)
			if err == nil {
				pointTime = time.Unix(vvvv, 0).Local()
			}
		}
		// 其它字段处理，识别出是否索引
		fmtVal, index := internal.FormatMongoValue(vv, kk, tableName)
		if index == 1 {
			tags[kk] = fmt.Sprint(fmtVal)
		} else {
			fields[kk] = fmtVal
		}
	}
	// 创建一个点
	influxBatchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  getInfluxdbName(),
		Precision: "ms",
	})
	if err != nil {
		internal.LogFile.E("influxdb选库失败:" + err.Error())
		log.Println("influxdb选库失败:" + err.Error())
	}
	pt, err := client.NewPoint(tableName, tags, fields, pointTime)
	if err != nil {
		internal.LogFile.W("influxdb创建点错误:" + err.Error())
	}
	influxBatchPoints.AddPoint(pt)
	// 写点
	err = idleConn.C.Write(influxBatchPoints)
	if err != nil {
		internal.LogFile.W("influxdb数据插入失败:" + err.Error())
	}
	return err
}

// 根据查询字符串查询
func QueryCmd(cmd string) (res []client.Result, err error) {
	idleConn, err := dbInflux()
	defer releaseClient(idleConn)
	q := client.Query{
		Command:  cmd,
		Database: getInfluxdbName(),
	}
	if response, err := idleConn.C.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		internal.LogFile.W("influxdb数据数据查询错误:" + err.Error())
		return res, err
	}
	return res, nil
}

// 返回类似数据库结果
func QueryResponse(cmd string) ([]map[string]interface{}, error) {
	// 查询数据
	response, err := QueryCmd(cmd)
	if err != nil {
		return nil, err
	}
	if len(response) < 1 {
		return make([]map[string]interface{}, 0), nil
	}
	if len(response[0].Series) < 1 {
		return make([]map[string]interface{}, 0), nil
	}
	// 字段列表和数据列表
	columns := response[0].Series[0].Columns
	values := response[0].Series[0].Values
	if len(values) == 0 {
		return nil, errors.New("No query to data")
	}
	// 最终返回数据
	var datas []map[string]interface{}
	for _, v := range values {
		// 每行数据
		item := make(map[string]interface{})
		for kk, vv := range columns {
			item[vv] = v[kk]
		}
		datas = append(datas, item)
	}
	return datas, nil
}

// zql查询, group, hostname, ip
func ZqlQueryCmd(zqlStr string) ([]map[string]interface{}, error) {
	i, group := analyseFieldValue(zqlStr, "group")
	_, hostname := analyseFieldValue(zqlStr, "hostname")
	_, ip := analyseFieldValue(zqlStr, "ip")
	if group == "" || hostname == "" || ip == "" {
		return nil, errors.New("条件中group、hostname、ip都不能为空")
	}
	// 处理group是关键词，导致的查询错误
	// log.Println(i)
	if i > 0 {
		zqlStr = zqlStr[:i] + "\"group\"" + zqlStr[(i+5):]
	}
	suffix := "_" + group + "_" + hostname + "_" + strings.Replace(ip, ".", "_", -1)
	// zql查询对象
	zqlObj, err := zql.New("", zqlStr)
	if err != nil {
		return nil, errors.New("创建查询对象错误:" + err.Error())
	}
	zqlQuery, err := zqlObj.GetInfluxdbQuery(suffix)
	if err != nil {
		return nil, errors.New("创建查询语句错误:" + err.Error())
	}
	// log.Println(zqlQuery)
	return QueryResponse(zqlQuery)
}

// 根据字段名分析出值
func analyseFieldValue(zql, field string) (int, string) {
	if field == "" {
		return 0, ""
	}
	i := strings.Index(zql, field)
	if i < 1 {
		return 0, ""
	}
	subZql := strings.TrimSpace(zql[(i + len(field)):])
	if subZql[:1] == "\"" {
		subZql = strings.TrimSpace(subZql[1:])
		i = 0
	}
	if subZql[:1] == "=" {
		subZql = strings.TrimSpace(subZql[1:])
		// 查询小括号位置
		j := strings.Index(subZql, ")")
		if j < 2 {
			return 0, ""
		}
		subZql = strings.TrimSpace(subZql[:j])
		return i, strings.Trim(subZql, "'")
	}
	return analyseFieldValue(subZql, field)
}
