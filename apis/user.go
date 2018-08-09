package apis

import (
	"errors"
	"net/http"
	"time"

	"53it.net/zues/models"
)

type User struct {
	Apis
}

// 获取用户信息
func (this *User) MyUserInfo(r *http.Request, args *map[string]interface{}, response *Response) error {
	user := this.getUserInfo("")
	// 查询当前用户信息
	u, err := models.GetIdUserInfo(user.Id)
	if err != nil {
		return errors.New("用户信息查询错误")
	}
	*response = u
	return nil
}

// 提交修改个人信息
func (this *User) SaveMyUserInfo(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	// 用户信息
	info := &models.User{
		Phone:    this.ToString((*args)["phone"]),
		Email:    this.ToString((*args)["email"]),
		Language: this.ToString((*args)["language"]),
	}
	// 判断两次输入密码
	rePassword := this.ToString((*args)["repassword"])
	password := this.ToString((*args)["password"])
	// 判断输入密码，且密码小于6位
	if password != "" && len(password) < 6 {
		return errors.New("密码长度不符合要求")
	}
	if password != "" && len(password) > 5 {
		if rePassword != password {
			return errors.New("两次输入密码不相同")
		} else {
			info.Password = password
		}
	}
	// 修改更新时间
	info.Uptime = time.Now().Unix()
	// 执行修改
	err := models.UpdateIdUserInfo(id, info)
	if err != nil {
		return errors.New("修改个人信息失败")
	}
	// 更新session信息-再查询一次库--redis
	go func() {
		u, _ := models.GetIdUserInfo(id)
		this.setSession("admin", u)
	}()

	return nil
}
