package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/redis"
)

type Public struct {
	Apis
}

// 登录
func (this *Public) Login(r *http.Request, args *map[string]interface{}, response *Response) error {
	username := this.ToString((*args)["username"])
	password := this.ToString((*args)["password"])
	u, err := models.VerifyLogin(username, password)
	if err != nil {
		return errors.New("用户名或密码")
	}
	// 写入redis
	token := internal.Rand().Hex()
	err = redis.SetSessionAdmin(token, u)
	if err != nil {
		return errors.New("糟糕，服务器出现错误，请联系开发者")
	}
	*response = map[string]interface{}{"user": u, "token": token}
	return nil
}

// 退出
func (this *Public) LogOut(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 退出类型pc or phone
	device := this.ToString((*args)["device"])
	token := this.ToString((*args)["token"])
	if device == "" {
		return errors.New("参数错误")
	}
	// 清除session
	var err error
	if device == "admin" {
		err = redis.DelSessionAdmin(token)
	} else if device == "phone" {
		err = redis.DelSessionPhone(token)
	} else {
		return errors.New("参数错误")
	}
	return err
}

// 递归菜单
func (this *Public) LeftNavTree(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 查询数据
	tree, err := models.GetMenuTree(0, "0")
	if err == nil {
		*response = tree
		return nil
	}
	return errors.New("左侧菜单获取错误")
}

// 级别列表
func (this *Public) RpcServerList(r *http.Request, args *map[string]interface{}, response *Response) error {
	list, err := models.GetAllRpcServerList()
	if err != nil {
		errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}
