package controllers

import (
	"image/png"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"github.com/astaxie/beego"
)

type PublicController struct {
	beego.Controller
}

// 默认执行方法，处理请求体json
func (this *PublicController) Prepare() {
	// 输出服务器信息，伪造
	this.Ctx.Output.Header("X-Powered-By", "PHP/7.0.0")
	this.Data["Lang"] = beego.AppConfig.String("lang::default") // 未登陆设置默认语言
}

// 登录
func (this *PublicController) Login() {
	if this.GetString("username", "") == "" {
		this.SetSession("admin", nil) // 只要访问登录页面，则删除session
	} else {
		username := this.GetString("username", "")
		password := this.GetString("password", "")
		u, err := models.VerifyLogin(username, password)
		this.Data["LoginError"] = false
		if err != nil && u.Id == 0 {
			internal.LogFile.E("登录失败:")
			this.Data["LoginError"] = true
		} else {
			// 存储session
			this.SetSession("admin", u)
			// 获取跳转地址,如果为空则跳转到首页
			url := this.GetString("url")
			if url == "" {
				url = "/"
			}
			this.Redirect(url, 302)
			return
		}
	}
	// 页面隐藏文本框赋值使用
	url := this.GetString("url")
	this.Data["Rurl"] = url
	this.TplName = "public/login.html"
}

// 退出
func (this *PublicController) LogOut() {
	// 清除session
	this.SetSession("admin", nil)
	this.Redirect("/login", 302)
}

// 递归菜单
func (this *PublicController) LeftNavTree() {
	// ajax返回数据
	data := AjaxData{State: 1, Msg: "左侧菜单获取失败，请联系管理员"}
	// 查询数据
	tree, err := models.GetMenuBeegoTree(0, "0")
	if err == nil {
		data.State = 0
		data.Msg = "左侧菜单获取成功"
		data.Data = tree
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 验证码
func (this *PublicController) VerifyImg() {
	captcha, _ := internal.GetCaptcha()
	key, _ := captcha.GetKey(4)
	this.Data["Key"] = key

	this.TplName = "public/code.html"
}

// 显示图片
func (this *PublicController) VerifyShow() {
	key := this.Ctx.Input.Param(":splat")
	captcha, _ := internal.GetCaptcha()
	img, _ := captcha.GetImage(key)

	this.Ctx.Output.Header("Content-Type", "image/png")
	png.Encode(this.Ctx.ResponseWriter, img)
}

// 验证二维码
func (this *PublicController) VerifyChk() {
	captcha, _ := internal.GetCaptcha()

	captchaId := this.GetString("captchaId")
	captchaSolution := this.GetString("captchaSolution")
	cc, ss := captcha.Verify(captchaId, captchaSolution)
	if cc {
		this.Data["Key"] = "ok" + ss
	} else {
		this.Data["Key"] = "no" + ss
	}

	this.TplName = "public/code.html"
}
