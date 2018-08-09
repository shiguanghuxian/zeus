package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type TopicsConfig struct {
	Id           int32  `json:"id" xorm:"int pk autoincr 'id'"`
	Topics       string `json:"topics" xorm:"varchar(8) 'topics'"`
	Channel      string `json:"channel" xorm:"varchar(20) 'channel'"`
	ChannelCount int32  `json:"channel_count" xorm:"int 'channel_count'"`
	Enable       int32  `json:"enable" xorm:"int 'enable'"`
	DataType     string `json:"data_type" xorm:"char(4) 'data_type'"`
	IsDelete     int32  `json:"is_delete" xorm:"int 'is_delete'"`
}

// 实际表名
func (this *TopicsConfig) TableName() string {
	return "zn_topics_config"
}

// 获取所有配置
func GetAllTopicsConfig() ([]*TopicsConfig, error) {
	list := make([]*TopicsConfig, 0)
	err := dbEngine().Where("(enable = '1') and (is_delete = 0)").Asc("id").Find(&list)
	if err != nil {
		internal.LogFile.E("查询话题配置信息失败:" + err.Error())
		return nil, err
	}
	return list, nil
}

// 获取指定话题的配置信息
func GetWhereTopicsConfig(topics string) ([]*TopicsConfig, error) {
	list := make([]*TopicsConfig, 0)
	err := dbEngine().Where("(`enable` = '1') and (`topics` = '" + topics + "') and (is_delete = 0)").Asc("id").Find(&list)
	if err != nil {
		internal.LogFile.E("查询话题配置信息失败:" + err.Error())
		return nil, err
	}
	return list, nil
}

// 通过关键词查询总数
func GetKeywordTopicsCount(keyword string) (c int64, err error) {
	topicsConfig := new(TopicsConfig)
	if keyword != "" {
		c, err = dbEngine().Where("(topics like '%" + keyword + "%' or channel like '%" + keyword + "%') and (is_delete = 0)").Count(topicsConfig)
	} else {
		c, err = dbEngine().Where("(is_delete = 0)").Count(topicsConfig)
	}
	if err != nil {
		internal.LogFile.W("获取话题列表总数错误:" + err.Error())
	}
	return c, err
}

// 通过关键词查询列表
func GetKeywordTopicsList(page, pageCount int, keyword string) ([]*TopicsConfig, error) {
	pageStart := (page - 1) * pageCount
	list := make([]*TopicsConfig, 0)
	err := dbEngine().Where("(topics like '%"+keyword+"%' or channel like '%"+keyword+"%') and (is_delete = 0)").Asc("id").Limit(pageCount, pageStart).Find(&list)
	if err != nil {
		internal.LogFile.W("获取话题列表错误:" + err.Error())
	}
	return list, err
}

// 修改状态
func UpdateTopicsEnable(id, enable int32) (int64, error) {
	topicsConfig := new(TopicsConfig)
	topicsConfig.Enable = enable
	affected, err := dbEngine().Cols("enable").Where("id = " + strconv.Itoa(int(id))).Update(topicsConfig)
	if err != nil {
		internal.LogFile.E("修改话题状态错误:" + err.Error())
	}
	return affected, err
}

// 添加数据
func AddOneTopicsConfig(topicsConfig *TopicsConfig) (int64, error) {
	affected, err := dbEngine().Insert(topicsConfig)
	if err != nil {
		internal.LogFile.E("添加话题配置错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsTopicsConfig(ids string) (int64, error) {
	topicsConfig := new(TopicsConfig)
	topicsConfig.IsDelete = 1
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + ids + ")").Update(topicsConfig)
	if err != nil {
		internal.LogFile.E("删除话题配置错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneTopicsInfo(id int32) (*TopicsConfig, error) {
	topicsConfig := new(TopicsConfig)
	has, err := dbEngine().Where("id = " + strconv.Itoa(int(id))).Get(topicsConfig)
	if !has {
		internal.LogFile.W("根据id查询话题配置错误:"+err.Error(), has)
		return topicsConfig, errors.New("No query to data")
	}
	return topicsConfig, nil
}

// 根据id修改信息
func UpdateIdTopicsInfo(id int32, topicsConfig *TopicsConfig) (int64, error) {
	affected, err := dbEngine().Cols("topics,channel,channel_count,enable,data_type").Where("id = " + strconv.Itoa(int(id))).Update(topicsConfig)
	if err != nil {
		internal.LogFile.E("根据id修改信息错误:" + err.Error())
	}
	return affected, err
}
