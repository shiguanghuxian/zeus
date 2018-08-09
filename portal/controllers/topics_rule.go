package controllers

import "53it.net/zues/models"

// 话题设置
type TopicsRuleController struct {
	BaseController
}

// 获取规则列表
func (this *TopicsRuleController) AjaxIdTopicsRules() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt("id", 0)
	if err != nil || id < 1 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	list, err := models.GetTCIdTopicsConfigRuleList(id)
	if err != nil {
		ajaxData.Msg = "查询解析规则列表失败"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功:", Data: list}
	this.AjaxReturn(ajaxData)
}

// 修改状态
func (this *TopicsRuleController) AjaxChangeEnable() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt("id", 0)
	enable, err := this.GetInt("enable", 0)
	if err != nil || id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	if enable == 0 {
		enable = 1
	} else {
		enable = 0
	}
	// 调用修改
	_, err = models.UpdateTopicsRuleEnable(id, enable)
	if err != nil {
		ajaxData.Msg = "修改话题状态错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 添加解析规则
func (this *TopicsRuleController) AddTopicsRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	topicsConfigId, _ := this.GetInt("topics_config_id", 0)
	if topicsConfigId == 0 {
		ajaxData.Msg = "话题配置id错误"
		this.AjaxReturn(ajaxData)
	}
	// 其它参数
	topicsConfigRule := new(models.TopicsConfigRule)
	topicsConfigRule.TopicsConfigId = topicsConfigId
	topicsConfigRule.AppName = this.GetString("app_name")
	topicsConfigRule.Tag = this.GetString("tag")
	topicsConfigRule.Mapped = this.GetString("mapped")
	topicsConfigRule.TextUnType = this.GetString("text_un_type")
	topicsConfigRule.TextUnRule = this.GetString("text_un_rule")
	topicsConfigRule.DateFormat = this.GetString("date_format")
	topicsConfigRule.Sort, _ = this.GetInt("sort", 1)
	topicsConfigRule.Enable, _ = this.GetInt("enable", 1)
	// 调用添加
	_, err := models.AddOneTopicsConfigRule(topicsConfigRule)
	if err != nil {
		ajaxData.Msg = "话题解析规则添加错误"
		this.AjaxReturn(ajaxData)
	}

	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *TopicsRuleController) AjaxDelTopicsRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.DelIdsTopicsConfigRule(id)
	if err != nil {
		ajaxData.Msg = "话题解析规则删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 编辑
func (this *TopicsRuleController) AjaxEditTopicsRuleInfo() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, _ := this.GetInt32("id", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneTopicsRuleInfo(id)
	if err != nil {
		ajaxData.Msg = "获取话题解析规则信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 保存修改
func (this *TopicsRuleController) AjaxUpTopicsRule() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, _ := this.GetInt("id", 0)
	if id == 0 {
		ajaxData.Msg = "话题解析规则id错误"
		this.AjaxReturn(ajaxData)
	}
	// 其它参数
	topicsConfigRule := new(models.TopicsConfigRule)

	topicsConfigRule.AppName = this.GetString("app_name")
	topicsConfigRule.Tag = this.GetString("tag")
	topicsConfigRule.Mapped = this.GetString("mapped")
	topicsConfigRule.TextUnType = this.GetString("text_un_type")
	topicsConfigRule.TextUnRule = this.GetString("text_un_rule")
	topicsConfigRule.DateFormat = this.GetString("date_format")
	topicsConfigRule.Sort, _ = this.GetInt("sort", 1)
	topicsConfigRule.Enable, _ = this.GetInt("enable", 1)
	// 调用添加
	_, err := models.UpdateIdTopicsRuleInfo(id, topicsConfigRule)
	if err != nil {
		ajaxData.Msg = "话题解析规则修改错误"
		this.AjaxReturn(ajaxData)
	}

	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
