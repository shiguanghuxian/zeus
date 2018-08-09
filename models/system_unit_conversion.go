package models

import (
	"errors"
	"fmt"

	"53it.net/zues/internal"
)

type SystemUnitConversion struct {
	Id           int     `json:"id" xorm:"not null pk autoincr INT(10)"`
	OriginalUnit string  `json:"original_unit" xorm:"VARCHAR(20)"`
	AfterUnit    string  `json:"after_unit" xorm:"VARCHAR(20)"`
	Multiple     float64 `json:"multiple" xorm:"double"`
	LuaCode      string  `json:"lua_code" xorm:"TEXT"`
	Type         int     `json:"type" xorm:"tinyint"`
}

func (this *SystemUnitConversion) TableName() string {
	return "zn_system_unit_conversion"
}

// 获取总数
func GetSystemUnitConversionCount(keyword string) (c int64, err error) {
	systemUnitConversion := new(SystemUnitConversion)
	if keyword == "" {
		c, err = dbEngine().Count(systemUnitConversion)
	} else {
		c, err = dbEngine().Where(fmt.Sprintf(" (original_unit like '%%%s%%') OR (after_unit like '%%%s%%') ", keyword, keyword)).Count(systemUnitConversion)
	}
	if err != nil {
		internal.LogFile.W("列表总数错误:" + err.Error())
	}
	return c, err
}

// 查询列表
func GetSystemUnitConversionList(page, pageCount int, keyword string) ([]*SystemUnitConversion, error) {
	pageStart := (page - 1) * pageCount
	list := make([]*SystemUnitConversion, 0)
	var err error
	if keyword == "" {
		err = dbEngine().Asc("id").Limit(pageCount, pageStart).Find(&list)
	} else {
		err = dbEngine().Where(fmt.Sprintf(" (original_unit like '%%%s%%') OR (after_unit like '%%%s%%') ", keyword, keyword)).Asc("id").Limit(pageCount, pageStart).Find(&list)
	}
	if err != nil {
		internal.LogFile.W("列表错误:" + err.Error())
	}
	return list, err
}

// 添加一条
func AddOneSystemUnitConversion(systemUnitConversion *SystemUnitConversion) (int64, error) {
	affected, err := dbEngine().Insert(systemUnitConversion)
	if err != nil {
		internal.LogFile.E("添加 SystemUnitConversion 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 编辑一条
func EditOneSystemUnitConversion(id int, systemUnitConversion *SystemUnitConversion) (int64, error) {
	affected, err := dbEngine().Cols("original_unit, after_unit, multiple, lua_code, type").Where(fmt.Sprintf("id = %d", id)).Update(systemUnitConversion)
	if err != nil {
		internal.LogFile.E("编辑 SystemUnitConversion 错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id删除数据
func DelOneSystemUnitConversion(id int) (int64, error) {
	systemUnitConversion := new(SystemUnitConversion)
	affected, err := dbEngine().Where(fmt.Sprintf("(id = %d)", id)).Delete(systemUnitConversion)
	if err != nil {
		internal.LogFile.E("删除 SystemUnitConversion 错误:"+err.Error(), id)
	}
	return affected, err
}

// 验证单位转换是否存在
func ChkSystemUnitConversion(original, after string) bool {
	if original == "" || after == "" {
		return false
	}
	systemUnitConversion := new(SystemUnitConversion)
	where := fmt.Sprintf(" ((original_unit = '%s') AND (after_unit = '%s')) OR ((original_unit = '%s') AND (after_unit = '%s')) ", original, after, after, original)
	c, err := dbEngine().Where(where).Count(systemUnitConversion)
	if err != nil || c > 0 {
		return true
	}
	return false
}

// 根据两个单位，获取单位转换规则
func GetDoubleUnitSystemUnitConversion(original, after string) (*SystemUnitConversion, error) {
	if original == "" || after == "" {
		return nil, errors.New("两个单位都不能为空")
	}
	systemUnitConversion := new(SystemUnitConversion)
	where := fmt.Sprintf(" ((original_unit = '%s') AND (after_unit = '%s')) OR ((original_unit = '%s') AND (after_unit = '%s')) ", original, after, after, original)
	has, err := dbEngine().Where(where).Get(systemUnitConversion)
	if err != nil {
		internal.LogFile.W("根据两个单位，获取单位转换规则错误:" + err.Error())
		return systemUnitConversion, err
	}
	if !has {
		return systemUnitConversion, errors.New("No query to data")
	}
	return systemUnitConversion, nil
}
