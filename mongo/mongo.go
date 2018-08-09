package mongo

import (
	"log"
	"os"
	"sync"
	"time"

	"53it.net/zues/internal"
	"gopkg.in/mgo.v2"
)

var mongoConn *mgo.Session
var mgodb *mgo.Database
var mutex sync.Mutex

func init() {
	newMongo()
}

// 创建mongodb链接
func newMongo() {
	defer func() {
		if e := recover(); e != nil {
			internal.LogFile.E("mongodb致命错误", e)
			log.Println("mongodb致命错误", e)
		}
	}()
	if internal.CFG == nil {
		return
	}
	// 是否启用
	enable, _ := internal.CFG.Bool("mongodb", "enable")
	if !enable {
		return
	}
	// 读取mongodb配置
	address, _ := internal.CFG.String("mongodb", "address")
	port, _ := internal.CFG.String("mongodb", "port")
	username, _ := internal.CFG.String("mongodb", "username")
	password, _ := internal.CFG.String("mongodb", "password")
	maxPoolSize, _ := internal.CFG.Int("mongodb", "max_pool_size")
	if maxPoolSize == 0 {
		maxPoolSize = 1024
	}
	log.Println("正在与mongodb建立连接：", address+":"+port)
	var err error
	// mongo数据库配置
	dialInfo, err := mgo.ParseURL(address + ":" + port)
	if err != nil {
		log.Println(err)
	}
	dialInfo.Username = username
	dialInfo.Password = password
	dialInfo.Timeout = 20 * time.Second
	dialInfo.PoolLimit = maxPoolSize
	//连接数据库
	mongoConn, err = mgo.DialWithInfo(dialInfo) //mgo.Dial(address + ":" + port)
	if err != nil {
		internal.LogFile.E("mongodb创建连接失败:" + err.Error())
		log.Println("mongodb创建连接失败:" + err.Error())
		os.Exit(2)
	}
	log.Println("已与mongodb建立连接")
	// 数据插入模式，是否强一致性
	mongoConn.SetMode(mgo.Monotonic, true)
	// 设置超时时间和ping
	setMongodbTimeOut()
}

// 释放mongo资源
func CloseMongoDb() {
	// 是否启用
	enable, _ := internal.CFG.Bool("mongodb", "enable")
	if !enable {
		return
	}
	log.Println("关闭MongoDb")
	mongoConn.Close()
}

// 每月一个表
func getMongodbTableName(tname string) string {
	return tname + "_" + getMongodbSubTableName()
}

// 表子部分
func getMongodbSubTableName(args ...time.Time) string {
	myTime := time.Now()
	if len(args) > 0 {
		myTime = args[0]
	}
	subTname, _ := internal.CFG.String("mongodb", "subtable")
	if subTname == "" {
		subTname = myTime.Format("2006_01")
	} else {
		subTname = myTime.Format(subTname)
	}
	return subTname
}

// 获取mongo session
func mgoSession() *mgo.Session {
	if mongoConn != nil {
		return mongoConn
	}
	// 是否启用
	enable, _ := internal.CFG.Bool("mongodb", "enable")
	if !enable {
		panic("mongodb config enable is false")
	}
	newMongo()
	return mongoConn
}

// mongodb操作对象－－每个通道似乎都会连接掉用一次
func dbmgo() (*mgo.Session, *mgo.Database) {
	// 读取数据库名
	dbname, _ := internal.CFG.String("mongodb", "dbname")
	// session
	sessionOne := mgoSession().Clone()
	return sessionOne, sessionOne.DB(dbname) //数据库名称
}

// 设置超时时间
func setMongodbTimeOut() {
	// ping时间，socket超时时间
	timeout, terr := internal.CFG.String("mongodb", "timeout")
	if terr != nil {
		timeout = "300"
	}
	ping, perr := internal.CFG.String("mongodb", "ping")
	if perr != nil {
		ping = "60"
	}
	timeoutT, err := time.ParseDuration(timeout + "s")
	if err != nil {
		panic("mongodb socket超时配置错误:" + err.Error())
	}
	pingT, err := time.ParseDuration(ping + "s")
	if err != nil {
		panic("mongodb心跳时间配置错误:" + err.Error())
	}
	// 设置超时
	mongoConn.SetSocketTimeout(timeoutT)
	// 定时ping
	go func() {
		defer func() {
			if e := recover(); e != nil {
				internal.LogFile.E("mongodb定时ping失败", e)
				log.Println("mongodb定时ping失败", e)
			}
		}()
		for {
			err := mongoConn.Ping()
			if err != nil {
				newMongo()
			}
			time.Sleep(pingT)
		}
	}()
}
