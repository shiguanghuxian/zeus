package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "get:Index")
	beego.Router("/home/ajax_zql_query", &controllers.HomeController{}, "get:AjaxZqlQuery")          // 测试
	beego.Router("/home/ajax_zql_list", &controllers.HomeController{}, "get:AjaxZqlList")            // 测试
	beego.Router("/home/ajax_zql_mongo_list", &controllers.HomeController{}, "get:AjaxZqlMongoList") // 测试mongodb
}
