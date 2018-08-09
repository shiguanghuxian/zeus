package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type EventTemplate struct {
	Id      int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name    string `json:"name" xorm:"index VARCHAR(30)"`
	Content string `json:"content" xorm:"TEXT"`
}

func (this *EventTemplate) TableName() string {
	return "zn_event_template"
}

// 告警级别列表
func GetAllEventTemplateList() ([]EventTemplate, error) {
	var list []EventTemplate
	err := dbEngine().OrderBy("id asc").Find(&list)
	if err != nil {
		internal.LogFile.E("查询告警级别列表:" + err.Error())
		return list, err
	}
	return list, nil
}

// 添加数据
func AddOneEventTemplate(eventTemplate *EventTemplate) (int64, error) {
	affected, err := dbEngine().Insert(eventTemplate)
	if err != nil {
		internal.LogFile.E("添加告警模板错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneEventTemplateInfo(id int) (*EventTemplate, error) {
	eventTemplate := new(EventTemplate)
	has, err := dbEngine().Where("id = " + strconv.Itoa(id)).Get(eventTemplate)
	if !has {
		internal.LogFile.W("根据id查询event_template错误:"+err.Error(), has)
		return eventTemplate, errors.New("No query to data")
	}
	return eventTemplate, nil
}

// 根据id修改信息
func UpdateIdEventTemplateInfo(id int, eventTemplate *EventTemplate) (int64, error) {
	affected, err := dbEngine().Cols("name,content").Where("id = " + strconv.Itoa(id)).Update(eventTemplate)
	if err != nil {
		internal.LogFile.E("根据id修改EventTemplate信息错误:" + err.Error())
	}
	return affected, err
}
