package routers

import (
	"53it.net/zues/portal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.PublicController{}, "get:Login;post:Login")
	beego.Router("/logout", &controllers.PublicController{}, "get:LogOut")
	// 左侧菜单
	beego.Router("/public/leftnavtree", &controllers.PublicController{}, "get:LeftNavTree")
	// 二维码
	beego.Router("/public/verify", &controllers.PublicController{}, "get:VerifyImg")
	beego.Router("/public/verifyshow/*", &controllers.PublicController{}, "get:VerifyShow")
	beego.Router("/public/verifychk", &controllers.PublicController{}, "post:VerifyChk")
}
