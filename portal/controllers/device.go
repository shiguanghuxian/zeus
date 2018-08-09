package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/redis"
	"github.com/astaxie/beego"
)

type DeviceController struct {
	BaseController
}

// ajax获取设备列表
func (this *DeviceController) AjaxDeviceList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}

	var totalRows int64 // 总行数
	var err error

	// 分组id列表
	groupIdStr := ""
	// 是否是not in
	notIn, _ := this.GetInt("not_group_device")
	if notIn == 0 {
		//分组参数
		groupId, _ := this.GetInt32("group_id")
		if groupId != 0 {
			// 根据分组id查询子分组id列表
			gids, err := models.GetGroupIdChildGroupIdList(groupId)
			if err != nil {
				ajaxData.Msg = "查询gids错误"
				this.AjaxReturn(ajaxData)
			}
			gids = append(gids, groupId)
			// 查询设备列表
			dids, err := models.GetGroupsDeviceGroupContrast(gids)
			if err == nil {
				groupIdStr = internal.IntArrayToString(dids)
			}
			if groupIdStr == "" {
				ajaxData.Msg = "该分组下不存在设备"
				this.AjaxReturn(ajaxData)
			}
		}
		// 分组参数，id列表形式
		groupIds := this.GetString("group_ids")
		if groupIds != "" {
			gids := strings.Split(groupIds, ",")
			gidss := make([]int32, 0)
			for _, v := range gids {
				vv, err := strconv.Atoi(v)
				if err == nil {
					gidss = append(gidss, int32(vv))
				}
			}
			dids, err := models.GetGroupsDeviceGroupContrast(gidss)
			if err == nil {
				groupIdStr = internal.IntArrayToString(dids)
			}
			if groupIdStr == "" {
				ajaxData.Msg = "该分组列表下不存在设备"
				this.AjaxReturn(ajaxData)
			}
		}
	} else if notIn == 1 {
		groupGroupId, err := this.GetInt32("group_group_id")
		if err != nil || (err == nil && groupGroupId == 0) {
			ajaxData.Msg = "参数错误：group_group_id"
			this.AjaxReturn(ajaxData)
		}
		dids, err := models.GetGroupGroupTypeDeviceIds(groupGroupId)
		if err == nil {
			groupIdStr = internal.IntArrayToString(dids)
		}
	}

	// 主机信息筛选
	hostname := this.GetString("host_name")     // 主机名
	deviceType := this.GetString("device_type") // 主机类型
	groupName := this.GetString("group_name")   // 主机分组
	sort := this.GetString("sort")              // 排序字段
	sortType, _ := this.GetInt("sort_type")     // 排序方向
	isDelete, _ := this.GetInt("is_delete")     // 排序方向
	totalRows, err = models.GetKeywordDeviceCount(hostname, deviceType, groupName, groupIdStr, isDelete, notIn)
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
	list, err := models.GetKeywordDeviceList(hostname, deviceType, groupName, groupIdStr, sort, isDelete, page, listRows, sortType, notIn)
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

