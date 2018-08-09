package controllers

import (
	"53it.net/zues/elasticsearch"
	"53it.net/zues/influxdb"
	"53it.net/zues/internal"
	"53it.net/zues/mongo"
)

type SearchController struct {
	BaseController
}

// 查询api
func (this *SearchController) AjaxZqlQueryV1() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	zqlStr := this.GetString("zql", "")
	if zqlStr == "" {
		ajaxData.Msg = "查询语句不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 获取数据源类型
	dataType := this.GetString("data_source")
	// 区分数据源
	var list interface{}
	var err error
	if dataType == "" {
		dataType = "influxdb"
	}
	if dataType == "influxdb" {
		// infludedb 附加参数
		group := this.GetString("group")
		hostname := this.GetString("hostname")
		ip := this.GetString("ip")
		list, err = influxdb.ZqlQueryCmd(zqlStr, group, hostname, ip)
	} else if dataType == "mongodb" {
		list, err = mongo.GetZqlList(zqlStr)
	} else if dataType == "elastic" {
		list, err = elasticsearch.GetZqlList(zqlStr)
	}

	if err != nil {
		ajaxData.Msg = "查询错误:" + err.Error()
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "ok", Data: list}
	this.AjaxReturn(ajaxData)
}

// 当前可用数据源
func (this *SearchController) AjaxDataSource() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	list := make([]map[string]string, 0)
	// 读配置组织数据
	ismgo, err := internal.CFG.Bool("mongodb", "enable")
	if err != nil {
		ajaxData.Msg = "读区配置文件错误"
		this.AjaxReturn(ajaxData)
	}
	if ismgo == true {
		list = append(list, map[string]string{
			"type": "mongodb",
			"name": "mongodb",
		})
	}
	isinf, _ := internal.CFG.Bool("influxdb", "enable")
	if isinf == true {
		list = append(list, map[string]string{
			"type": "influxdb",
			"name": "influxdb",
		})
	}
	isela, _ := internal.CFG.Bool("elasticsearch", "enable")
	if isela == true {
		list = append(list, map[string]string{
			"type": "elastic",
			"name": "elastic",
		})
	}
	ajaxData = &AjaxData{State: 0, Msg: "ok", Data: list}
	this.AjaxReturn(ajaxData)
}
