package controllers

import (
	"time"

	"53it.net/zues/models"
)

type UserController struct {
	BaseController
}

// 获取用户信息
func (this *UserController) MyUserInfo() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	user := this.getUserInfo()
	// 查询当前用户信息
	u, err := models.GetIdUserInfo(user.Id)
	if err != nil {
		ajaxData.Msg = "用户信息查询错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: u}
	this.AjaxReturn(ajaxData)
}

// 提交修改个人信息
func (this *UserController) SaveMyUserInfo() {
	ajaxData := &AjaxData{State: 1, Msg: "失败"}
	id, err := this.GetInt("id")
	if err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 用户信息
	info := &models.User{
		Phone:    this.GetString("phone"),
		Email:    this.GetString("email"),
		Language: this.GetString("language"),
	}
	// 判断两次输入密码
	rePassword := this.GetString("repassword")
	password := this.GetString("password")
	// 判断输入密码，且密码小于6位
	if password != "" && len(password) < 6 {
		ajaxData.Msg = "密码长度不符合要求"
		this.AjaxReturn(ajaxData)
	}
	if password != "" && len(password) > 5 {
		if rePassword != password {
			ajaxData.Msg = "两次输入密码不相同"
			this.AjaxReturn(ajaxData)
		} else {
			info.Password = password
		}
	}
	// 修改更新时间
	info.Uptime = time.Now().Unix()
	// 执行修改
	err = models.UpdateIdUserInfo(id, info)
	if err != nil {
		ajaxData.Msg = "修改个人信息失败"
		this.AjaxReturn(ajaxData)
	}
	// 更新session信息-再查询一次库
	go func() {
		u, _ := models.GetIdUserInfo(id)
		this.SetSession("admin", u)
	}()

	ajaxData = &AjaxData{State: 0, Msg: "个人信息保存成功"}
	this.AjaxReturn(ajaxData)
}
