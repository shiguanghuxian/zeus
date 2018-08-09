package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/redis"
)

type Device struct {
	Apis
}

// ajax获取设备列表
func (this *Device) DeviceList(r *http.Request, args *map[string]interface{}, response *Response) error {
	var totalRows int64 // 总行数
	var err error

	// 分组id列表
	groupIdStr := ""
	// 是否是not in
	notIn, _ := this.ToInt((*args)["not_group_device"], 0)
	if notIn == 0 {
		//分组参数
		groupId, _ := this.ToInt32((*args)["group_id"], 0)
		if groupId != 0 {
			// 根据分组id查询子分组id列表
			gids, err := models.GetGroupIdChildGroupIdList(groupId)
			if err != nil {
				return errors.New("查询gids错误")
			}
			gids = append(gids, groupId)
			// 查询设备列表
			dids, err := models.GetGroupsDeviceGroupContrast(gids)
			if err == nil {
				groupIdStr = internal.IntArrayToString(dids)
			}
			if groupIdStr == "" {
				return errors.New("该分组下不存在设备")
			}
		}
		// 分组参数，id列表形式
		groupIds := this.ToString((*args)["group_ids"])
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
				return errors.New("该分组列表下不存在设备")
			}
		}
	} else if notIn == 1 {
		groupGroupId, err := this.ToInt32((*args)["group_group_id"], 0)
		if err != nil || (err == nil && groupGroupId == 0) {
			return errors.New("参数错误：group_group_id")
		}
		dids, err := models.GetGroupGroupTypeDeviceIds(groupGroupId)
		if err == nil {
			groupIdStr = internal.IntArrayToString(dids)
		}
	}

	// 主机信息筛选
	hostname := this.ToString((*args)["host_name"])     // 主机名
	deviceType := this.ToString((*args)["device_type"]) // 主机类型
	groupName := this.ToString((*args)["group_name"])   // 主机分组
	sort := this.ToString((*args)["sort"])              // 排序字段
	sortType, _ := this.ToInt((*args)["sort_type"], 0)  // 排序方向
	isDelete, _ := this.ToInt((*args)["is_delete"], 0)  // 排序方向
	totalRows, err = models.GetKeywordDeviceCount(hostname, deviceType, groupName, groupIdStr, isDelete, notIn)
	if err != nil {
		return errors.New("服务端错误 count")
	}
	// 每页行数
	listRows, err := internal.CFG.Int("apis", "pagecount")
	if err != nil {
		listRows = 10
	}
	// 当前页码
	page, _ := this.ToInt((*args)["page"], 1)
	// 查询列表
	list, err := models.GetKeywordDeviceList(hostname, deviceType, groupName, groupIdStr, sort, isDelete, page, listRows, sortType, notIn)
	if err != nil {
		return errors.New("服务端错误 list")
	}
	// 页面数据
	data := make(map[string]interface{})
	data["page"] = page
	data["total_rows"] = totalRows
	data["list_rows"] = listRows
	data["list"] = list

	*response = data
	return nil
}

// 删除设备
func (this *Device) DelIdsDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsDevice(ids)
	if err != nil {
		return errors.New("设备删除错误")
	}
	return nil
}

// 还原删除的设备
func (this *Device) RestoreIdsDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.RestoreIdsDevice(ids)
	if err != nil {
		return errors.New("设备还原错误")
	}
	return nil
}

// 保存描述信息和删除与否
func (this *Device) UpDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt((*args)["id"], 0)
	description := this.ToString((*args)["description"])
	isDelete, _ := this.ToInt32((*args)["is_delete"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	// 话题models对象
	device := new(models.Device)
	device.Description = description
	device.IsDelete = isDelete
	// 执行修改
	_, err := models.UpdateIdDeviceInfo(id, device)
	if err != nil {
		return errors.New("设备信息修改错误")
	}
	return nil
}

// AutoDiscoveryDevice 设备发现页面
func (this *Device) AutoDiscoveryDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 获取是否只查看未保存
	isSave, _ := this.ToBool((*args)["is_save"])
	var list []string
	var err error
	if isSave {
		list, err = redis.GetKeysList("devicelist:0:")
	} else {
		list, err = redis.GetKeysList("devicelist:")
	}
	if err != nil {
		return err
	}
	// 关键词
	keyword := this.ToString((*args)["keyword"])
	// 当前页码
	page, _ := this.ToInt((*args)["page"], 1)
	// 每页行数
	listRows, err := internal.CFG.Int("apis", "pagecount")
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

	*response = data
	return nil
}

func (this *Device) SaveOneDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	deviceInfo := this.ToString((*args)["device"])
	description := this.ToString((*args)["description"])
	if deviceInfo == "" {
		return errors.New("参数不能为空")
	}
	dL := strings.Split(deviceInfo, ":")
	// 判断是否存在id，存在则已保存
	if dL[1] != "0" {
		return errors.New("设备已保存，如果设备列表不存在，可在已删除列表中恢复")
	}
	if len(dL) != 6 {
		return errors.New("参数错误")
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
		return errors.New("设备信息保存错误")
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
		return errors.New(msg + ",但更新redis失败")
	}
	// 删除原来的
	deviceMap["id"] = 0
	redis.DelDeviceInfo(deviceMap)

	return nil
}

// 同步设备id列表到redis
func (this *Device) SynchroDeviceids(r *http.Request, args *map[string]interface{}, response *Response) error {
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
		return errors.New("查询设备列表为空")
	}
	for _, v := range dList {
		rkey := fmt.Sprintf("deviceids:%s:%s:%s", v.GroupName, v.HostName, v.Ip)
		err := redis.SetKeyVal(rkey, fmt.Sprint(v.Id))
		tCount++
		if err == nil {
			okCount++
		}
	}

	*response = &map[string]string{"message": fmt.Sprintf("共[%d]个设备，成功同步[%d]个设备id", tCount, okCount)}
	return nil
}

// 数据上传分组列表
func (this *Device) DeviceNativeGroupList(r *http.Request, args *map[string]interface{}, response *Response) error {
	keys, err := redis.GetKeysList("deviceids:")
	if err != nil {
		return errors.New("列表查询为空")
	}
	gList := make(map[string]string, 0)
	for _, v := range keys {
		vv := strings.Split(v, ":")
		if len(vv) < 2 {
			continue
		}
		gList[vv[1]] = vv[1]
	}

	*response = gList
	return nil
}
