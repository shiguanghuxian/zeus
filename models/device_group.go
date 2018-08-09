package models

import (
	"strconv"

	"fmt"

	"53it.net/zues/internal"
)

type DeviceGroup struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(60)"`
	ParentId    int32  `json:"parent_id" xorm:"int 'parent_id'"`
	Description string `json:"description" xorm:"TEXT"`
	Type        int32  `json:"type" xorm:"tinyint 'type'"`
	IsDelete    int32  `json:"is_delete" xorm:"tinyint 'is_delete'"`
}

func (this *DeviceGroup) TableName() string {
	return "zn_device_group"
}

// 根据关键词获取总数
func GetKeywordDeviceGroupCount(keyword, typeStr string) (c int64, err error) {
	deviceGroup := new(DeviceGroup)
	where := "(is_delete = 0)"
	if keyword != "" {
		where += fmt.Sprintf(" and (name like '%%%s%%')", keyword)
	} else if typeStr != "" {
		where += fmt.Sprintf(" and (type = %s)", typeStr)
	}
	c, err = dbEngine().Where(where).Count(deviceGroup)
	if err != nil {
		internal.LogFile.W("列表总数错误:" + err.Error())
	}
	return c, err
}

// 通过关键词查询列表
func GetKeywordDeviceGroupList(page, pageCount int, keyword, typeStr string) (list []*DeviceGroup, err error) {
	pageStart := (page - 1) * pageCount
	where := "(is_delete = 0)"
	if keyword != "" {
		where += fmt.Sprintf(" and (name like '%%%s%%')", keyword)
	} else if typeStr != "" {
		where += fmt.Sprintf(" and (type = %s)", typeStr)
	}
	err = dbEngine().Where(where).Asc("id").Limit(pageCount, pageStart).Find(&list)
	if err != nil {
		internal.LogFile.W("DeviceGroup 列表错误:" + err.Error())
	}
	return list, err
}

// 根据id列表获取数据
func GetIdsDeviceGroupList(ids string) (list []*DeviceGroup, err error) {
	err = dbEngine().Where(fmt.Sprintf("id in (%s)", ids)).Find(&list)
	if err != nil {
		internal.LogFile.W("DeviceGroup 根据ids获取列表错误:" + err.Error())
	}
	return list, err
}

// 根据type和pid获取分组列表
func GetTypeAndPidDeviceGroupList(typeStr, pid string, option ...string) (list []*DeviceGroup, err error) {
	whereDelete := "and (is_delete = 0)"
	if len(option) > 0 {
		if option[0] == "true" {
			whereDelete = ""
		}
	}
	err = dbEngine().Where(fmt.Sprintf("(type = '%s') and (parent_id = %s) %s", typeStr, pid, whereDelete)).Find(&list)
	if err != nil {
		internal.LogFile.W("DeviceGroup 根据type获取列表错误:" + err.Error())
	}
	return list, err
}

type DeviceGroupTree struct {
	DeviceGroup
	Child interface{} `json:"child"`
}

// 根据type获取带层级列表
func GetTypeDeviceGroupList(typeStr string, pid string, option ...string) (list []*DeviceGroupTree, err error) {
	if pid == "" {
		pid = "0"
	}
	isDelete := ""
	if len(option) > 0 {
		isDelete = option[0]
	}
	ll, err := GetTypeAndPidDeviceGroupList(typeStr, pid, isDelete)
	if err != nil {
		return nil, err
	}
	for _, v := range ll {
		lll, _ := GetTypeDeviceGroupList(typeStr, fmt.Sprint(v.Id), isDelete)
		v := v
		list = append(list, &DeviceGroupTree{DeviceGroup: *v, Child: lll})
	}
	return list, nil
}

