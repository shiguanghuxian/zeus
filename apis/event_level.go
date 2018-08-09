package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
)

type EventLevel struct {
	Apis
}

// 级别列表
func (this *EventLevel) EventLevelList(r *http.Request, args *map[string]interface{}, response *Response) error {
	list, err := models.GetAllEventLevelList()
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// 删除
func (this *EventLevel) DelEventLevel(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.DelIdEventLevel(id)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}

// 添加
func (this *EventLevel) AddEventLevel(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	level, err := this.ToInt((*args)["level"], 0)
	name := this.ToString((*args)["name"])
	if err != nil || name == "" {
		return errors.New("参数错误")
	}
	// models对象
	eventLevel := new(models.EventLevel)
	eventLevel.Name = name
	eventLevel.Level = level
	// 执行插入
	_, err = models.AddOneEventLevel(eventLevel)
	if err != nil {
		return errors.New("告警级别添加错误")
	}
	return nil
}

// 获取信息
func (this *EventLevel) InfoEventLevel(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, err := this.ToInt((*args)["id"], 0)
	if id == 0 || err != nil {
		return errors.New("参数错误")
	}
	info, err := models.GetOneEventLevelInfo(id)
	if err != nil {
		return errors.New("获取信息错误")
	}
	*response = info
	return nil
}

// 编辑
func (this *EventLevel) EditEventLevel(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, err := this.ToInt((*args)["id"], 0)
	if err != nil || id == 0 {
		return errors.New("参数id错误")
	}
	level, err := this.ToInt((*args)["level"], 0)
	name := this.ToString((*args)["name"])
	if err != nil || name == "" {
		return errors.New("参数name错误")
	}
	// models对象
	eventLevel := new(models.EventLevel)
	eventLevel.Name = name
	eventLevel.Level = level
	// 执行插入
	_, err = models.UpdateIdEventLevelInfo(id, eventLevel)
	if err != nil {
		return errors.New("告警级别编辑错误")
	}
	return nil
}
