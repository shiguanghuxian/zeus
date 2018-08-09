package controllers

import "53it.net/zues/models"

type EventPushController struct {
	BaseController
}

// 级别列表
func (this *EventPushController) AjaxEventPushList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	esid, err := this.GetInt("esid")
	if esid == 0 && err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	list, err := models.GetAllEventPushList(esid)
	if err != nil {
		ajaxData.Msg = "查询列表出现错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{Data: list, State: 0, Msg: "列表获取成功"}
	this.AjaxReturn(ajaxData)
}

// AjaxAddEventPush 添加
func (this *EventPushController) AjaxAddEventPush() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	name := this.GetString("name")
	purl := this.GetString("url")
	esid, _ := this.GetInt("event_seting_id", 0)
	dataType, _ := this.GetInt("data_type", 0)
	if esid == 0 || name == "" || purl == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	eventPush := new(models.EventPush)
	eventPush.Name = name
	eventPush.Url = purl
	eventPush.DataType = dataType
	eventPush.EventSetingId = esid
	// 执行插入
	_, err := models.AddOneEventPush(eventPush)
	if err != nil {
		ajaxData.Msg = "告警推送添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *EventPushController) AjaxDelEventPush() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.DelIdEventPush(id)
	if err != nil {
		ajaxData.Msg = "删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// AjaxUpEventPush 编辑
func (this *EventPushController) AjaxUpEventPush() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	pId, _ := this.GetInt("id", 0)
	name := this.GetString("name")
	purl := this.GetString("url")
	dataType, _ := this.GetInt("data_type", 0)
	if pId == 0 || name == "" || purl == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	eventPush := new(models.EventPush)
	eventPush.Name = name
	eventPush.Url = purl
	eventPush.DataType = dataType
	// 执行插入
	_, err := models.UpdateIdEventPushInfo(pId, eventPush)
	if err != nil {
		ajaxData.Msg = "告警推送编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
