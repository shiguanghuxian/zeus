package controllers

import "53it.net/zues/models"

type EventLevelController struct {
	BaseController
}

// 级别列表
func (this *EventLevelController) AjaxEventLevelList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	list, err := models.GetAllEventLevelList()
	if err != nil {
		ajaxData.Msg = "查询列表出现错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{Data: list, State: 0, Msg: "列表获取成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *EventLevelController) AjaxDelEventLevel() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.DelIdEventLevel(id)
	if err != nil {
		ajaxData.Msg = "删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 添加
func (this *EventLevelController) AjaxAddEventLevel() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	level, err := this.GetInt("level", 0)
	name := this.GetString("name")
	if err != nil || name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	eventLevel := new(models.EventLevel)
	eventLevel.Name = name
	eventLevel.Level = level
	// 执行插入
	_, err = models.AddOneEventLevel(eventLevel)
	if err != nil {
		ajaxData.Msg = "告警级别添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取信息
func (this *EventLevelController) AjaxInfoEventLevel() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt("id", 0)
	if id == 0 || err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneEventLevelInfo(id)
	if err != nil {
		ajaxData.Msg = "获取信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 编辑
func (this *EventLevelController) AjaxEditEventLevel() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, err := this.GetInt("id")
	if err != nil || id == 0 {
		ajaxData.Msg = "参数id错误"
		this.AjaxReturn(ajaxData)
	}
	level, err := this.GetInt("level", 0)
	name := this.GetString("name")
	if err != nil || name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	eventLevel := new(models.EventLevel)
	eventLevel.Name = name
	eventLevel.Level = level
	// 执行插入
	_, err = models.UpdateIdEventLevelInfo(id, eventLevel)
	if err != nil {
		ajaxData.Msg = "告警级别编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
