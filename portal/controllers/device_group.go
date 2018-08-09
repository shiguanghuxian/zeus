package controllers

import (
	"fmt"

	"53it.net/zues/models"
)

type DeviceGroupController struct {
	BaseController
}

// ajax获取设备分组列表
func (this *DeviceGroupController) AjaxGroupList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	typeStr := this.GetString("type") // 分组
	isDelete := this.GetString("is_delete")
	// 最终data列表
	var dataList []interface{}
	var err error
	if typeStr == "" {
		ggList, err := models.GetDeviceGroupGroupAll()
		if err == nil {
			for _, v := range ggList {
				list, err := models.GetTypeDeviceGroupList(fmt.Sprint(v.Id), "0", isDelete)
				if err == nil && len(list) > 0 {
					dataList = append(dataList, list)
				}
			}
		}
	} else {
		list, err := models.GetTypeDeviceGroupList(typeStr, "0", isDelete)
		if err == nil && len(list) > 0 {
			dataList = append(dataList, list)
		}
	}
	if err != nil {
		ajaxData.Msg = "查询错误"
		this.AjaxReturn(ajaxData)
	}

	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: dataList}
	this.AjaxReturn(ajaxData)
}

// 添加
func (this *DeviceGroupController) AjaxAddDeviceGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	typeVal, _ := this.GetInt32("type", 0)
	pid, _ := this.GetInt32("parent_id", 0)
	name := this.GetString("name")
	description := this.GetString("description")
	if name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 如果parent_id=0判断改分组的分组也就是type是否存在顶级分组
	if pid == 0 {
		c, err := models.ChkTypeDeviceGroup(typeVal)
		if err != nil || c > 0 {
			ajaxData.Msg = "该类型顶级分组已存在，只能存在一个顶级分组"
			this.AjaxReturn(ajaxData)
		}
	}
	// models对象
	deviceGroup := new(models.DeviceGroup)
	deviceGroup.Name = name
	deviceGroup.Type = typeVal
	deviceGroup.ParentId = pid
	deviceGroup.Description = description
	deviceGroup.IsDelete = 0
	// 执行插入
	_, err := models.AddOneDeviceGroup(deviceGroup)
	if err != nil {
		ajaxData.Msg = "设备分组添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 编辑
func (this *DeviceGroupController) AjaxEditDeviceGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	lastId, _ := this.GetInt("id", 0)
	typeVal, _ := this.GetInt32("type", 0)
	pid, _ := this.GetInt32("parent_id", 0)
	name := this.GetString("name")
	description := this.GetString("description")
	if lastId == 0 || name == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// models对象
	deviceGroup := new(models.DeviceGroup)
	deviceGroup.Name = name
	deviceGroup.Type = typeVal
	deviceGroup.ParentId = pid
	deviceGroup.Description = description
	// 执行插入
	_, err := models.UpOneDeviceGroup(lastId, deviceGroup)
	if err != nil {
		ajaxData.Msg = "设备分组编辑错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *DeviceGroupController) AjaxDelDeviceGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.DelIdsDeviceGroup(id)
	if err != nil {
		ajaxData.Msg = "删除错误"
		this.AjaxReturn(ajaxData)
	}
	// 删除分组关联设备对照列表
	models.DelGidDeviceGroupContrasts(id)

	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

//还原删除
func (this *DeviceGroupController) AjaxRestoreDeviceGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id := this.GetString("id")
	if id == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	_, err := models.RestoreIdsDeviceGroup(id)
	if err != nil {
		ajaxData.Msg = "还原错误"
		this.AjaxReturn(ajaxData)
	}

	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 根据type查询父id列表
func (this *DeviceGroupController) AjaxGetGroupTypeList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	typeStr := this.GetString("type")
	if typeStr == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	parentId := this.GetString("parent_id")
	if parentId == "" {
		parentId = "0"
	}
	list, err := models.GetTypeDeviceGroupList(typeStr, parentId)
	if err != nil {
		ajaxData.Msg = err.Error()
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: list}
	this.AjaxReturn(ajaxData)
}

// 获取分组的分组
func (this *DeviceGroupController) AjaxGetGroupGroupList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ggList, err := models.GetDeviceGroupGroupAll()
	if err != nil {
		ajaxData.Msg = "获取列表错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: ggList}
	this.AjaxReturn(ajaxData)
}

// 移除设备在某个分组
func (this *DeviceGroupController) AjaxRemoveDeviceOnGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	did, _ := this.GetInt32("device_id")
	gid, _ := this.GetInt32("group_id")
	if did == 0 || gid == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 调用修改
	err := models.DelDidGidDeviceGroupContrast(did, gid)
	if err != nil {
		ajaxData.Msg = "移除设备错误"
	} else {
		ajaxData = &AjaxData{State: 0, Msg: "成功"}
	}
	this.AjaxReturn(ajaxData)
}

// 添加设备到分组
func (this *DeviceGroupController) AjaxAddDeviceOnGroup() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	did, _ := this.GetInt32("device_id")
	gid, _ := this.GetInt32("group_id")
	if did == 0 || gid == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 调用修改
	err := models.AddDidGidDeviceGroupContrast(did, gid)
	if err != nil {
		ajaxData.Msg = "添加设备错误"
	} else {
		ajaxData = &AjaxData{State: 0, Msg: "成功"}
	}
	this.AjaxReturn(ajaxData)
}

// 设备设备分组
func (this *DeviceGroupController) AjaxGetDeviceGroupTypeList() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 根分组id
	groupId, _ := this.GetInt32("group_id")
	if groupId == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 根据分组id查询子分组id列表
	gids, err := models.GetGroupIdChildGroupIdList(groupId)
	if err != nil {
		ajaxData.Msg = "查询gids错误"
		this.AjaxReturn(ajaxData)
	}
	gids = append(gids, groupId)
	dgList, err := models.GetGidsDeviceGroupNameList(gids)
	if err != nil {
		ajaxData.Msg = "查询错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: dgList}
	this.AjaxReturn(ajaxData)
}
