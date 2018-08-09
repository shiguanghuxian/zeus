package models

import (
	"errors"
	"strconv"

	"53it.net/zues/internal"
)

type EventSeting struct {
	Id        int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"VARCHAR(60)"`
	AppName   string `json:"app_name" xorm:"VARCHAR(60)"`
	Field     string `json:"field" xorm:"VARCHAR(30)"`
	ValueType string `json:"value_type" xorm:"index ENUM('中位值','当前值','统计值','平均值','最大值','最小值')"`
	Describe  string `json:"describe" xorm:"TEXT"`
	// ContinuedCount  int    `json:"continued_count" xorm:"default 1 index INT(11)"`
	ContinuedTime   int    `json:"continued_time" xorm:"default 60 index INT(11)"`
	CycleTime       string `json:"cycle_time" xorm:"default 300 index INT(11)"`
	EventTemplateId int    `json:"event_template_id" xorm:"default 0 index INT(11)"`
	Enable          int32  `json:"enable" xorm:"int 'enable'"`
	IsDelete        int32  `json:"is_delete" xorm:"int 'is_delete'"`
}

func (this *EventSeting) TableName() string {
	return "zn_event_seting"
}

// 联合查询
type EventSetingTemplate struct {
	EventSeting     `xorm:"extends"`
	TemplateName    string `json:"template_name" xorm:"index VARCHAR(30)"`
	TemplateContent string `json:"template_content" xorm:"TEXT"`
}

func (this *EventSetingTemplate) TableName() string {
	return "zn_event_seting"
}

// 无条件查询全部数据
func GetEventSetingAll(field string) ([]EventSetingTemplate, error) {
	eventSeting := make([]EventSetingTemplate, 0)
	var err error
	where := " (es.enable = 1) and (es.is_delete = 0) "
	if field == "" {
		err = dbEngine().Sql("SELECT es.*, et.name as template_name, et.content as template_content FROM zn_event_seting as es LEFT JOIN zn_event_template as et ON es.event_template_id = et.id WHERE" + where).Find(&eventSeting)
	} else {
		err = dbEngine().Sql("SELECT es.*, et.name as template_name, et.content as template_content FROM zn_event_seting as es LEFT JOIN zn_event_template as et ON es.event_template_id = et.id WHERE es.field = '" + field + "' and " + where).Find(&eventSeting)
	}
	if err != nil {
		internal.LogFile.E("查询告警设置错误:" + err.Error())
		return eventSeting, err
	}
	return eventSeting, nil
}

// 无条件查询全部数据
func GetRpcEventSetingAll(field string) ([]EventSeting, error) {
	eventSeting := make([]EventSeting, 0)
	var err error
	where := " (es.enable = 1) and (es.is_delete = 0) "
	if field == "" {
		err = dbEngine().Sql("SELECT es.* FROM zn_event_seting as es WHERE" + where).Find(&eventSeting)
	} else {
		err = dbEngine().Sql("SELECT es.* FROM zn_event_seting as es WHERE es.field = '" + field + "' and " + where).Find(&eventSeting)
	}
	if err != nil {
		internal.LogFile.E("rpc查询告警设置错误:" + err.Error())
		return eventSeting, err
	}
	return eventSeting, nil
}

// 根据关键词获取总数
func GetKeywordEventSetingCount(keyword string) (c int64, err error) {
	eventSeting := new(EventSeting)
	if keyword != "" {
		c, err = dbEngine().Where("(name like '%" + keyword + "%' or field like '%" + keyword + "%' or app_name like '%" + keyword + "%') and (is_delete = 0) ").Count(eventSeting)
	} else {
		c, err = dbEngine().Where("(is_delete = 0) ").Count(eventSeting)
	}
	if err != nil {
		internal.LogFile.W("列表总数错误:" + err.Error())
	}
	return c, err
}

// 通过关键词查询列表
func GetKeywordEventSetingList(page, pageCount int, keyword string) ([]*EventSetingTemplate, error) {
	pageStart := (page - 1) * pageCount
	list := make([]*EventSetingTemplate, 0)
	var err error
	sql := " LIMIT " + strconv.Itoa(pageStart) + "," + strconv.Itoa(pageCount)
	if keyword == "" {
		where := " WHERE  (es.is_delete = 0) "
		sql := "SELECT es.*, et.name as template_name, et.content as template_content FROM zn_event_seting as es LEFT JOIN zn_event_template as et ON es.event_template_id = et.id " + where + sql
		err = dbEngine().Sql(sql).Find(&list)
	} else {
		where := " WHERE (es.app_name like '%" + keyword + "%' or es.field like '%" + keyword + "%' or es.name like '%" + keyword + "%') and (es.is_delete = 0) "
		sql := "SELECT es.*, et.name as template_name, et.content as template_content FROM zn_event_seting as es LEFT JOIN zn_event_template as et ON es.event_template_id = et.id " + where + sql
		err = dbEngine().Sql(sql).Find(&list)
	}
	if err != nil {
		internal.LogFile.W("列表错误:" + err.Error())
	}
	return list, err
}

// 修改状态
func UpdateEventSetingEnable(id, enable int32) (int64, error) {
	eventSeting := new(EventSeting)
	eventSeting.Enable = enable
	affected, err := dbEngine().Cols("enable").Where("id = " + strconv.Itoa(int(id))).Update(eventSeting)
	if err != nil {
		internal.LogFile.E("修改告警状态错误:" + err.Error())
	}
	return affected, err
}

// 添加数据
func AddOneEventSeting(eventSeting *EventSeting) (int64, error) {
	affected, err := dbEngine().Insert(eventSeting)
	if err != nil {
		internal.LogFile.E("添加告警配置错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsEventSeting(ids string) (int64, error) {
	eventSeting := new(EventSeting)
	eventSeting.IsDelete = 1
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + ids + ")").Update(eventSeting)
	if err != nil {
		internal.LogFile.E("删除告警设置配置错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneSetingsEventInfo(id int) (*EventSeting, error) {
	eventSeting := new(EventSeting)
	has, err := dbEngine().Where("id = " + strconv.Itoa(id)).Get(eventSeting)
	if !has {
		internal.LogFile.W("根据id查询setings_event错误:"+err.Error(), has)
		return eventSeting, errors.New("No query to data")
	}
	return eventSeting, nil
}

// 根据id修改信息
func UpdateIdSetingsEventInfo(id int, eventSeting *EventSeting) (int64, error) {
	affected, err := dbEngine().Cols("name,app_name,field,value_type,continued_time,cycle_time,enable,describe").Where("id = " + strconv.Itoa(id)).Update(eventSeting) // continued_count
	if err != nil {
		internal.LogFile.E("根据id修改setings_event信息错误:" + err.Error())
	}
	return affected, err
}

// 根据id修改模板
func UpdateIdSetingsEventTemplateId(sid, tid int) (int64, error) {
	eventSeting := new(EventSeting)
	eventSeting.EventTemplateId = tid
	affected, err := dbEngine().Cols("event_template_id").Where("id = " + strconv.Itoa(sid)).Update(eventSeting)
	if err != nil {
		internal.LogFile.E("根据id修改setings_event模板id错误:" + err.Error())
	}
	return affected, err
}
