package models

import (
	"strconv"

	"53it.net/zues/internal"
)

type EventPush struct {
	Id            int    `json:"id" xorm:"not null pk autoincr INT(10)"`
	EventSetingId int    `json:"event_seting_id" xorm:"index INT(11)"`
	Url           string `json:"url" xorm:"VARCHAR(300)"`
	Name          string `json:"name" xorm:"VARCHAR(60)"`
	DataType      int    `json:"data_type" xorm:"default 0 INT(11)"`
}

// 真实表名
func (this *EventPush) TableName() string {
	return "zn_event_push"
}

// 告警级别列表
func GetAllEventPushList(esid int) ([]EventPush, error) {
	var list []EventPush
	err := dbEngine().Where("(event_seting_id = " + strconv.Itoa(esid) + ")").Find(&list)
	if err != nil {
		internal.LogFile.E("查询告警推送列表:" + err.Error())
		return list, err
	}
	return list, nil
}

// 添加一条
func AddOneEventPush(eventPush *EventPush) (int64, error) {
	affected, err := dbEngine().Insert(eventPush)
	if err != nil {
		internal.LogFile.E("添加 event push 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdEventPush(id string) (int64, error) {
	eventPush := new(EventPush)
	affected, err := dbEngine().Where("id = " + id).Delete(eventPush)
	if err != nil {
		internal.LogFile.E("删除event push错误:"+err.Error(), id)
	}
	return affected, err
}

// 根据id修改信息
func UpdateIdEventPushInfo(id int, eventPush *EventPush) (int64, error) {
	affected, err := dbEngine().
		Cols("name,url,data_type").
		Where("id = " + strconv.Itoa(id)).
		Update(eventPush)
	if err != nil {
		internal.LogFile.E("根据id修改 event push 信息错误:" + err.Error())
	}
	return affected, err
}
