package controllers

import (
	"strings"

	"53it.net/zues/models"
)

type EventRuleController struct {
	BaseController
}

// 规则列表
func (this *EventRuleController) AjaxRuleList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	esid, err := this.GetInt("esid")
	if err != nil || esid == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	list, err := models.GetSetingIdEventRuleLevelList(esid)
	if err != nil {
		ajaxData.Msg = "查询列表出现错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{Data: list, State: 0, Msg: "列表获取成功"}
	this.AjaxReturn(ajaxData)
}

// 调整排序
func (this *EventRuleController) AjaxEventRuleChageSort() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	rsid, err1 := this.GetInt("rsid")
	sort, err2 := this.GetInt("sort")
	if err1 != nil || err2 != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	err := models.UpEventRuleChageSort(rsid, sort)
	if err != nil {
		ajaxData.Msg = "排序修改错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "排序修改成功"}
	this.AjaxReturn(ajaxData)
}

// 添加告警规则
func (this *EventRuleController) AjaxAddOneEventRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	eventSetingId, err := this.GetInt("event_seting_id", 0)
	value := this.GetString("value")
	if eventSetingId == 0 || err != nil || value == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	eventLevelId, _ := this.GetInt("event_level_id", 0)
	// models对象
	eventRule := new(models.EventRule)
	eventRule.EventLevelId = eventLevelId
	eventRule.EventSetingId = eventSetingId
	eventRule.Value = value
	eventRule.Expression = this.GetString("expression")
	eventRule.Sort, _ = this.GetInt("sort")
	eventRule.Unit = this.GetString("unit")
	// 执行插入
	_, err = models.AddOneEventRule(eventRule)
	if err != nil {
		ajaxData.Msg = "告警规则添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *EventRuleController) AjaxDelEventRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsEventRule(ids)
	if err != nil {
		ajaxData.Msg = "删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取信息
func (this *EventRuleController) AjaxInfoEventRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt("id", 0)
	if id == 0 || err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneEventRuleInfo(id)
	if err != nil {
		ajaxData.Msg = "获取信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 保存编辑信息
func (this *EventRuleController) AjaxUpEventRuleInfo() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, err := this.GetInt("id")
	if err != nil {
		ajaxData.Msg = "参数错误id"
		this.AjaxReturn(ajaxData)
	}
	value := this.GetString("value")
	if value == "" {
		ajaxData.Msg = "参数错误value"
		this.AjaxReturn(ajaxData)
	}
	eventLevelId, _ := this.GetInt("event_level_id", 0)
	// models对象
	eventRule := new(models.EventRule)
	eventRule.EventLevelId = eventLevelId
	eventRule.Value = value
	eventRule.Expression = this.GetString("expression")
	eventRule.Sort, _ = this.GetInt("sort")
	eventRule.Unit = this.GetString("unit")
	// 执行插入
	_, err = models.UpdateIdEventRuleInfo(id, eventRule)
	if err != nil {
		ajaxData.Msg = "告警规则编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
