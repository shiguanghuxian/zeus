package models

import (
	"errors"
	"fmt"
	"strconv"

	"53it.net/zues/internal"
)

type Device struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	HostName    string `json:"host_name" xorm:"VARCHAR(200) 'hostname'"`
	Ip          string `json:"ip" xorm:"VARCHAR(60) 'ip'"`
	DeviceType  string `json:"device_type" xorm:"VARCHAR(60) 'device_type'"`
	GroupName   string `json:"group_name" xorm:"VARCHAR(60) 'group_name'"`
	Description string `json:"description" xorm:"VARCHAR(500) 'description'"`
	IsDelete    int32  `json:"is_delete" xorm:"int 'is_delete'"`
}

func (this *Device) TableName() string {
	return "zn_device"
}

// GetKeywordDeviceCount 根据关键词获取总数
func GetKeywordDeviceCount(hostname, deviceType, group, dids string, isDelete, notIn int) (c int64, err error) {
	device := new(Device)
	where := "(is_delete = " + strconv.Itoa(isDelete) + ")"
	if dids != "" {
		if notIn == 0 {
			where += fmt.Sprintf(" and (id in (%s))", dids)
		} else {
			where += fmt.Sprintf(" and (id not in (%s))", dids)
		}
	}
	if hostname != "" {
		where += " and (hostname like '%" + hostname + "%')"
	}
	if deviceType != "" {
		where += " and (device_type = '" + deviceType + "')"
	}
	if group != "" {
		where += " and (group_name = '" + group + "')"
	}
	c, err = dbEngine().Where(where).Count(device)
	if err != nil {
		internal.LogFile.W("列表总数错误:" + err.Error())
	}
	return c, err
}

// GetKeywordDeviceList 根据关键词查询设备列表
func GetKeywordDeviceList(hostname, deviceType, group, dids, sort string, isDelete, page, pageCount, sortType, notIn int) (devices []Device, err error) {
	if sort == "" {
		sort = "id"
	}
	pageStart := (page - 1) * pageCount
	where := "(is_delete = " + strconv.Itoa(isDelete) + ")"
	if dids != "" {
		if notIn == 0 {
			where += fmt.Sprintf(" and (id in (%s))", dids)
		} else {
			where += fmt.Sprintf(" and (id not in (%s))", dids)
		}
	}
	if hostname != "" {
		where += " and (hostname like '%" + hostname + "%')"
	}
	if deviceType != "" {
		where += " and (device_type = '" + deviceType + "')"
	}
	if group != "" {
		where += " and (group_name = '" + group + "')"
	}
	if sortType == 0 {
		err = dbEngine().Where(where).Asc(sort).Limit(pageCount, pageStart).Find(&devices)
	} else {
		err = dbEngine().Where(where).Desc(sort).Limit(pageCount, pageStart).Find(&devices)
	}
	if err != nil {
		internal.LogFile.W("查询设备列表数据失败：" + err.Error())
	}
	return devices, err
}

// 根据id列表删除数据
func DelIdsDevice(ids string) (int64, error) {
	device := new(Device)
	device.IsDelete = 1
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + ids + ")").Update(device)
	if err != nil {
		internal.LogFile.E("删除设备错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id列表还原删除数据
func RestoreIdsDevice(ids string) (int64, error) {
	device := new(Device)
	device.IsDelete = 0
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + ids + ")").Update(device)
	if err != nil {
		internal.LogFile.E("还原删除设备错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id修改信息
func UpdateIdDeviceInfo(id int, device *Device) (int64, error) {
	affected, err := dbEngine().Cols("description,is_delete").Where("id = " + strconv.Itoa(id)).Update(device)
	if err != nil {
		internal.LogFile.E("根据id修改信息错误:" + err.Error())
	}
	return affected, err
}

// 添加一条
func AddOneDeviceInfo(device *Device) (int64, error) {
	affected, err := dbEngine().Insert(device)
	if err != nil {
		internal.LogFile.E("添加 device 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 获取所有设备
func GetAllDeviceLists() (devices []Device, err error) {
	err = dbEngine().Find(&devices)
	if err != nil {
		internal.LogFile.E("获取所有设备:" + err.Error())
	}
	return devices, err
}

// 根据id获取一条设备信息
func GetOneDeviceInfo(id int) (*Device, error) {
	device := new(Device)
	ok, _ := dbEngine().Where(fmt.Sprintf("(id = %d)", id)).Get(device)
	if ok {
		return device, nil
	}
	return nil, errors.New("未查询到设备信息")
}
