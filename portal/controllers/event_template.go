package controllers

import "53it.net/zues/models"

type SetingsTemplateController struct {
	BaseController
}

// 添加
func (this *SetingsTemplateController) AjaxAddSetingsTemplate() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	sid, err := this.GetInt("seting_event_id", 0)
	if err != nil || sid == 0 {
		ajaxData.Msg = "参数错误sid"
		this.AjaxReturn(ajaxData)
	}
	name := this.GetString("name")
	content := this.GetString("content")
	if name == "" || content == "" {
		ajaxData.Msg = "模板名和内容不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	eventTemplate := new(models.EventTemplate)
	eventTemplate.Name = name
	eventTemplate.Content = content
	// 执行插入
	_, err = models.AddOneEventTemplate(eventTemplate)
	if err != nil {
		ajaxData.Msg = "告警模板添加错误"
		this.AjaxReturn(ajaxData)
	}
	// 修改告警设置的模板id
	_, err = models.UpdateIdSetingsEventTemplateId(sid, int(eventTemplate.Id))
	if err != nil {
		ajaxData.Msg = "模板已添加，保存到告警设置错误，请联系开发者"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取模板信息
func (this *SetingsTemplateController) AjaxInfoTemplate() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, _ := this.GetInt("id", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneEventTemplateInfo(id)
	if err != nil {
		ajaxData.Msg = "获取信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 编辑
func (this *SetingsTemplateController) AjaxUpSetingsTemplate() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, err := this.GetInt("id", 0)
	if err != nil || id == 0 {
		ajaxData.Msg = "参数错误id"
		this.AjaxReturn(ajaxData)
	}
	name := this.GetString("name")
	content := this.GetString("content")
	if name == "" || content == "" {
		ajaxData.Msg = "模板名和内容不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	eventTemplate := new(models.EventTemplate)
	eventTemplate.Name = name
	eventTemplate.Content = content
	// 执行修改
	_, err = models.UpdateIdEventTemplateInfo(id, eventTemplate)
	if err != nil {
		ajaxData.Msg = "告警模板编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
