package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"53it.net/zues/internal"
	"53it.net/zues/models"
)

// 保存session pc端管理页面
func SetSessionAdmin(token string, u *models.User) error {
	return setSession("session:admin:"+token, u)
}

// 保存session 手机app
func SetSessionPhone(token string, u *models.User) error {
	return setSession("session:phone:"+token, u)
}

// 保存session
func setSession(key string, u *models.User) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	// 设置登录时间
	mObj := make(map[string]interface{}, 0)
	mObj["user"], _ = json.Marshal(u)
	mObj["login_time"] = time.Now().Unix() // session写入时间戳
	resp := dbRedis().Cmd("HMSET", key, mObj)
	if resp.Err != nil {
		internal.LogFile.E("session保存错误:" + resp.Err.Error())
		return resp.Err
	}
	return expireSession(key)
}

// 设置过期时间
func expireSession(key string) error {
	expire, err := internal.CFG.Int("apis", "session_expire")
	if err != nil || expire == 0 {
		expire = 1440
	}
	resp := dbRedis().Cmd("EXPIRE", key, expire)
	if resp.Err != nil {
		internal.LogFile.E("session设置过期时间错误:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}

// 获取session pc端管理页面
func GetSessionAdmin(token string) (u *models.User, loginTime string, err error) {
	return getSession("session:admin:" + token)
}

// 获取session 手机app
func GetSessionPhone(token string) (u *models.User, loginTime string, err error) {
	return getSession("session:phone:" + token)
}

// 获取session
func getSession(key string) (u *models.User, loginTime string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	// 用户信息
	userBytes, err := dbRedis().Cmd("HGET", key, "user").Bytes()
	if err != nil || len(userBytes) == 0 {
		internal.LogFile.E("session读取用户信息错误1:" + err.Error())
		return nil, "", err
	}
	err = json.Unmarshal(userBytes, &u)
	if err != nil {
		internal.LogFile.E("session读取用户信息错误2:" + err.Error())
		return nil, "", err
	}
	// 最后登录时间
	loginTime, err = dbRedis().Cmd("HGET", key, "login_time").Str()
	if err != nil {
		internal.LogFile.E("session读取最后登录时间错误:" + err.Error())
		loginTime = ""
	}
	return u, loginTime, expireSession(key)
}

// 删除session pc端管理页面
func DelSessionAdmin(token string) error {
	return delSession("session:admin:" + token)
}

// 删除session 手机端
func DelSessionPhone(token string) error {
	return delSession("session:phone:" + token)
}

// 删除session
func delSession(key string) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Println(err1)
			internal.LogFile.E("致命错误：" + time.Now().Format("2006-01-02 15:04:05"))
			err = errors.New(fmt.Sprint(err1))
		}
		return
	}()
	// 最后登录时间
	resp := dbRedis().Cmd("DEL", key)
	if resp.Err != nil {
		internal.LogFile.E("session删除错误:" + resp.Err.Error())
		return resp.Err
	}
	return nil
}
