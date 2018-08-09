package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type EventLevel struct {
	Id    int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name  string `json:"name" xorm:"VARCHAR(30)"`
	Level int    `json:"level" xorm:"index INT(11)"`
}

func (this *EventLevel) TableName() string {
	return "zn_event_level"
}

// 告警级别列表
func GetAllEventLevelList() ([]EventLevel, error) {
	var list []EventLevel
	err := dbEngine().OrderBy("level desc").Find(&list)
	if err != nil {
		internal.LogFile.E("查询告警级别列表:" + err.Error())
		return list, err
	}
	return list, nil
}

// 根据id列表删除数据
func DelIdEventLevel(id string) (int64, error) {
	eventLevel := new(EventLevel)
	affected, err := dbEngine().Where("id = " + id).Delete(eventLevel)
	if err != nil {
		internal.LogFile.E("删除event level错误:"+err.Error(), id)
	}
	return affected, err
}

// 添加一条
func AddOneEventLevel(eventLevel *EventLevel) (int64, error) {
	affected, err := dbEngine().Insert(eventLevel)
	if err != nil {
		internal.LogFile.E("添加 event rule 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneEventLevelInfo(id int) (*EventLevel, error) {
	eventLevel := new(EventLevel)
	has, err := dbEngine().Where("id = " + strconv.Itoa(id)).Get(eventLevel)
	if !has {
		internal.LogFile.W("根据id查询 event level 错误:"+err.Error(), has)
		return eventLevel, errors.New("No query to data")
	}
	return eventLevel, nil
}

// 根据id修改信息
func UpdateIdEventLevelInfo(id int, eventLevel *EventLevel) (int64, error) {
	affected, err := dbEngine().
		Cols("name,level").
		Where("id = " + strconv.Itoa(id)).
		Update(eventLevel)
	if err != nil {
		internal.LogFile.E("根据id修改 event level 信息错误:" + err.Error())
	}
	return affected, err
}
