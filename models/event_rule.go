package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type EventRule struct {
	Id            int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	EventLevelId  int    `json:"event_level_id" xorm:"index INT(11)"`
	EventSetingId int    `json:"event_seting_id" xorm:"index INT(11)"`
	Value         string `json:"value" xorm:"index VARCHAR(60)"`
	Expression    string `json:"expression" xorm:"default '=' index ENUM('=','>','<','>=','<=','!=')"`
	Sort          int    `json:"sort" xorm:"default 0 index INT(11)"`
	Unit          string `json:"unit" xorm:"VARCHAR(30)"`
}

func (this *EventRule) TableName() string {
	return "zn_event_rule"
}

// 告警规则和告警级别
type EventRuleLevel struct {
	Id            int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	EventLevelId  int    `json:"event_level_id" xorm:"index INT(11)"`
	EventSetingId int    `json:"event_seting_id" xorm:"index INT(11)"`
	Value         string `json:"value" xorm:"index VARCHAR(60)"`
	Expression    string `json:"expression" xorm:"default '=' index ENUM('=','>','<','>=','<=','!=')"`
	Sort          int    `json:"sort" xorm:"default 0 index INT(11)"`
	Name          string `json:"level_name" xorm:"VARCHAR(30)"`
	Level         int    `json:"level_level" xorm:"index INT(11)"`
	Unit          string `json:"unit" xorm:"VARCHAR(30)"`
}

func (this *EventRuleLevel) TableName() string {
	return "zn_event_rule"
}

// 根据告警设置id查询规则列表
func GetSetingIdEventRuleLevelList(esid int) ([]EventRuleLevel, error) {
	eventRuleLevel := make([]EventRuleLevel, 0)
	err := dbEngine().
		Cols("zn_event_rule.*, zn_event_level.name,zn_event_level.level").
		Join("INNER", "zn_event_level", "zn_event_rule.event_level_id = zn_event_level.id").
		Where("(event_seting_id = " + strconv.Itoa(esid) + ")").
		OrderBy("zn_event_rule.sort desc, zn_event_level.level desc, zn_event_rule.id asc").
		Find(&eventRuleLevel)
	if err != nil {
		internal.LogFile.E("查询告警规则错误:" + err.Error())
		return eventRuleLevel, err
	}
	return eventRuleLevel, nil
}

// 调整排序
func UpEventRuleChageSort(rsid int, sort int) error {
	eventRule := new(EventRule)
	eventRule.Sort = sort
	_, err := dbEngine().Cols("sort").Id(rsid).Update(eventRule)
	if err != nil {
		internal.LogFile.E("修改排名错误:" + err.Error())
		return err
	}
	return nil
}

// 添加一条
func AddOneEventRule(eventRule *EventRule) (int64, error) {
	affected, err := dbEngine().Insert(eventRule)
	if err != nil {
		internal.LogFile.E("添加 event rule 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsEventRule(ids string) (int64, error) {
	eventRule := new(EventRule)
	affected, err := dbEngine().Where("id in (" + ids + ")").Delete(eventRule)
	if err != nil {
		internal.LogFile.E("删除event_rule错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneEventRuleInfo(id int) (*EventRule, error) {
	eventRule := new(EventRule)
	has, err := dbEngine().Where("id = " + strconv.Itoa(id)).Get(eventRule)
	if !has {
		internal.LogFile.W("根据id查询 event rule 错误:"+err.Error(), has)
		return eventRule, errors.New("No query to data")
	}
	return eventRule, nil
}

// 根据id修改信息
func UpdateIdEventRuleInfo(id int, eventRule *EventRule) (int64, error) {
	affected, err := dbEngine().
		Cols("event_level_id,value,expression,sort,unit").
		Where("id = " + strconv.Itoa(id)).
		Update(eventRule)
	if err != nil {
		internal.LogFile.E("根据id修改 event rule 信息错误:" + err.Error())
	}
	return affected, err
}
