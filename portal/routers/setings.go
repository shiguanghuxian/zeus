package routers

import (
	"53it.net/zues/portal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// setings event
	beego.Router("/setings/ajaxevent", &controllers.SetingsController{}, "get:AjaxEventList")                      // ajax获取列表
	beego.Router("/setings/ajaxchangeenable", &controllers.SetingsController{}, "get:AjaxChangeEnable")            // ajax 修改是否启用
	beego.Router("/setings/ajaxaddsetingevent", &controllers.SetingsController{}, "post:AjaxAddSetingEvent")       // ajax 添加
	beego.Router("/setings/ajaxdelsetingevent", &controllers.SetingsController{}, "get:AjaxDelSetingsEvent")       // ajax 删除
	beego.Router("/setings/ajax_info_setings_event", &controllers.SetingsController{}, "get:AjaxInfoSetingsEvent") // ajax 获取用于编辑
	beego.Router("/setings/ajax_up_setings_event", &controllers.SetingsController{}, "post:AjaxUpSetingEvent")     // ajax 获取用于编辑

	// setings event template
	beego.Router("/setings/ajax_add_event_template", &controllers.SetingsTemplateController{}, "post:AjaxAddSetingsTemplate") // ajax 添加
	beego.Router("/setings/ajax_info_template", &controllers.SetingsTemplateController{}, "get:AjaxInfoTemplate")             // ajax 获取信息
	beego.Router("/setings/ajax_up_event_template", &controllers.SetingsTemplateController{}, "post:AjaxUpSetingsTemplate")   // ajax 保存信息

	// setings event rule
	beego.Router("/setings/ajax_rule_list", &controllers.EventRuleController{}, "get:AjaxRuleList")                       // ajax 告警规则列表
	beego.Router("/setings/ajax_event_rule_chage_sort", &controllers.EventRuleController{}, "get:AjaxEventRuleChageSort") // ajax 更新排序
	beego.Router("/setings/ajax_add_one_event_rule", &controllers.EventRuleController{}, "post:AjaxAddOneEventRule")      // ajax 添加一条
	beego.Router("/setings/ajax_del_event_rule", &controllers.EventRuleController{}, "get:AjaxDelEventRule")              // ajax 删除
	beego.Router("/setings/ajax_info_event_rule", &controllers.EventRuleController{}, "get:AjaxInfoEventRule")            // ajax 获取信息
	beego.Router("/setings/ajax_up_event_rule_info", &controllers.EventRuleController{}, "post:AjaxUpEventRuleInfo")      // ajax 编辑一条信息

	// setings event level
	beego.Router("/setings/ajax_event_level_list", &controllers.EventLevelController{}, "get:AjaxEventLevelList")  // ajax 告警级别列表
	beego.Router("/setings/ajax_del_event_level", &controllers.EventLevelController{}, "get:AjaxDelEventLevel")    // ajax 删除
	beego.Router("/setings/ajax_add_event_level", &controllers.EventLevelController{}, "post:AjaxAddEventLevel")   // ajax 添加
	beego.Router("/setings/ajax_info_event_level", &controllers.EventLevelController{}, "get:AjaxInfoEventLevel")  // ajax 获取单条信息
	beego.Router("/setings/ajax_edit_event_level", &controllers.EventLevelController{}, "post:AjaxEditEventLevel") // ajax 修改保存

	// setings event push
	beego.Router("/setings/ajax_event_push_list", &controllers.EventPushController{}, "get:AjaxEventPushList") // ajax 告警级别列表
	beego.Router("/setings/ajax_add_event_push", &controllers.EventPushController{}, "post:AjaxAddEventPush")  // ajax 添加
	beego.Router("/setings/ajax_del_event_push", &controllers.EventPushController{}, "get:AjaxDelEventPush")   // ajax 添加
	beego.Router("/setings/ajax_up_event_push", &controllers.EventPushController{}, "post:AjaxUpEventPush")    // ajax 添加
}
