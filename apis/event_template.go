package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
)

type SetingsTemplate struct {
	Apis
}

// 获取模版列表-全部
func (this *SetingsTemplate) GetAllSetingsTemplate(r *http.Request, args *map[string]interface{}, response *Response) error {
	list, err := models.GetAllEventTemplateList()
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// 添加
func (this *SetingsTemplate) AddSetingsTemplate(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt((*args)["id"], 0)
	sid, err := this.ToInt((*args)["seting_event_id"], 0)
	if err != nil || sid == 0 {
		return errors.New("参数[seting_event_id]不能为空")
	}
	if id == 0 {
		name := this.ToString((*args)["name"])
		content := this.ToString((*args)["content"])
		if name == "" || content == "" {
			return errors.New("模板名和内容不能为空")
		}
		// 话题models对象
		eventTemplate := new(models.EventTemplate)
		eventTemplate.Name = name
		eventTemplate.Content = content
		// 执行插入
		_, err = models.AddOneEventTemplate(eventTemplate)
		if err != nil {
			return errors.New("告警模板添加错误")
		}
		// 修改告警设置的模板id
		_, err = models.UpdateIdSetingsEventTemplateId(sid, int(eventTemplate.Id))
		if err != nil {
			return errors.New("模板已添加，保存到告警设置错误，请联系开发者")
		}
	} else {
		// 修改告警设置的模板id
		_, err := models.UpdateIdSetingsEventTemplateId(sid, id)
		if err != nil {
			return errors.New("保存到告警设置错误，请联系开发者")
		}
	}
	return nil
}

// 获取模板信息
func (this *SetingsTemplate) InfoTemplate(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	info, err := models.GetOneEventTemplateInfo(id)
	if err != nil {
		return errors.New("获取信息错误")
	}
	*response = info
	return nil
}

// 编辑
func (this *SetingsTemplate) UpSetingsTemplate(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, err := this.ToInt((*args)["id"], 0)
	if err != nil || id == 0 {
		return errors.New("参数id错误")
	}
	setingEventId, _ := this.ToInt((*args)["seting_event_id"], 0)
	if setingEventId == 0 {
		return errors.New("参数seting_event_id错误")
	}
	name := this.ToString((*args)["name"])
	content := this.ToString((*args)["content"])
	if name == "" || content == "" {
		return errors.New("模板名和内容不能为空")
	}
	// 话题models对象
	eventTemplate := new(models.EventTemplate)
	eventTemplate.Name = name
	eventTemplate.Content = content
	// 执行修改
	_, err = models.UpdateIdEventTemplateInfo(id, eventTemplate)
	if err != nil {
		return errors.New("告警模板编辑错误")
	}
	// 修改告警设置的模板id
	_, err = models.UpdateIdSetingsEventTemplateId(setingEventId, id)
	if err != nil {
		return errors.New("模板已修改，保存到告警设置错误，请联系开发者")
	}
	return nil
}
