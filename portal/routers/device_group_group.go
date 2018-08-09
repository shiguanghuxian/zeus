package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/device_group/types_list", &controllers.DeviceGroupGroupController{}, "get:AjaxGetGroupTypes")          // 设备分组类型列表
	beego.Router("/device_group/ajax_add_one_types", &controllers.DeviceGroupGroupController{}, "post:AjaxAddGroupTypes") // 添加设备分组类型
	beego.Router("/device_group/ajax_del_types", &controllers.DeviceGroupGroupController{}, "get:AjaxDelGroupTypes")      // 删除设备分组类型
	beego.Router("/device_group/ajax_edit_types", &controllers.DeviceGroupGroupController{}, "post:AjaxEditGroupTypes")   // 添加设备分组类型
}
