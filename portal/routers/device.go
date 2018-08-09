package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/device/ajax_device_list", &controllers.DeviceController{}, "get:AjaxDeviceList")             // 设备列表
	beego.Router("/device/ajax_del_device", &controllers.DeviceController{}, "get:AjaxDelIdsDevice")            // 删除设备列表
	beego.Router("/device/ajax_restore_device", &controllers.DeviceController{}, "get:AjaxRestoreIdsDevice")    // 还原删除设备列表
	beego.Router("/device/ajax_up_device", &controllers.DeviceController{}, "post:AjaxUpDevice")                // 修改设备信息
	beego.Router("/device/ajax_synchro_deviceids", &controllers.DeviceController{}, "get:AjaxSynchroDeviceids") // 同步设备ids列表

	beego.Router("/device/ajax_auto_discovery_list", &controllers.DeviceController{}, "get:AjaxAutoDiscoveryDevice") // 设备发现列表
	beego.Router("/device/ajax_save_one_device", &controllers.DeviceController{}, "post:AjaxSaveOneDevice")          // 保存自动发现的设备

	beego.Router("/device/ajax_get_device_native_group_list", &controllers.DeviceController{}, "get:AjaxGetDeviceNativeGroupList") // 设备原始分组列表
}
