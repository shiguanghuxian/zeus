package apis

import (
	"errors"
	"net/http"
	"strings"

	"53it.net/zues/models"
)

type EventRule struct {
	Apis
}

// 规则列表
func (this *EventRule) RuleList(r *http.Request, args *map[string]interface{}, response *Response) error {
	esid, err := this.ToInt((*args)["esid"], 0)
	if err != nil || esid == 0 {
		return errors.New("参数错误")
	}
	list, err := models.GetSetingIdEventRuleLevelList(esid)
	if err != nil {
		return errors.New("查询列表出现错误")
	}
	*response = list
	return nil
}

// 调整排序
func (this *EventRule) EventRuleChageSort(r *http.Request, args *map[string]interface{}, response *Response) error {
	rsid, err1 := this.ToInt((*args)["rsid"], 0)
	sort, err2 := this.ToInt((*args)["sort"], 0)
	if err1 != nil || err2 != nil {
		return errors.New("参数错误")
	}
	err := models.UpEventRuleChageSort(rsid, sort)
	if err != nil {
		return errors.New("排序修改错误")
	}
	return nil
}

// 添加告警规则
func (this *EventRule) AddOneEventRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	eventSetingId, err := this.ToInt((*args)["event_seting_id"], 0)
	value := this.ToString((*args)["value"])
	if eventSetingId == 0 || err != nil || value == "" {
		return errors.New("参数错误")
	}
	eventLevelId, _ := this.ToInt((*args)["event_level_id"], 0)
	expression := this.ToString((*args)["expression"])
	if expression == "inexistence" {
		if value != "true" && value != "false" {
			return errors.New("比较值只能是true或false")
		}
	}
	if expression == "diff" {
		values := strings.Split(value, "|")
		values[0] = strings.TrimSpace(values[0])
		if len(values) != 3 {
			if values[0] != "true" && values[0] != "false" {
				return errors.New("比较值错误，请查看格式要求")
			}
		} else {
			exps := map[string]string{
				"=":  "=",
				"!=": "!=",
				">":  ">",
				">=": ">=",
				"<":  "<",
				"<=": "<=",
			}
			if exps[values[2]] == "" {
				return errors.New("比较值中的expression格式错误")
			}
		}
	}
	// models对象
	eventRule := new(models.EventRule)
	eventRule.EventLevelId = eventLevelId
	eventRule.EventSetingId = eventSetingId
	eventRule.Value = value
	eventRule.Expression = expression
	eventRule.Sort, _ = this.ToInt((*args)["sort"])
	eventRule.Unit = this.ToString((*args)["unit"])
	// 执行插入
	_, err = models.AddOneEventRule(eventRule)
	if err != nil {
		return errors.New("告警规则添加错误")
	}
	return nil
}

// 删除
func (this *EventRule) DelEventRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsEventRule(ids)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}

// 获取信息
func (this *EventRule) InfoEventRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, err := this.ToInt((*args)["id"], 0)
	if id == 0 || err != nil {
		return errors.New("参数错误")
	}
	info, err := models.GetOneEventRuleInfo(id)
	if err != nil {
		return errors.New("获取信息错误")
	}
	*response = info
	return nil
}

// 保存编辑信息
func (this *EventRule) UpEventRuleInfo(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, err := this.ToInt((*args)["id"], 0)
	if err != nil {
		return errors.New("参数id错误")
	}
	value := this.ToString((*args)["value"])
	if value == "" {
		return errors.New("参数value错误")
	}
	eventLevelId, _ := this.ToInt((*args)["event_level_id"], 0)
	expression := this.ToString((*args)["expression"])
	if expression == "inexistence" {
		if value != "true" && value != "false" {
			return errors.New("比较值只能是true或false")
		}
	}
	if expression == "diff" {
		values := strings.Split(value, "|")
		values[0] = strings.TrimSpace(values[0])
		if len(values) != 3 {
			if values[0] != "true" && values[0] != "false" {
				return errors.New("比较值错误，请查看格式要求")
			}
		} else {
			exps := map[string]string{
				"=":  "=",
				"!=": "!=",
				">":  ">",
				">=": ">=",
				"<":  "<",
				"<=": "<=",
			}
			if exps[values[2]] == "" {
				return errors.New("比较值中的expression格式错误")
			}
		}
	}
	// models对象
	eventRule := new(models.EventRule)
	eventRule.EventLevelId = eventLevelId
	eventRule.Value = value
	eventRule.Expression = expression
	eventRule.Sort, _ = this.ToInt((*args)["sort"])
	eventRule.Unit = this.ToString((*args)["unit"])
	// 执行插入
	_, err = models.UpdateIdEventRuleInfo(id, eventRule)
	if err != nil {
		return errors.New("告警规则编辑错误")
	}
	return nil
}
