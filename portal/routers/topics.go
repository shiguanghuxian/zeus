package routers

import (
	"53it.net/zues/portal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/topics/ajaxtopicslist", &controllers.TopicsController{}, "get:AjaxTopicsList")         // ajax获取话题列表
	beego.Router("/topics/ajaxchangeenable", &controllers.TopicsController{}, "get:AjaxChangeEnable")     // ajax修改状态
	beego.Router("/topics/ajaxaddtopics", &controllers.TopicsController{}, "post:AjaxAddTopics")          // ajax添加
	beego.Router("/topics/ajaxdeltopics", &controllers.TopicsController{}, "get:AjaxDelTopics")           // ajax删除
	beego.Router("/topics/ajaxinfotopics", &controllers.TopicsController{}, "get:AjaxInfoTopics")         // ajax获取信息
	beego.Router("/topics/ajaxuptopics", &controllers.TopicsController{}, "post:AjaxUpTopics")            // ajax保存修改信息
	beego.Router("/topics/ajaxrestartserverd", &controllers.TopicsController{}, "get:AjaxRestartServerd") // ajax同步话题配置，重启服务

	beego.Router("/topics/nsqtopics", &controllers.TopicsController{}, "get:NsqTopics") // nsq所有话题列表
}
