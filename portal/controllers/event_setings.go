package controllers

import (
	"strings"

	"53it.net/zues/models"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

type SetingsController struct {
	BaseController
}

// 告警列表
func (this *SetingsController) AjaxEventList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}

	var totalRows int64 // 总行数
	var err error

	keyword := this.GetString("keyword") // 关键次参数
	totalRows, err = models.GetKeywordEventSetingCount(keyword)
	if err != nil {
		ajaxData.Msg = "服务端错误 count"
		this.AjaxReturn(ajaxData)
	}
	// 每页行数
	listRows, err := beego.AppConfig.Int("AdminListPagesCount")
	if err != nil {
		listRows = 10
	}
	// 当前页码
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	// 查询列表
	list, err := models.GetKeywordEventSetingList(page, listRows, keyword)
	if err != nil {
		ajaxData.Msg = "服务端错误 list"
		this.AjaxReturn(ajaxData)
	}
	// 页面数据
	data := make(map[string]interface{})
	data["page"] = page
	data["total_rows"] = totalRows
	data["list_rows"] = listRows
	data["list"] = list
	ajaxData = &AjaxData{State: 0, Msg: "数据获取成功", Data: data}

	// 返回数据
	this.AjaxReturn(ajaxData)
}

// 修改启用状态
func (this *SetingsController) AjaxChangeEnable() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt32("id", 0)
	enable, err := this.GetInt32("enable", 0)
	if err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	if enable == 0 {
		enable = 1
	} else {
		enable = 0
	}
	// 调用修改
	_, err = models.UpdateEventSetingEnable(id, enable)
	if err != nil {
		ajaxData.Msg = "修改告警配置状态错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 添加告警配置
func (this *SetingsController) AjaxAddSetingEvent() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	name := this.GetString("name")
	appName := this.GetString("app_name")
	field := this.GetString("field")
	if name == "" || appName == "" || field == "" {
		ajaxData.Msg = "告警标题、APPNAM和字段名 不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	eventSeting := new(models.EventSeting)
	eventSeting.Name = name
	eventSeting.AppName = appName
	eventSeting.Field = field
	eventSeting.ValueType = this.GetString("value_type")
	// 数值型判断
	continuedCount, err := this.GetInt("continued_count", 0)
	if err != nil || continuedCount == 0 {
		ajaxData.Msg = "出现次数输入错误"
		this.AjaxReturn(ajaxData)
	}
	eventSeting.ContinuedCount = continuedCount
	// 步长时间
	continuedTime, err := this.GetInt("continued_time", 0)
	if err != nil || continuedTime == 0 {
		ajaxData.Msg = "步长时间输入错误"
		this.AjaxReturn(ajaxData)
	}
	eventSeting.ContinuedTime = continuedTime
	// 执行周期
	cycleTime := this.GetString("cycle_time")
	_, err = cron.Parse(cycleTime)
	if err != nil {
		ajaxData.Msg = "执行周期输入错误:" + err.Error()
		this.AjaxReturn(ajaxData)
	}
	eventSeting.CycleTime = cycleTime
	//  启用情况
	enable, _ := this.GetInt32("enable", 0)
	eventSeting.Enable = enable
	// 描述
	eventSeting.Describe = this.GetString("describe")

	// 执行插入
	_, err = models.AddOneEventSeting(eventSeting)
	if err != nil {
		ajaxData.Msg = "告警配置添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除告警设置
func (this *SetingsController) AjaxDelSetingsEvent() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsEventSeting(ids)
	if err != nil {
		ajaxData.Msg = "告警配置删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取单条消息
func (this *SetingsController) AjaxInfoSetingsEvent() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, _ := this.GetInt("id", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneSetingsEventInfo(id)
	if err != nil {
		ajaxData.Msg = "获取setings_event信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 保存编辑告警
func (this *SetingsController) AjaxUpSetingEvent() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, err := this.GetInt("id", 0)
	if err != nil || id == 0 {
		ajaxData.Msg = "ID参数错误"
		this.AjaxReturn(ajaxData)
	}
	name := this.GetString("name")
	appName := this.GetString("app_name")
	field := this.GetString("field")
	if name == "" || appName == "" || field == "" {
		ajaxData.Msg = "告警标题、APPNAM和字段名 不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	eventSeting := new(models.EventSeting)
	eventSeting.Name = name
	eventSeting.AppName = appName
	eventSeting.Field = field
	eventSeting.ValueType = this.GetString("value_type")
	// 数值型判断
	continuedCount, err := this.GetInt("continued_count", 0)
	if err != nil || continuedCount == 0 {
		ajaxData.Msg = "出现次数输入错误"
		this.AjaxReturn(ajaxData)
	}
	eventSeting.ContinuedCount = continuedCount
	// 步长时间
	continuedTime, err := this.GetInt("continued_time", 0)
	if err != nil || continuedTime == 0 {
		ajaxData.Msg = "步长时间输入错误"
		this.AjaxReturn(ajaxData)
	}
	eventSeting.ContinuedTime = continuedTime
	// 执行周期，定时执行
	cycleTime := this.GetString("cycle_time")
	_, err = cron.Parse(cycleTime)
	if err != nil {
		ajaxData.Msg = "执行周期输入错误:" + err.Error()
		this.AjaxReturn(ajaxData)
	}
	eventSeting.CycleTime = cycleTime

	//  启用情况
	enable, _ := this.GetInt32("enable", 0)
	eventSeting.Enable = enable
	// 描述
	eventSeting.Describe = this.GetString("describe")

	// 执行插入
	_, err = models.UpdateIdSetingsEventInfo(id, eventSeting)
	if err != nil {
		ajaxData.Msg = "告警配置编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