// 删除设备
func (this *DeviceController) AjaxDelIdsDevice() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsDevice(ids)
	if err != nil {
		ajaxData.Msg = "设备删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 还原删除的设备
func (this *DeviceController) AjaxRestoreIdsDevice() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.RestoreIdsDevice(ids)
	if err != nil {
		ajaxData.Msg = "设备还原错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 保存描述信息和删除与否
func (this *DeviceController) AjaxUpDevice() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, _ := this.GetInt("id", 0)
	description := this.GetString("description")
	isDelete, _ := this.GetInt32("is_delete", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	device := new(models.Device)
	device.Description = description
	device.IsDelete = isDelete
	// 执行修改
	_, err := models.UpdateIdDeviceInfo(id, device)
	if err != nil {
		ajaxData.Msg = "设备信息修改错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// AutoDiscoveryDevice 设备发现页面
func (this *DeviceController) AjaxAutoDiscoveryDevice() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	// 获取是否只查看未保存
	isSave, _ := this.GetBool("is_save")
	var list []string
	var err error
	if isSave {
		list, err = redis.GetKeysList("devicelist:0:")
	} else {
		list, err = redis.GetKeysList("devicelist:")
	}
	if err != nil {
		ajaxData.Msg = err.Error()
		this.AjaxReturn(ajaxData)
	}
	// 关键词
	keyword := this.GetString("keyword")
	// 当前页码
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	// 每页行数
	listRows, err := beego.AppConfig.Int("AdminListPagesCount")
	if err != nil {
		listRows = 10
	}
	// 处理数据
	var list1 []string
	var deviceList []map[string]string
	if keyword != "" {
		for _, v := range list {
			// 查找关键词
			if ok := strings.Index(v, keyword); ok < 1 {
				continue
			}
			list1 = append(list1, v)
		}
	} else {
		list1 = list
	}
	totalRows := len(list1)
	// 切割数组，要现实的页
	pageEnd := page * listRows
	if totalRows < pageEnd {
		pageEnd = totalRows
	}
	list1 = list1[(page-1)*listRows : pageEnd]
	for _, v := range list1 {
		// 切割数组
		dL := strings.Split(v, ":")
		if len(dL) != 6 {
			continue
		}
		deviceList = append(deviceList, map[string]string{
			"key":         v,
			"id":          dL[1],
			"group":       dL[2],
			"device_type": dL[3],
			"hostname":    dL[4],
			"ip":          dL[5],
		})
	}
	// 页面数据
	data := make(map[string]interface{})
	data["page"] = page
	data["total_rows"] = totalRows
	data["list_rows"] = listRows
	data["list"] = deviceList
	ajaxData = &AjaxData{State: 0, Msg: "数据获取成功", Data: data}
	this.AjaxReturn(ajaxData)
}

func (this *DeviceController) AjaxSaveOneDevice() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	deviceInfo := this.GetString("device")
	description := this.GetString("description")
	if deviceInfo == "" {
		ajaxData.Msg = "参数不能为空"
		this.AjaxReturn(ajaxData)
	}
	dL := strings.Split(deviceInfo, ":")
	// 判断是否存在id，存在则已保存
	if dL[1] != "0" {
		ajaxData.Msg = "设备已保存，如果设备列表不存在，可在已删除列表中恢复"
		this.AjaxReturn(ajaxData)
	}
	if len(dL) != 6 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	device := new(models.Device)
	device.IsDelete = 0
	device.HostName = dL[4]
	device.Ip = dL[5]
	device.DeviceType = dL[3]
	device.GroupName = dL[2]
	device.Description = description
	// 执行修改
	_, err := models.AddOneDeviceInfo(device)
	if err != nil {
		ajaxData.Msg = "设备信息保存错误"
		this.AjaxReturn(ajaxData)
	}
	msg := "设备保存成功" // 返回信息
	// 更新redis缓存
	deviceMap := map[string]interface{}{
		"id":          device.Id,
		"group":       dL[2],
		"device_type": dL[3],
		"hostname":    dL[4],
		"ip":          dL[5],
		"description": description,
	}
	err = redis.SaveDeviceInfo(deviceMap)
	if err != nil {
		ajaxData.Msg = msg + ",但更新redis失败"
		this.AjaxReturn(ajaxData)
	}
	// 删除原来的
	deviceMap["id"] = 0
	redis.DelDeviceInfo(deviceMap)

	ajaxData = &AjaxData{State: 0, Msg: msg}
	this.AjaxReturn(ajaxData)
}

// 同步设备id列表到redis
func (this *DeviceController) AjaxSynchroDeviceids() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	// 删除列表
	keyList, _ := redis.GetKeysList("deviceids:*")
	for _, vv := range keyList {
		redis.DelKeyVal(vv)
	}
	// 统计错误
	okCount := 0
	tCount := 0
	//添加数据
	dList, err := models.GetAllDeviceLists()
	if err != nil {
		ajaxData.Msg = "查询设备列表为空"
		this.AjaxReturn(ajaxData)
	}
	for _, v := range dList {
		rkey := fmt.Sprintf("deviceids:%s:%s:%s", v.GroupName, v.HostName, v.Ip)
		err := redis.SetKeyVal(rkey, fmt.Sprint(v.Id))
		tCount++
		if err == nil {
			okCount++
		}
	}
	ajaxData = &AjaxData{State: 0, Msg: fmt.Sprintf("共[%d]个设备，成功同步[%d]个设备id", tCount, okCount)}
	this.AjaxReturn(ajaxData)
}

// 数据上传分组列表
func (this *DeviceController) AjaxGetDeviceNativeGroupList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	keys, err := redis.GetKeysList("deviceids:")
	if err != nil {
		ajaxData.Msg = "列表查询为空"
		this.AjaxReturn(ajaxData)
	}
	gList := make(map[string]string, 0)
	for _, v := range keys {
		vv := strings.Split(v, ":")
		if len(vv) < 2 {
			continue
		}
		gList[vv[1]] = vv[1]
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: gList}
	this.AjaxReturn(ajaxData)
}
