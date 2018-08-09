package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type TopicsConfigRule struct {
	Id             int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Mapped         string `json:"mapped" xorm:"default '' VARCHAR(255)"`
	TextUnType     string `json:"text_un_type" xorm:"default 'char' ENUM('regular','char')"`
	TextUnRule     string `json:"text_un_rule" xorm:"VARCHAR(255)"`
	DateFormat     string `json:"date_format" xorm:"default '2006-01-02 15:04:05' VARCHAR(60)"`
	TopicsConfigId int    `json:"topics_config_id" xorm:"INT(11)"`
	AppName        string `json:"appname" xorm:"default 'zn_raw_data' VARCHAR(60)"`
	Tag            string `json:"tag" xorm:"default '' VARCHAR(30)"`
	Sort           int    `json:"sort" xorm:"int 'sort'"`
	Enable         int    `json:"enable" xorm:"int 'enable'"`
}

// 实际表名
func (this *TopicsConfigRule) TableName() string {
	return "zn_topics_config_rule"
}

// 根据topics_config_id查询列表
func GetTCIdTopicsConfigRuleList(tc_id int) ([]*TopicsConfigRule, error) {
	list := make([]*TopicsConfigRule, 0)
	err := dbEngine().Where("(topics_config_id = " + strconv.Itoa(tc_id) + ")").Desc("sort").Find(&list)
	if err != nil {
		internal.LogFile.E("查询解析规则列表失败:"+err.Error(), tc_id)
		return nil, err
	}
	return list, nil
}

// 根据topics_config_id查询列表，启用的
func GetTCIdEnableTopicsConfigRuleList(tc_id int) ([]*TopicsConfigRule, error) {
	list := make([]*TopicsConfigRule, 0)
	err := dbEngine().Where(" (enable = 1) and (topics_config_id = " + strconv.Itoa(tc_id) + ")").Desc("sort").Find(&list)
	if err != nil {
		internal.LogFile.E("查询解析规则列表失败:"+err.Error(), tc_id)
		return nil, err
	}
	return list, nil
}

// 修改状态
func UpdateTopicsRuleEnable(id, enable int) (int64, error) {
	topicsConfigRule := new(TopicsConfigRule)
	topicsConfigRule.Enable = enable
	affected, err := dbEngine().Cols("enable").Where("id = " + strconv.Itoa(id)).Update(topicsConfigRule)
	if err != nil {
		internal.LogFile.E("修改解析规则状态错误:" + err.Error())
	}
	return affected, err
}

// 添加解析规则
func AddOneTopicsConfigRule(topicsConfigRule *TopicsConfigRule) (int64, error) {
	affected, err := dbEngine().Insert(topicsConfigRule)
	if err != nil {
		internal.LogFile.E("添加话题解析规则错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsTopicsConfigRule(ids string) (int64, error) {
	topicsConfigRule := new(TopicsConfigRule)
	affected, err := dbEngine().Where("id in (" + ids + ")").Delete(topicsConfigRule)
	if err != nil {
		internal.LogFile.E("删除话题配置错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneTopicsRuleInfo(id int32) (*TopicsConfigRule, error) {
	topicsConfigRule := new(TopicsConfigRule)
	has, err := dbEngine().Where("id = " + strconv.Itoa(int(id))).Get(topicsConfigRule)
	if !has {
		internal.LogFile.W("根据id查询话题解析规则错误:"+err.Error(), has)
		return topicsConfigRule, errors.New("No query to data")
	}
	return topicsConfigRule, nil
}

// 根据id修改信息
func UpdateIdTopicsRuleInfo(id int, topicsConfigRule *TopicsConfigRule) (int64, error) {
	affected, err := dbEngine().Cols("app_name,tag,mapped,text_un_type,text_un_rule,date_format,sort,enable").Where("id = " + strconv.Itoa(int(id))).Update(topicsConfigRule)
	if err != nil {
		internal.LogFile.E("根据id修改信息错误:" + err.Error())
	}
	return affected, err
}
