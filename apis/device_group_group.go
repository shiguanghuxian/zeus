package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
)

type DeviceGroupGroup struct {
	Apis
}

func (this *DeviceGroupGroup) GetGroupTypes(r *http.Request, args *map[string]interface{}, response *Response) error {
	list, err := models.GetDeviceGroupGroupAll()
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// 添加
func (this *DeviceGroupGroup) AddGroupTypes(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	name := this.ToString((*args)["name"])
	if name == "" {
		return errors.New("参数错误")
	}
	description := this.ToString((*args)["description"])
	// models对象
	deviceGroupGroup := new(models.DeviceGroupGroup)
	deviceGroupGroup.Name = name
	deviceGroupGroup.Description = description
	// 执行插入
	_, err := models.AddOneDeviceGroupGroup(deviceGroupGroup)
	if err != nil {
		return errors.New("添加错误")
	}
	return nil
}

// 删除
func (this *DeviceGroupGroup) DelGroupTypes(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 检查是否有分组信息，有则不可以删除
	chkId, _ := this.ToInt32((*args)["id"], 0)
	c, _ := models.ChkTrueTypeDeviceGroup(chkId)
	if c > 0 {
		return errors.New("该类型存在分组，请先删除分组")
	}
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.DelIdDeviceGroupGroup(id)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}

// 编辑
func (this *DeviceGroupGroup) EditGroupTypes(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, err := this.ToInt((*args)["id"], 0)
	if err != nil || id == 0 {
		return errors.New("参数[id]错误")
	}
	name := this.ToString((*args)["name"])
	if name == "" {
		return errors.New("参数[name]错误")
	}
	description := this.ToString((*args)["description"])
	// models对象
	deviceGroupGroup := new(models.DeviceGroupGroup)
	deviceGroupGroup.Name = name
	deviceGroupGroup.Description = description
	// 执行插入
	_, err = models.UpdateIdDeviceGroupGroup(id, deviceGroupGroup)
	if err != nil {
		return errors.New("分组类型编辑错误")
	}
	return nil
}
