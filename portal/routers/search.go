package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/search/v1/query", &controllers.SearchController{}, "get:AjaxZqlQueryV1")    // zql查询列表
	beego.Router("/search/data_source", &controllers.SearchController{}, "get:AjaxDataSource") // 数据源类型
}
