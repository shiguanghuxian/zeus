package models

import (
	"strconv"

	"53it.net/zues/internal"
)

type DeviceGroupGroup struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(60)"`
	Description string `json:"description" xorm:"text"`
}

func (this *DeviceGroupGroup) TableName() string {
	return "zn_device_group_group"
}

// 获取所有分组的分组
func GetDeviceGroupGroupAll() (lists []*DeviceGroupGroup, err error) {
	err = dbEngine().Find(&lists)
	if err != nil {
		internal.LogFile.W("DeviceGroupGroup 列表错误:" + err.Error())
	}
	return lists, err
}

// 添加一条
func AddOneDeviceGroupGroup(deviceGroupGroup *DeviceGroupGroup) (int64, error) {
	affected, err := dbEngine().Insert(deviceGroupGroup)
	if err != nil {
		internal.LogFile.E("添加 分组类型 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdDeviceGroupGroup(id string) (int64, error) {
	deviceGroupGroup := new(DeviceGroupGroup)
	affected, err := dbEngine().Where("id = " + id).Delete(deviceGroupGroup)
	if err != nil {
		internal.LogFile.E("删除 分组类型 错误:"+err.Error(), id)
	}
	return affected, err
}

// 根据id修改信息
func UpdateIdDeviceGroupGroup(id int, deviceGroupGroup *DeviceGroupGroup) (int64, error) {
	affected, err := dbEngine().
		Cols("name, description").
		Where("id = " + strconv.Itoa(id)).
		Update(deviceGroupGroup)
	if err != nil {
		internal.LogFile.E("根据id修改 分组类型 信息错误:" + err.Error())
	}
	return affected, err
}
