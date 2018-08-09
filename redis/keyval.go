package redis

import (
	"errors"
	"fmt"
	"log"
	"time"

	"53it.net/zues/internal"
)

// 根据key保存值
func SetKeyVal(key, val string) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	resp := dbRedis().Cmd("SET", key, val)
	// 添加情况
	if resp.Err != nil {
		internal.LogFile.E("redis 根据key保存val失败:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}

// 根据key获取val
func GetKeyVal(key string) (val string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	val, err = dbRedis().Cmd("GET", key).Str()
	if err != nil {
		internal.LogFile.E("redis 根据key读取数据失败:" + err.Error())
		return "", err
	}
	return val, nil
}

// 根据key 删除值
func DelKeyVal(key string) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	resp := dbRedis().Cmd("DEL", key)
	// 添加情况
	if resp.Err != nil {
		internal.LogFile.E("redis 根据key删除val失败:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}
