package redis

import (
	"log"
	"os"
	"sync"
	"time"

	"53it.net/zues/internal"
)

var redisDb *Pool
var mutex sync.Mutex

func init() {
	newRedis()
}

func newRedis() {
	mutex.Lock()
	defer mutex.Unlock()
	if internal.CFG == nil {
		return
	}
	// 读取mongodb配置
	address, _ := internal.CFG.String("redis", "address")
	port, _ := internal.CFG.String("redis", "port")
	maxIdle, _ := internal.CFG.Int("redis", "max_idle")
	if maxIdle == 0 {
		maxIdle = 32
	}
	log.Println("正在与redis建立连接")
	// 自定义初始化函数，防止连接错误
	var err error
	redisDb, err = NewPool("tcp", address+":"+port, maxIdle)
	if err != nil {
		internal.LogFile.E("redis创建连接失败:" + err.Error())
		log.Println("redis创建连接失败:" + err.Error())
		log.Println(address + ":" + port)
		os.Exit(2)
	}
	// 防止断开
	pingTime, _ := internal.CFG.String("redis", "ping")
	if pingTime == "" {
		pingTime = "30"
	}
	timeoutT, err := time.ParseDuration(pingTime + "s")
	go func() {
		for {
			redisDb.Cmd("PING")
			time.Sleep(timeoutT)
		}
	}()
	log.Println("已与redis建立连接")
}

// 释放redis资源
func CloseRedis() {
	dbRedis().Empty()
}

// redis操作对象
func dbRedis() *Pool {
	if redisDb == nil {
		log.Println("redis 重连")
		newRedis()
	}
	return redisDb
}
