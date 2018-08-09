package models

import (
	"errors"
	"fmt"

	"53it.net/zues/internal"
)

type DeviceGroupContrast struct {
	Id            int   `json:"id" xorm:"not null pk autoincr INT(11)"`
	DeviceId      int32 `json:"device_id" xorm:"int 'device_id'"`
	DeviceGroupId int32 `json:"device_group_id" xorm:"int 'device_group_id'"`
}

func (this *DeviceGroupContrast) TableName() string {
	return "zn_device_group_contrast"
}

// 根据分组id列表获取设备id列表
func GetGroupsDeviceGroupContrast(gids []int32) ([]int32, error) {
	deviceGroupContrast := make([]*DeviceGroupContrast, 0)
	gidsStr := internal.IntArrayToString(gids)
	if gidsStr == "" {
		return nil, errors.New("Group ID list is empty")
	}
	err := dbEngine().Cols("device_id").Where(fmt.Sprintf("(device_group_id in (%s))", gidsStr)).Find(&deviceGroupContrast)
	if err != nil {
		internal.LogFile.E("根据分组id列表获取设备id列表 错误:" + err.Error())
		return nil, err
	}
	var ids []int32
	for _, v := range deviceGroupContrast {
		ids = append(ids, v.DeviceId)
	}
	return ids, nil
}

// 根据设备id和分组id删除数据
func DelDidGidDeviceGroupContrast(did, gid int32) error {
	deviceGroupContrast := new(DeviceGroupContrast)
	_, err := dbEngine().Where(fmt.Sprintf("(device_id = %d) and (device_group_id = %d)", did, gid)).Delete(deviceGroupContrast)
	if err != nil {
		internal.LogFile.E("根据设备id和分组id删除数据 错误:" + err.Error())
	}
	return err
}

// 根据设备id删除分组对照数据
func DelGidDeviceGroupContrasts(gid string) error {
	deviceGroupContrast := new(DeviceGroupContrast)
	_, err := dbEngine().Where(fmt.Sprintf("(device_group_id = %s)", gid)).Delete(deviceGroupContrast)
	if err != nil {
		internal.LogFile.E("根据分组id删除数据 错误:" + err.Error())
	}
	return err
}

// 添加设备到本分组
func AddDidGidDeviceGroupContrast(did, gid int32) error {
	deviceGroupContrast := new(DeviceGroupContrast)
	deviceGroupContrast.DeviceId = did
	deviceGroupContrast.DeviceGroupId = gid
	_, err := dbEngine().Insert(deviceGroupContrast)
	if err != nil {
		internal.LogFile.E("添加设备到本分组 错误:" + err.Error())
	}
	return err
}

func GetGidsDeviceGroupNameList(gdis []int32) (l []map[string]string, err error) {
	gidsStr := internal.IntArrayToString(gdis)
	sql := fmt.Sprintf(`SELECT
	dgc.device_id,
	dgc.device_group_id,
	dg. NAME AS device_group_name
FROM
	zn_device_group_contrast AS dgc
LEFT JOIN zn_device_group dg ON dgc.device_group_id = dg.id
WHERE
	(dgc.device_group_id IN (%s))
AND (dg.is_delete = 0)`, gidsStr)
	results, err := engine.Query(sql)
	if err != nil {
		internal.LogFile.E("查询设备id和分组对照 错误:" + err.Error())
		return nil, err
	}
	for _, v := range results {
		val := make(map[string]string, 0)
		for kk, vv := range v {
			val[kk] = string(vv)
		}
		l = append(l, val)
	}
	return l, nil
}
