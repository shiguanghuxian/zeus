package routers

import (
	"53it.net/zues/portal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/topics_rule/ajaxidtopicsrules", &controllers.TopicsRuleController{}, "get:AjaxIdTopicsRules")           // ajax获取解析规则列表
	beego.Router("/topics_rule/ajaxchangeenable", &controllers.TopicsRuleController{}, "get:AjaxChangeEnable")             // ajax修改状态
	beego.Router("/topics_rule/addtopicsrule", &controllers.TopicsRuleController{}, "post:AddTopicsRule")                  // ajax添加解析规则
	beego.Router("/topics_rule/ajaxdeltopicsrule", &controllers.TopicsRuleController{}, "get:AjaxDelTopicsRule")           // ajax删除
	beego.Router("/topics_rule/ajaxedittopicsruleinfo", &controllers.TopicsRuleController{}, "get:AjaxEditTopicsRuleInfo") // ajax获取信息
	beego.Router("/topics_rule/ajaxuptopicsrule", &controllers.TopicsRuleController{}, "post:AjaxUpTopicsRule")            // ajax获取信息
}
