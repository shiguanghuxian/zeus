package apis

import (
	"errors"
	"fmt"
	"net/http"

	"53it.net/zues/models"
)

type DeviceGroup struct {
	Apis
}

// ajax获取设备分组列表
func (this *DeviceGroup) GroupList(r *http.Request, args *map[string]interface{}, response *Response) error {
	typeStr := this.ToString((*args)["type"]) // 分组
	isDelete := this.ToString((*args)["is_delete"])
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
		return err
	}

	*response = dataList
	return nil
}

// 添加
func (this *DeviceGroup) AddDeviceGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	typeVal, _ := this.ToInt32((*args)["type"], 0)
	pid, _ := this.ToInt32((*args)["parent_id"], 0)
	name := this.ToString((*args)["name"])
	description := this.ToString((*args)["description"])
	if name == "" {
		return errors.New("参数错误")
	}
	// 如果parent_id=0判断改分组的分组也就是type是否存在顶级分组
	if pid == 0 {
		c, err := models.ChkTypeDeviceGroup(typeVal)
		if err != nil || c > 0 {
			return errors.New("该类型顶级分组已存在，只能存在一个顶级分组")
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
		return errors.New("设备分组添加错误")
	}
	return nil
}

// 编辑
func (this *DeviceGroup) EditDeviceGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	lastId, _ := this.ToInt((*args)["id"], 0)
	typeVal, _ := this.ToInt32((*args)["type"], 0)
	pid, _ := this.ToInt32((*args)["parent_id"], 0)
	name := this.ToString((*args)["name"])
	description := this.ToString((*args)["description"])
	if lastId == 0 || name == "" {
		return errors.New("参数错误")
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
		return errors.New("设备分组编辑错误")
	}
	return nil
}

// 删除
func (this *DeviceGroup) DelDeviceGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.DelIdsDeviceGroup(id)
	if err != nil {
		return errors.New("删除错误")
	}
	// 删除分组关联设备对照列表
	models.DelGidDeviceGroupContrasts(id)

	return nil
}

//还原删除
func (this *DeviceGroup) RestoreDeviceGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.RestoreIdsDeviceGroup(id)
	if err != nil {
		return errors.New("还原错误")
	}

	return nil
}

// 真实删除
func (this *DeviceGroup) DelTrueDeviceGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	// 检查子级是否存在，存在则提示
	c, _ := models.ChkParentIdDeviceGroup(id)
	if c > 0 {
		return errors.New("此分组包含子分组，请先删除子分组")
	}
	_, err := models.DelTrueIdDeviceGroup(id)
	if err != nil {
		return errors.New("删除错误")
	}

	return nil
}

// 根据type查询父id列表
func (this *DeviceGroup) GetGroupTypeList(r *http.Request, args *map[string]interface{}, response *Response) error {
	typeStr := this.ToString((*args)["type"])
	if typeStr == "" {
		return errors.New("参数错误")
	}
	parentId := this.ToString((*args)["parent_id"])
	if parentId == "" {
		parentId = "0"
	}
	list, err := models.GetTypeDeviceGroupList(typeStr, parentId)
	if err != nil {
		return err
	}
	*response = list
	return nil
}

// 获取分组的分组
func (this *DeviceGroup) GetGroupGroupList(r *http.Request, args *map[string]interface{}, response *Response) error {
	ggList, err := models.GetDeviceGroupGroupAll()
	if err != nil {
		return errors.New("获取列表错误")
	}
	*response = ggList
	return nil
}

// 移除设备在某个分组
func (this *DeviceGroup) RemoveDeviceOnGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	did, _ := this.ToInt32((*args)["device_id"])
	gid, _ := this.ToInt32((*args)["group_id"])
	if did == 0 || gid == 0 {
		return errors.New("参数错误")
	}
	// 调用修改
	err := models.DelDidGidDeviceGroupContrast(did, gid)
	if err != nil {
		return errors.New("移除设备错误")
	}
	return nil
}

// 添加设备到分组
func (this *DeviceGroup) AddDeviceOnGroup(r *http.Request, args *map[string]interface{}, response *Response) error {
	did, _ := this.ToInt32((*args)["device_id"])
	gid, _ := this.ToInt32((*args)["group_id"])
	if did == 0 || gid == 0 {
		return errors.New("参数错误")
	}
	// 调用修改
	err := models.AddDidGidDeviceGroupContrast(did, gid)
	if err != nil {
		return errors.New("添加设备错误")
	}
	return nil
}

// 设备设备分组
func (this *DeviceGroup) GetDeviceGroupTypeList(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 根分组id
	groupId, _ := this.ToInt32((*args)["group_id"])
	if groupId == 0 {
		return errors.New("参数错误")
	}
	// 根据分组id查询子分组id列表
	gids, err := models.GetGroupIdChildGroupIdList(groupId)
	if err != nil {
		return errors.New("查询gids错误")
	}
	gids = append(gids, groupId)
	dgList, err := models.GetGidsDeviceGroupNameList(gids)
	if err != nil {
		return errors.New("查询错误")
	}
	*response = dgList
	return nil
}
