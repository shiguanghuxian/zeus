package apis

import (
	"errors"
	"net/http"

	"53it.net/zues/models"
)

// 话题设置
type TopicsRule struct {
	Apis
}

// 获取规则列表
func (this *TopicsRule) IdTopicsRules(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, err := this.ToInt((*args)["id"], 0)
	if err != nil || id < 1 {
		return errors.New("参数错误")
	}
	list, err := models.GetTCIdTopicsConfigRuleList(id)
	if err != nil {
		return errors.New("查询解析规则列表失败")
	}

	*response = list
	return nil
}

// 修改状态
func (this *TopicsRule) ChangeEnable(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, err := this.ToInt((*args)["id"], 0)
	enable, err := this.ToInt((*args)["enable"], 0)
	if err != nil || id == 0 {
		return errors.New("参数错误")
	}
	if enable == 0 {
		enable = 1
	} else {
		enable = 0
	}
	// 调用修改
	_, err = models.UpdateTopicsRuleEnable(id, enable)
	if err != nil {
		return errors.New("修改话题状态错误")
	}
	return nil
}

// 添加解析规则
func (this *TopicsRule) AddTopicsRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	topicsConfigId, _ := this.ToInt((*args)["topics_config_id"], 0)
	if topicsConfigId == 0 {
		return errors.New("话题配置id错误")
	}
	// 其它参数
	topicsConfigRule := new(models.TopicsConfigRule)
	topicsConfigRule.TopicsConfigId = topicsConfigId
	topicsConfigRule.AppName = this.ToString((*args)["app_name"])
	topicsConfigRule.Tag = this.ToString((*args)["tag"])
	topicsConfigRule.Mapped = this.ToString((*args)["mapped"])
	topicsConfigRule.TextUnType = this.ToString((*args)["text_un_type"])
	topicsConfigRule.TextUnRule = this.ToString((*args)["text_un_rule"])
	topicsConfigRule.DateFormat = this.ToString((*args)["date_format"])
	topicsConfigRule.Sort, _ = this.ToInt((*args)["sort"], 1)
	topicsConfigRule.Enable, _ = this.ToInt((*args)["enable"], 1)
	// 调用添加
	_, err := models.AddOneTopicsConfigRule(topicsConfigRule)
	if err != nil {
		return errors.New("话题解析规则添加错误")
	}

	return nil
}

// 删除
func (this *TopicsRule) DelTopicsRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	id := this.ToString((*args)["id"])
	if id == "" {
		return errors.New("参数错误")
	}
	_, err := models.DelIdsTopicsConfigRule(id)
	if err != nil {
		return errors.New("话题解析规则删除错误")
	}
	return nil
}

// 编辑
func (this *TopicsRule) EditTopicsRuleInfo(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt32((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	info, err := models.GetOneTopicsRuleInfo(id)
	if err != nil {
		return errors.New("获取话题解析规则信息错误")
	}
	*response = info
	return nil
}

// 保存修改
func (this *TopicsRule) UpTopicsRule(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("话题解析规则id错误")
	}
	// 其它参数
	topicsConfigRule := new(models.TopicsConfigRule)

	topicsConfigRule.AppName = this.ToString((*args)["app_name"])
	topicsConfigRule.Tag = this.ToString((*args)["tag"])
	topicsConfigRule.Mapped = this.ToString((*args)["mapped"])
	topicsConfigRule.TextUnType = this.ToString((*args)["text_un_type"])
	topicsConfigRule.TextUnRule = this.ToString((*args)["text_un_rule"])
	topicsConfigRule.DateFormat = this.ToString((*args)["date_format"])
	topicsConfigRule.Sort, _ = this.ToInt((*args)["sort"], 1)
	topicsConfigRule.Enable, _ = this.ToInt((*args)["enable"], 1)
	// 调用添加
	_, err := models.UpdateIdTopicsRuleInfo(id, topicsConfigRule)
	if err != nil {
		return errors.New("话题解析规则修改错误")
	}

	return nil
}
