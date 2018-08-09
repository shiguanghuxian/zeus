package models

import (
	"errors"
	"strconv"

	"fmt"

	"53it.net/zues/internal"
)

type AppnameFieldType struct {
	Id       int32  `json:"id" xorm:"not null pk autoincr INT(11)"`
	AppName  string `json:"app_name" xorm:"default 'zn_raw_data' VARCHAR(60)"`
	Field    string `json:"field" xorm:"VARCHAR(60)"`
	Type     string `json:"type" xorm:"default 'string' VARCHAR(30)"`
	Unit     string `json:"unit" xorm:"VARCHAR(30)"`
	Index    int32  `json:"index" xorm:"INT(11)"`
	IsDelete int32  `json:"is_delete" xorm:"int 'is_delete'"`
}

// 实际表名
func (this *AppnameFieldType) TableName() string {
	return "zn_appname_field_type"
}

// 获取全部列表－未删除
func GetNotDeleteAppnameFieldTypeList() ([]*AppnameFieldType, error) {
	list := make([]*AppnameFieldType, 0)
	err := dbEngine().Where("(is_delete = 0)").Asc("id").Find(&list)
	if err != nil {
		internal.LogFile.E("查询appname列表失败:" + err.Error())
		return nil, err
	}
	return list, nil
}

// 根据appname查询数据
func GetNameNotDeleteAppnameFieldTypeList(name string) ([]*AppnameFieldType, error) {
	list := make([]*AppnameFieldType, 0)
	err := dbEngine().Where("(is_delete = 0) and (app_name = '" + name + "')").Asc("id").Find(&list)
	if err != nil {
		internal.LogFile.E("根据name查询appname列表失败:" + err.Error())
		return nil, err
	}
	return list, nil
}

// 根据关键词获取总数
func GetKeywordAppNameCount(keyword string) (c int64, err error) {
	appnameFieldType := new(AppnameFieldType)
	if keyword != "" {
		c, err = dbEngine().Where("(app_name like '%" + keyword + "%' or field like '%" + keyword + "%' or type like '%" + keyword + "%') and (is_delete = 0)").Count(appnameFieldType)
	} else {
		c, err = dbEngine().Where("(is_delete = 0)").Count(appnameFieldType)
	}
	if err != nil {
		internal.LogFile.W("获取appname列表总数错误:" + err.Error())
	}
	return c, err
}

// 通过关键词查询列表
func GetKeywordAppNameList(page, pageCount int, keyword string) ([]*AppnameFieldType, error) {
	pageStart := (page - 1) * pageCount
	list := make([]*AppnameFieldType, 0)
	var err error
	if keyword == "" {
		err = dbEngine().Where("(is_delete = 0)").Asc("id").Limit(pageCount, pageStart).Find(&list)
	} else {
		err = dbEngine().Where("(app_name like '%"+keyword+"%' or field like '%"+keyword+"%' or type like '%"+keyword+"%') and (is_delete = 0)").Asc("id").Limit(pageCount, pageStart).Find(&list)
	}
	if err != nil {
		internal.LogFile.W("获取appname列表错误:" + err.Error())
	}
	return list, err
}

// 添加一条
func AddOneAppName(appName *AppnameFieldType) (int64, error) {
	affected, err := dbEngine().Insert(appName)
	if err != nil {
		internal.LogFile.E("添加appname错误:"+err.Error(), affected)
	}
	return affected, err
}

// 根据id列表删除数据
func DelIdsAppName(ids string) (int64, error) {
	appNameFieldType := new(AppnameFieldType)
	appNameFieldType.IsDelete = 1
	affected, err := dbEngine().Cols("is_delete").Where("id in (" + ids + ")").Update(appNameFieldType)
	if err != nil {
		internal.LogFile.E("删除appname错误:"+err.Error(), ids)
	}
	return affected, err
}

// 根据id查询单条数据
func GetOneAppNameInfo(id int32) (*AppnameFieldType, error) {
	appNameFieldType := new(AppnameFieldType)
	has, err := dbEngine().Where("id = " + strconv.Itoa(int(id))).Get(appNameFieldType)
	if !has {
		internal.LogFile.W("根据id查询appname错误:"+err.Error(), has)
		return appNameFieldType, errors.New("No query to data")
	}
	return appNameFieldType, nil
}

// 根据id修改信息
func UpdateIdAppNameInfo(id int32, appNameFieldType *AppnameFieldType) (int64, error) {
	affected, err := dbEngine().Cols("app_name,field,type,unit,index").Where("id = " + strconv.Itoa(int(id))).Update(appNameFieldType)
	if err != nil {
		internal.LogFile.E("根据id修改appname信息错误:" + err.Error())
	}
	return affected, err
}

// 根据appname和field查询信息
func GetAppnameFieldAppNameInfo(appname, field string) (*AppnameFieldType, error) {
	appNameFieldType := new(AppnameFieldType)
	has, err := dbEngine().Where(fmt.Sprintf("(app_name = '%s') AND (field = '%s')", appname, field)).Get(appNameFieldType)
	if err != nil {
		internal.LogFile.W("根据appname和field查询信息错误:" + err.Error())
		return appNameFieldType, err
	}
	if !has {
		return appNameFieldType, errors.New("No query to data")
	}
	return appNameFieldType, nil
}
