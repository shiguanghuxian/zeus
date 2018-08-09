package influxdb

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"

	"fmt"

	"53it.net/influxdb_pool"
	"53it.net/zues/internal"
)

var influxClient client.Client
var influxdbPool *influxdb_pool.InfluxdbPool

var influxDbName string
var mutex sync.Mutex

func init() {
	newInfluxdb()
}

// 创建数据库连接
func newInfluxdb() {
	defer func() {
		if e := recover(); e != nil {
			internal.LogFile.E("influxdb致命错误", e)
			log.Println("influxdb致命错误", e)
		}
	}()
	if internal.CFG == nil {
		return
	}
	// 是否启用
	enable, _ := internal.CFG.Bool("influxdb", "enable")
	if !enable {
		return
	}
	// 数据连接配置
	address, _ := internal.CFG.String("influxdb", "address")
	port, _ := internal.CFG.String("influxdb", "port")
	username, _ := internal.CFG.String("influxdb", "username")
	password, _ := internal.CFG.String("influxdb", "password")
	maxIdle, _ := internal.CFG.Int("influxdb", "max_idle")
	// maxPoolSize, _ := internal.CFG.Int("influxdb", "max_pool_size")
	connType, _ := internal.CFG.Int("influxdb", "conn_type")
	payloadSize, _ := internal.CFG.Int("influxdb", "payload_size")
	timeout, _ := internal.CFG.Int("influxdb", "timeout")
	if maxIdle == 0 {
		maxIdle = 64
	}
	if payloadSize == 0 {
		payloadSize = 512
	}
	if timeout == 0 {
		timeout = 10
	}

	log.Println("创建连接池influxdb：", address+":"+port)
	var err error
	// 创建连接对象
	influxdbPool, err = influxdb_pool.NewInfluxdbPool(influxdb_pool.SetUrl(address, port),
		influxdb_pool.SetAuth(username, password),
		influxdb_pool.SetConnType(connType),
		influxdb_pool.SetMaxPoolIdle(maxIdle),
		influxdb_pool.SetPayloadSize(payloadSize),
		influxdb_pool.SetShowLog(true),
		influxdb_pool.SetTimeout(timeout))
	if err != nil {
		internal.LogFile.E("influxdb创建连接池失败:" + err.Error())
		log.Println("influxdb创建连接池失败:" + err.Error())
		os.Exit(2)
	}
	// 保存数据库名
	influxDbName = getInfluxdbName()
	log.Println("已与influxdb建立连接池")
}

// 关闭数据库连接
func CloseInfluxdb() error {
	log.Println("关闭influxdb")
	return influxClient.Close()
}

// 防止数据库连接断开
func dbInflux() (*influxdb_pool.IdleConn, error) {
	return influxdbPool.Get()
}

// 归还连接
func releaseClient(idleConn *influxdb_pool.IdleConn) {
	influxdbPool.Release(idleConn)
}

// 获取数据库名
func getInfluxdbName() string {
	idleConn, err := dbInflux()
	defer func() {
		releaseClient(idleConn)
		return
	}()
	if influxDbName != "" {
		return influxDbName
	}
	dbname, _ := internal.CFG.String("influxdb", "dbname")
	subTname, _ := internal.CFG.String("influxdb", "subtable")
	if subTname == "" {
		subTname = time.Now().Format("2006")
	} else {
		subTname = time.Now().Format(subTname)
	}
	// 判断数据库是否存在，不存在则创建
	influxDbName = "_internal"
	database, err := QueryResponse("SHOW DATABASES")
	// 数据库名
	influxDbName := dbname + "_" + subTname
	if err != nil {
		internal.LogFile.E("influxdb 查询数据库列表:" + err.Error())
		log.Println("influxdb查询数据库列表:" + err.Error())
		os.Exit(2)
	}
	isDatabase := false
	for _, v := range database {
		if fmt.Sprint(v["name"]) == influxDbName {
			isDatabase = true
		}
	}
	if isDatabase == false {
		q := client.Query{
			Command: fmt.Sprintf("CREATE DATABASE \"%s\"", influxDbName),
		}
		_, err := idleConn.C.Query(q)
		if err != nil {
			internal.LogFile.E("influxdb 创建数据库:" + err.Error())
			log.Println("influxdb 创建数据库:" + err.Error())
			os.Exit(2)
		}
	}
	return influxDbName
}

// 获取表名
func getInfluxdbTableName(tname string, deviceInfo ...string) string {
	// 表名拼接
	for _, v := range deviceInfo {
		tname += "_" + v
	}
	return tname
}
