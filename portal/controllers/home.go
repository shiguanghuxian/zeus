package controllers

import (
	"log"

	"53it.net/zql"
	"53it.net/zues/influxdb"
	"53it.net/zues/mongo"
)

type HomeController struct {
	BaseController
}

func (this *HomeController) Index() {
	this.TplName = "home/index.html"
}

// 根据zql查询数据
func (this *HomeController) AjaxZqlQuery() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	// 接收查询语句
	zqlStr := this.GetString("zql", "")
	if zqlStr == "" {
		ajaxData.Msg = "查询语句不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 这里需要判断用的什么存储
	zqlObj, err := zql.New("", zqlStr)
	if err != nil {
		ajaxData.Msg = "创建查询对象错误"
		this.AjaxReturn(ajaxData)
	}
	zqlQuery, err := zqlObj.GetInfluxdbQuery("_default_http-server3_101.200.174.134")
	log.Println(err)
	if err != nil {
		ajaxData.Msg = "创建查询语句错误"
		this.AjaxReturn(ajaxData)
	}
	log.Println(zqlQuery)
	response, err := influxdb.QueryCmd(zqlQuery)
	if err != nil {
		ajaxData.Msg = "执行查询错误"
		this.AjaxReturn(ajaxData)
	}
	//	log.Println(response)
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: response[0].Series}
	this.AjaxReturn(ajaxData)
}

// 测试查询结果
func (this *HomeController) AjaxZqlList() {
	// zql := "SELECT * FROM \"zn_raw_data_default_http-server3_101_200_174_134\""
	zql := "SELECT count(value) FROM \"zn_raw_data_default_http-server3_101_200_174_134\"  where time > '2016-08-16'  group by time(10m)"
	response, _ := influxdb.QueryResponse(zql)
	ajaxData := &AjaxData{State: 0, Msg: "成功", Data: response}
	this.AjaxReturn(ajaxData)
}

// 测试mongodb查询
func (this *HomeController) AjaxZqlMongoList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}
	zqlStr := this.GetString("zql", "")
	if zqlStr == "" {
		ajaxData.Msg = "查询语句不能为空"
		this.AjaxReturn(ajaxData)
	}
	list, _ := mongo.GetZqlList(zqlStr)
	ajaxData = &AjaxData{State: 0, Msg: "ok", Data: list}
	this.AjaxReturn(ajaxData)
}
