package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/device_group/ajax_list", &controllers.DeviceGroupController{}, "get:AjaxGroupList")                                // ajax 设备分组列表
	beego.Router("/device_group/ajax_add_device_group", &controllers.DeviceGroupController{}, "post:AjaxAddDeviceGroup")              // ajax 添加设备分组
	beego.Router("/device_group/ajax_edit_device_group", &controllers.DeviceGroupController{}, "post:AjaxEditDeviceGroup")            // ajax 编辑设备分组
	beego.Router("/device_group/ajax_del_device_group", &controllers.DeviceGroupController{}, "get:AjaxDelDeviceGroup")               // ajax 删除设备分组
	beego.Router("/device_group/ajax_restore_device_group", &controllers.DeviceGroupController{}, "get:AjaxRestoreDeviceGroup")       // ajax 还原设备分组
	beego.Router("/device_group/ajax_get_device_group_type_list", &controllers.DeviceGroupController{}, "get:AjaxGetGroupTypeList")   // ajax 根据type查询分组列表
	beego.Router("/device_group/ajax_get_device_group_group_list", &controllers.DeviceGroupController{}, "get:AjaxGetGroupGroupList") // ajax 分组的分组，type

	beego.Router("/device_group/ajax_remove_device_on_group", &controllers.DeviceGroupController{}, "get:AjaxRemoveDeviceOnGroup") // ajax 移除设备
	beego.Router("/device_group/ajax_add_device_on_group", &controllers.DeviceGroupController{}, "get:AjaxAddDeviceOnGroup")       // ajax 添加设备到分组

	beego.Router("/device_group/ajax_get_device_groupname_type_list", &controllers.DeviceGroupController{}, "get:AjaxGetDeviceGroupTypeList") // ajax 设备id和设备名对照

}