// 添加一条
func AddOneDeviceGroup(deviceGroup *DeviceGroup) (int64, error) {
	affected, err := dbEngine().Insert(deviceGroup)
	if err != nil {
		internal.LogFile.E("添加 deviceGroup 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 编辑
func UpOneDeviceGroup(id int, deviceGroup *DeviceGroup) (int64, error) {
	affected, err := dbEngine().
		Cols("name,description,type,parent_id").
		Where("id = " + strconv.Itoa(id)).
		Update(deviceGroup)
	if err != nil {
		internal.LogFile.E("根据id修改 DeviceGroup 信息错误:" + err.Error())
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsDeviceGroup(id string) (int64, error) {
	deviceGroup := new(DeviceGroup)
	deviceGroup.IsDelete = 1
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + id + ")").Update(deviceGroup)
	if err != nil {
		internal.LogFile.E("删除 DeviceGroup 错误:"+err.Error(), id)
	}
	return affected, err
}

// 根据id列表还原数据
func RestoreIdsDeviceGroup(id string) (int64, error) {
	deviceGroup := new(DeviceGroup)
	deviceGroup.IsDelete = 0
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + id + ")").Update(deviceGroup)
	if err != nil {
		internal.LogFile.E("还原 DeviceGroup 错误:"+err.Error(), id)
	}
	return affected, err
}

// 根据id 真实删除数据
func DelTrueIdDeviceGroup(id string) (int64, error) {
	deviceGroup := new(DeviceGroup)
	affected, err := dbEngine().Where("id in (" + id + ")").Delete(deviceGroup)
	if err != nil {
		internal.LogFile.E("真实删除 DeviceGroup 错误:"+err.Error(), id)
	}
	return affected, err
}

// 根据parent_id 检查是否存在设备
func ChkParentIdDeviceGroup(id string) (c int64, err error) {
	deviceGroup := new(DeviceGroup)
	c, err = dbEngine().Where(fmt.Sprintf("(parent_id = %s)", id)).Count(deviceGroup)
	if err != nil {
		internal.LogFile.E("根据parent_id查询顶级分组 DeviceGroup 错误:" + err.Error())
	}
	return c, err
}

// 根据分组id获取设备id列表
func GetGroupIdChildGroupIdList(gid int32) (gids []int32, err error) {
	deviceGroups := make([]*DeviceGroup, 0)
	err = dbEngine().Cols("id").Where(fmt.Sprintf("(parent_id = %d) and (is_delete = 0)", gid)).Find(&deviceGroups)
	if err != nil {
		return nil, err
	}
	for _, v := range deviceGroups {
		gids = append(gids, int32(v.Id))
		cgids, err := GetGroupIdChildGroupIdList(int32(v.Id))
		if err == nil {
			for _, vv := range cgids {
				gids = append(gids, vv)
			}
		}
	}
	return gids, nil
}

// 获取某个分组的分组下已分组的设备
func GetGroupGroupTypeDeviceIds(t int32) (dids []int32, err error) {
	sql := fmt.Sprintf("SELECT dgc.* FROM zn_device_group_contrast as dgc LEFT JOIN zn_device_group as dg on dgc.device_group_id = dg.id WHERE (dg.type = %d) and (dg.is_delete = 0)", t)
	results, err := engine.Query(sql)
	if err != nil {
		internal.LogFile.E("获取某个分组的分组下已分组的设备 错误:" + err.Error())
		return nil, err
	}
	for _, v := range results {
		vv, _ := strconv.Atoi(string(v["device_id"]))
		if vv != 0 {
			dids = append(dids, int32(vv))
		}
	}
	return dids, nil
}

// 根据type获取是否存在
func ChkTypeDeviceGroup(t int32) (c int64, err error) {
	deviceGroup := new(DeviceGroup)
	c, err = dbEngine().Where(fmt.Sprintf("(is_delete = 0) and (type = %d) and (parent_id = 0)", t)).Count(deviceGroup)
	if err != nil {
		internal.LogFile.E("根据type查询顶级分组 DeviceGroup 错误:" + err.Error())
	}
	return c, err
}

// 根据type查询真实删除的数据
func ChkTrueTypeDeviceGroup(t int32) (c int64, err error) {
	deviceGroup := new(DeviceGroup)
	c, err = dbEngine().Where(fmt.Sprintf("(type = %d)", t)).Count(deviceGroup)
	if err != nil {
		internal.LogFile.E("根据type查询顶级分组 DeviceGroup 错误:" + err.Error())
	}
	return c, err
}
