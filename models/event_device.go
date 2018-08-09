package models

import "53it.net/zues/internal"
import "fmt"

type EventDevice struct {
	Id            int `json:"id" xorm:"not null pk autoincr INT(11)"`
	EventSetingId int `json:"event_seting_id" xorm:"INT(11) 'event_seting_id'"`
	DeviceId      int `json:"device_id" xorm:"INT(11) 'device_id'"`
}

func (this *EventDevice) TableName() string {
	return "zn_event_device"
}

type EventDeviceInfo struct {
	EventDevice `xorm:"extends"`
	HostName    string `json:"host_name" xorm:"VARCHAR(200) 'hostname'"`
	Ip          string `json:"ip" xorm:"VARCHAR(60) 'ip'"`
	DeviceType  string `json:"device_type" xorm:"VARCHAR(60) 'device_type'"`
	GroupName   string `json:"group_name" xorm:"VARCHAR(60) 'group_name'"`
}

// 根据告警id获取列表
func GetESIDEventDeviceAll(esid int) ([]EventDevice, error) {
	var list []EventDevice
	err := dbEngine().Where(fmt.Sprintf("(event_seting_id = %d)", esid)).OrderBy("id asc").Find(&list)
	if err != nil {
		internal.LogFile.E("查询告警设备列表错误:" + err.Error())
		return list, err
	}
	return list, nil
}

// 添加设备、告警对照
func AddOneEventDevice(eventDevice *EventDevice) (int64, error) {
	affected, err := dbEngine().Insert(eventDevice)
	if err != nil {
		internal.LogFile.E("添加 zn_event_device 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据esid和device_id删除数据
func DelOneEventDevice(esid, deviceId int) (int64, error) {
	eventDevice := new(EventDevice)
	affected, err := dbEngine().Where(fmt.Sprintf("(event_seting_id = %d) and (device_id = %d)", esid, deviceId)).Delete(eventDevice)
	if err != nil {
		internal.LogFile.E("删除 zn_event_device 错误:" + err.Error())
	}
	return affected, err
}

// 根据告警id获取列表
func GetRpcESIDEventDeviceAll(esid int) ([]EventDeviceInfo, error) {
	var list []EventDeviceInfo
	mySql := fmt.Sprintf("SELECT ed.*, d.hostname, d.ip, d.device_type,d.group_name from zn_event_device as ed, zn_device as d WHERE (ed.device_id = d.id) AND (ed.event_seting_id = %d) AND (d.is_delete = 0)", esid)
	err := dbEngine().Sql(mySql).Find(&list)
	if err != nil {
		internal.LogFile.E("rpc查询告警设备列表错误:" + err.Error())
		return list, err
	}
	return list, nil
}
