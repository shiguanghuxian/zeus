package controllers

import (
	"net/url"

	"53it.net/zues/internal"
	"53it.net/zues/models"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

// ajax返回数据
type AjaxData struct {
	State int         `json:"state"` // 非0即为错
	Msg   string      `json:"msg"`   // 提示文本
	Data  interface{} `json:"data"`  // 返回实际需要使用部分数据
}

// 默认执行方法，处理请求体json
func (this *BaseController) Prepare() {
	// 输出服务器信息，伪造
	this.Ctx.Output.Header("X-Powered-By", "PHP/7.0.0")
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")

	// 登录情况
	if this.GetSession("admin") == nil {
		internal.LogFile.W("访问首页，且未登录")
		// 跳转url
		this.Redirect("/login?url="+url.QueryEscape(this.Ctx.Request.RequestURI), 302)
	}

	// 分配用户信息
	if this.GetSession("admin") != nil {
		userInfo := this.getUserInfo()
		this.Data["SessionAdmin"] = userInfo
	} else {
		this.Data["SessionAdmin"] = new(*models.User)
	}
}

// 调用结束要做的事
func (this *BaseController) Finish() {

}

// 获取当前用户登录信息
func (this *BaseController) getUserInfo() *models.User {
	return this.GetSession("admin").(*models.User)
}

// ajax返回数据
func (this *BaseController) AjaxReturn(ajaxData *AjaxData) {
	this.Data["json"] = ajaxData
	this.ServeJSON()
	this.StopRun()
}
