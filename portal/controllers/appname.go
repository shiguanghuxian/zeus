package controllers

import (
	"strings"

	"53it.net/zues/models"
	"github.com/astaxie/beego"
)

type AppNameController struct {
	BaseController
}

// appname 列表
func (this *AppNameController) AjaxAppNameList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}

	var totalRows int64 // 总行数
	var err error

	keyword := this.GetString("keyword") // 关键次参数
	totalRows, err = models.GetKeywordAppNameCount(keyword)
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
	list, err := models.GetKeywordAppNameList(page, listRows, keyword)
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

// 添加appname字段配置
func (this *AppNameController) AjaxAddAppName() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	appname := this.GetString("app_name")
	field := this.GetString("field")
	if appname == "" || field == "" {
		ajaxData.Msg = "应用名和字段不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	appNameFieldType := new(models.AppnameFieldType)
	appNameFieldType.AppName = appname
	appNameFieldType.Field = field
	// 其它参数
	appNameFieldType.Type = this.GetString("type")
	appNameFieldType.Unit = this.GetString("unit")
	appNameFieldType.Index, _ = this.GetInt32("index", 0)
	// 执行插入
	_, err := models.AddOneAppName(appNameFieldType)
	if err != nil {
		ajaxData.Msg = "AppName添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *AppNameController) AjaxDelAppName() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsAppName(ids)
	if err != nil {
		ajaxData.Msg = "AppName删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取单条消息
func (this *AppNameController) AjaxInfoAppName() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, _ := this.GetInt32("id", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneAppNameInfo(id)
	if err != nil {
		ajaxData.Msg = "获取AppName信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 编辑appname字段配置
func (this *AppNameController) AjaxEditAppName() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, _ := this.GetInt32("id", 0)
	appname := this.GetString("app_name")
	field := this.GetString("field")
	if id == 0 || appname == "" || field == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	appNameFieldType := new(models.AppnameFieldType)
	appNameFieldType.AppName = appname
	appNameFieldType.Field = field
	// 其它参数
	appNameFieldType.Type = this.GetString("type")
	appNameFieldType.Unit = this.GetString("unit")
	appNameFieldType.Index, _ = this.GetInt32("index", 0)
	// 执行插入
	_, err := models.UpdateIdAppNameInfo(id, appNameFieldType)
	if err != nil {
		ajaxData.Msg = "AppName编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
