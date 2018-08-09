package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
)

type EventPush struct {
	Apis
}

// 级别列表
func (this *EventPush) EventPushList(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	esid, err := this.ToInt((*args)["esid"], 0)
	if esid == 0 && err != nil {
		return errors.New("参数错误")
	}
	list, err := models.GetAllEventPushList(esid)
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// AjaxAddEventPush 添加
func (this *EventPush) AddEventPush(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	name := this.ToString((*args)["name"])
	purl := this.ToString((*args)["url"])
	esid, _ := this.ToInt((*args)["event_seting_id"], 0)
	dataType, _ := this.ToInt((*args)["data_type"], 0)
	if esid == 0 || name == "" || purl == "" {
		return errors.New("参数错误")
	}
	// models对象
	eventPush := new(models.EventPush)
	eventPush.Name = name
	eventPush.Url = purl
	eventPush.DataType = dataType
	eventPush.EventSetingId = esid
	// 执行插入
	_, err := models.AddOneEventPush(eventPush)
	if err != nil {
		return errors.New("告警推送添加错误")
	}
	return nil
}

// 删除
func (this *EventPush) DelEventPush(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.DelIdEventPush(id)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}

// AjaxUpEventPush 编辑
func (this *EventPush) UpEventPush(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	pId, _ := this.ToInt((*args)["id"], 0)
	name := this.ToString((*args)["name"])
	purl := this.ToString((*args)["url"])
	dataType, _ := this.ToInt((*args)["data_type"], 0)
	if pId == 0 || name == "" || purl == "" {
		return errors.New("参数错误")
	}
	// models对象
	eventPush := new(models.EventPush)
	eventPush.Name = name
	eventPush.Url = purl
	eventPush.DataType = dataType
	// 执行插入
	_, err := models.UpdateIdEventPushInfo(pId, eventPush)
	if err != nil {
		return errors.New("告警推送编辑错误")
	}
	return nil
}
