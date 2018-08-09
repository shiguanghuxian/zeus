package routers

import (
	"53it.net/zues/portal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/user/myuserinfo", &controllers.UserController{}, "get:MyUserInfo;post:SaveMyUserInfo") // 跟人信息修改显示
}
