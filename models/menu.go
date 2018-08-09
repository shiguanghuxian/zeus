package models

import "errors"

type Menu struct {
	Id         int    `json:"id" xorm:"not null pk autoincr INT(10)"`
	Name       string `json:"name" xorm:"VARCHAR(60)"`
	NameEn     string `json:"name_en" xorm:"VARCHAR(60)"`
	ParentId   int    `json:"parent_id" xorm:"default 0 index INT(11)"`
	Sort       int    `json:"sort" xorm:"default 0 index INT(11)"`
	Url        string `json:"url" xorm:"index VARCHAR(180)"`
	Icon       string `json:"icon" xorm:"VARCHAR(255)"`
	Other      string `json:"other" xorm:"VARCHAR(255)"`
	RoleTypeId int    `json:"role_type_id" xorm:"default 2 index INT(11)"`
	Show       int    `json:"show" xorm:"default 1 index INT(11)"`
}

func (this *Menu) TableName() string {
	return "zn_menu"
}

// 获取左侧菜单树
func GetMenuTree(pid int, rtid string) ([]map[string]interface{}, error) {
	tree := make([]map[string]interface{}, 0)
	var oneTree map[string]interface{}

	list := make([]Menu, 0)
	err := engine.Where("(`show` = ?) and (parent_id = ?) and (role_type_id in (?))", 1, pid, rtid).Desc("sort").Find(&list)
	if err != nil {
		return nil, errors.New("左侧菜单获取失败")
	}

	// 查询子集
	for _, v := range list {
		oneTree = map[string]interface{}{
			"id":      v.Id,
			"name":    v.Name,
			"name_en": v.NameEn,
			"url":     v.Url,
			"icon":    v.Icon,
			"other":   v.Other,
		}
		child, err := GetMenuTree(v.Id, rtid) // 子级
		if err == nil && len(child) > 0 {
			oneTree["child"] = child
		}
		tree = append(tree, oneTree) // 追加到数组
	}

	return tree, nil
}
