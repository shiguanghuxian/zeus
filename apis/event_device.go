package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
	"53it.net/zues/redis"
)

type EventDevice struct {
	Apis
}

// 告警设备列表
func (this *EventDevice) EventDeviceList(r *http.Request, args *map[string]interface{}, response *Response) error {
	esid, err := this.ToInt((*args)["esid"], 0)
	if err != nil || esid == 0 {
		return errors.New("参数错误")
	}
	list, err := models.GetESIDEventDeviceAll(esid)
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// 添加设备到告警
func (this *EventDevice) AddEventDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	esid, err := this.ToInt((*args)["esid"], 0)
	if err != nil || esid == 0 {
		return errors.New("参数错误:esid")
	}
	deviceId, err := this.ToInt((*args)["device_id"], 0)
	if err != nil || deviceId == 0 {
		return errors.New("参数错误:device_id")
	}
	/* 查询设备是否有该appname中的字段 */
	// 告警信息
	eventInfo, err := models.GetOneSetingsEventInfo(esid)
	if err != nil {
		return errors.New("查询告警信息错误:" + err.Error())
	}
	// 查询设备信息
	deviceInfo, err := models.GetOneDeviceInfo(deviceId)
	if err != nil {
		return errors.New("查询设备信息错误:" + err.Error())
	}
	rawData := map[string]string{
		"group":    deviceInfo.GroupName,
		"hostname": deviceInfo.HostName,
		"ip":       deviceInfo.Ip,
	}
	kpiInfo, err := redis.GetOneNewestData(rawData, eventInfo.AppName, eventInfo.Field)
	// 判断信息是否存在
	if err != nil || kpiInfo == "" {
		return errors.New("改设备不存改告警字段数据，或为空")
	}
	eventDevice := new(models.EventDevice)
	eventDevice.EventSetingId = esid
	eventDevice.DeviceId = deviceId
	_, err = models.AddOneEventDevice(eventDevice)
	if err != nil {
		return errors.New("添加设备到告警错误")
	}
	return nil
}

// 删除
func (this *EventDevice) DelEventDevice(r *http.Request, args *map[string]interface{}, response *Response) error {
	esid, err := this.ToInt((*args)["esid"], 0)
	if err != nil || esid == 0 {
		return errors.New("参数错误:esid")
	}
	deviceId, err := this.ToInt((*args)["device_id"], 0)
	if err != nil || deviceId == 0 {
		return errors.New("参数错误:device_id")
	}
	_, err = models.DelOneEventDevice(esid, deviceId)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}
