package apis

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"53it.net/zues/internal"
	"53it.net/zues/models"
)

type AppName struct {
	Apis
}

// appname 列表
func (this *AppName) AppNameList(r *http.Request, args *map[string]interface{}, response *Response) error {
	var totalRows int64 // 总行数
	var err error

	keyword := this.ToString((*args)["keyword"]) // 关键次参数
	totalRows, err = models.GetKeywordAppNameCount(keyword)
	if err != nil {
		return errors.New("服务端错误 count")
	}
	// 每页行数
	listRows, err := internal.CFG.Int("apis", "pagecount")
	if err != nil {
		listRows = 10
	}
	// 当前页码
	page, err := this.ToInt((*args)["page"], 1)
	if err != nil {
		page = 1
	}
	// 查询列表
	list, err := models.GetKeywordAppNameList(page, listRows, keyword)
	if err != nil {
		return errors.New("服务端错误 list")
	}
	// 页面数据
	data := make(map[string]interface{})
	data["page"] = page
	data["total_rows"] = totalRows
	data["list_rows"] = listRows
	data["list"] = list
	*response = data
	return nil
}

// 添加appname字段配置
func (this *AppName) AddAppName(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	appname := this.ToString((*args)["app_name"])
	field := this.ToString((*args)["field"])
	if appname == "" || field == "" {
		return errors.New("应用名和字段不能为空")
	}
	// 话题models对象
	appNameFieldType := new(models.AppnameFieldType)
	appNameFieldType.AppName = appname
	appNameFieldType.Field = field
	// 其它参数
	appNameFieldType.Type = this.ToString((*args)["type"])
	appNameFieldType.Unit = this.ToString((*args)["unit"])
	appNameFieldType.Index, _ = this.ToInt32((*args)["index"], 0)
	// 执行插入
	_, err := models.AddOneAppName(appNameFieldType)
	if err != nil {
		return errors.New("AppName添加错误")
	}
	return nil
}

// 删除
func (this *AppName) DelAppName(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsAppName(ids)
	if err != nil {
		return errors.New("AppName删除错误")
	}
	return nil
}

// 获取单条消息
func (this *AppName) InfoAppName(r *http.Request, args *map[string]interface{}, response *Response) error {
	log.Println(*args)
	id, err := this.ToInt32((*args)["id"], 0)
	log.Println(err)
	if id == 0 {
		return errors.New("参数错误")
	}
	info, err := models.GetOneAppNameInfo(id)
	if err != nil {
		return errors.New("获取AppName信息错误")
	}
	*response = info
	return nil
}

// 编辑appname字段配置
func (this *AppName) EditAppName(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt32((*args)["id"], 0)
	appname := this.ToString((*args)["app_name"])
	field := this.ToString((*args)["field"])
	if id == 0 || appname == "" || field == "" {
		return errors.New("参数错误")
	}
	// 话题models对象
	appNameFieldType := new(models.AppnameFieldType)
	appNameFieldType.AppName = appname
	appNameFieldType.Field = field
	// 其它参数
	appNameFieldType.Type = this.ToString((*args)["type"])
	appNameFieldType.Unit = this.ToString((*args)["unit"])
	appNameFieldType.Index, _ = this.ToInt32((*args)["index"], 0)
	// 执行插入
	_, err := models.UpdateIdAppNameInfo(id, appNameFieldType)
	if err != nil {
		return errors.New("AppName编辑错误")
	}
	return nil
}
