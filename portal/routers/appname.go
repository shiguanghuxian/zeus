package routers

import (
	"53it.net/zues/portal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/appname/ajaxlist", &controllers.AppNameController{}, "get:AjaxAppNameList")         // ajax获取列表
	beego.Router("/appname/ajaxaddappname", &controllers.AppNameController{}, "post:AjaxAddAppName")   // ajax添加数据
	beego.Router("/appname/ajaxdelappname", &controllers.AppNameController{}, "get:AjaxDelAppName")    // ajax删除数据
	beego.Router("/appname/ajaxinfoappname", &controllers.AppNameController{}, "get:AjaxInfoAppName")  // ajax获取单条信息
	beego.Router("/appname/ajaxeditappname", &controllers.AppNameController{}, "post:AjaxEditAppName") // ajax提交修改
}
