package controllers

import "53it.net/zues/models"

type DeviceGroupGroupController struct {
	BaseController
}

func (this *DeviceGroupGroupController) AjaxGetGroupTypes() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	list, err := models.GetDeviceGroupGroupAll()
	if err != nil {
		ajaxData.Msg = "查询列表出现错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{Data: list, State: 0, Msg: "列表获取成功"}
	this.AjaxReturn(ajaxData)
}

// 添加
func (this *DeviceGroupGroupController) AjaxAddGroupTypes() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	name := this.GetString("name")
	if name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	deviceGroupGroup := new(models.DeviceGroupGroup)
	deviceGroupGroup.Name = name
	// 执行插入
	_, err := models.AddOneDeviceGroupGroup(deviceGroupGroup)
	if err != nil {
		ajaxData.Msg = "添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *DeviceGroupGroupController) AjaxDelGroupTypes() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.DelIdDeviceGroupGroup(id)
	if err != nil {
		ajaxData.Msg = "删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 编辑
func (this *DeviceGroupGroupController) AjaxEditGroupTypes() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, err := this.GetInt("id")
	if err != nil || id == 0 {
		ajaxData.Msg = "参数id错误"
		this.AjaxReturn(ajaxData)
	}
	name := this.GetString("name")
	if name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	deviceGroupGroup := new(models.DeviceGroupGroup)
	deviceGroupGroup.Name = name
	// 执行插入
	_, err = models.UpdateIdDeviceGroupGroup(id, deviceGroupGroup)
	if err != nil {
		ajaxData.Msg = "分组类型编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}
