package models

import (
	"fmt"
	"strconv"

	"53it.net/zues/internal"
)

type Autocomplete struct {
	Id      int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name    string `json:"name" xorm:"VARCHAR(60)"`
	Desc    string `json:"desc" xorm:"TEXT"`
	Addtime int64  `json:"addtime" xorm:"INT"`
	Count   int    `json:"count" xorm:"INT"`
}

func (this *Autocomplete) TableName() string {
	return "zn_autocomplete"
}

// 获取最近提示信息
func GetAllAutocompleteList(name string) (list []Autocomplete, err error) {
	if name == "" {
		err = dbEngine().OrderBy("count desc,addtime desc").Limit(50).Find(&list)
	} else {
		err = dbEngine().Where(fmt.Sprintf("name like '%%%s%%'", name)).OrderBy("count desc,addtime desc").Limit(50).Find(&list)
	}
	if err != nil {
		internal.LogFile.E("查询搜索提示列表:" + err.Error())
		return list, err
	}
	return list, nil
}

// 添加一条
func AddOneAutocomplete(autocomplete *Autocomplete) (int64, error) {
	affected, err := dbEngine().Insert(autocomplete)
	if err != nil {
		internal.LogFile.E("添加 autocomplete 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id修改信息
func UpdateIdAutocomplete(id int, autocomplete *Autocomplete) (int64, error) {
	affected, err := dbEngine().
		Cols("count").
		Where("id = " + strconv.Itoa(id)).
		Update(autocomplete)
	if err != nil {
		internal.LogFile.E("根据id修改 autocomplete 信息错误:" + err.Error())
	}
	return affected, err
}
