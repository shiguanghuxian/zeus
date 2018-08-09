package models

import (
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"53it.net/zues/internal"
)

var engine *xorm.Engine
var mutex sync.Mutex

func init() {
	newDb()
}

// 创建数据库连接对象
func newDb() {
	if internal.CFG == nil {
		return
	}
	var err error
	// 读取配置文件
	dbtype, _ := internal.CFG.String("db", "dbtype")
	if dbtype == "" {
		dbtype = "mysql"
	}
	dbuser, _ := internal.CFG.String("db", "dbuser")
	dbpasswd, _ := internal.CFG.String("db", "dbpasswd")
	dbaddress, _ := internal.CFG.String("db", "dbaddress")
	dbport, _ := internal.CFG.String("db", "dbport")
	dbname, _ := internal.CFG.String("db", "dbname")
	debug, eerr := internal.CFG.Bool("db", "debug")
	if eerr != nil {
		debug = false
	}
	// 数据库连接字符串
	dbConnStr := dbuser + ":" + dbpasswd + "@tcp(" + dbaddress + ":" + dbport + ")/" + dbname + "?charset=utf8"
	log.Println("正在与mysql建立连接：", dbConnStr)
	// 连接数据库
	engine, err = xorm.NewEngine(dbtype, dbConnStr)
	if err != nil {
		internal.LogFile.E("创建orm对象失败:" + err.Error())
		log.Println("连接mysql创建orm对象失败:" + err.Error())
		panic(err)
	}
	log.Println("已与mysql建立连接")
	if debug == true {
		engine.Logger().SetLevel(core.LOG_DEBUG) // 调试信息
		engine.ShowSQL(true)                     // 显示sql
	}
	engine.SetMaxIdleConns(10)           // 空闲连接池数量
	engine.SetMaxOpenConns(40)           // 最大连接数
	engine.SetMapper(core.GonicMapper{}) // 命名规则
	// 表前缀
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "zn_")
	engine.SetTableMapper(tbMapper)
	// 设置数据库时区
	location, err := time.LoadLocation("Asia/Shanghai")
	engine.TZLocation = location
}

// 获取数据库操作对象
func dbEngine() *xorm.Engine {
	mutex.Lock()
	defer mutex.Unlock()
	if engine != nil {
		return engine
	}
	newDb()
	return engine
}

// 释放数据库资源
func CloseDb() {
	engine.Close()
}
