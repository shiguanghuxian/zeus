package controllers

import (
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/proto"
	"github.com/astaxie/beego"
)

// 话题设置
type TopicsController struct {
	BaseController
}

// 获取话题列表
func (this *TopicsController) AjaxTopicsList() {
	ajaxData := &AjaxData{State: 1, Msg: "数据获取失败"}

	var totalRows int64 // 总行数
	var err error

	keyword := this.GetString("keyword") // 关键次参数
	if keyword != "" {
		totalRows, err = models.GetKeywordTopicsCount(keyword)
	} else {
		totalRows, err = models.GetKeywordTopicsCount("")
	}
	if err != nil {
		ajaxData.Msg = "服务端错误 count"
		this.Data["json"] = ajaxData
		this.ServeJSON()
		this.StopRun()
	}
	// 每页行数
	listRows, err := beego.AppConfig.Int("AdminListPagesCount")
	if err != nil {
		listRows = 10
	}
	// 当前页码
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	// 查询列表
	list, err := models.GetKeywordTopicsList(page, listRows, keyword)
	if err != nil {
		ajaxData.Msg = "服务端错误 list"
		this.Data["json"] = ajaxData
		this.ServeJSON()
		this.StopRun()
	}
	// 页面数据
	data := make(map[string]interface{})
	data["page"] = page
	data["total_rows"] = totalRows
	data["list_rows"] = listRows
	data["list"] = list
	ajaxData = &AjaxData{State: 0, Msg: "数据获取成功", Data: data}

	// 返回数据
	this.Data["json"] = ajaxData
	this.ServeJSON()
}

// 修改话题状态
func (this *TopicsController) AjaxChangeEnable() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, err := this.GetInt32("id", 0)
	enable, err := this.GetInt32("enable", 0)
	if err != nil {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	if enable == 0 {
		enable = 1
	} else {
		enable = 0
	}
	// 调用修改
	_, err = models.UpdateTopicsEnable(id, enable)
	if err != nil {
		ajaxData.Msg = "修改话题状态错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 新增话题
func (this *TopicsController) AjaxAddTopics() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	topics := this.GetString("topics")
	channel := this.GetString("channel")
	if topics == "" || channel == "" {
		ajaxData.Msg = "话题和通道不能为空"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	topicsConfig := new(models.TopicsConfig)
	topicsConfig.Topics = topics
	topicsConfig.Channel = channel
	// 其它参数
	topicsConfig.ChannelCount, _ = this.GetInt32("channel_count", 0)
	topicsConfig.DataType = this.GetString("data_type")
	topicsConfig.Enable, _ = this.GetInt32("enable", 0)
	// 执行插入
	_, err := models.AddOneTopicsConfig(topicsConfig)
	if err != nil {
		ajaxData.Msg = "话题配置添加错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 删除
func (this *TopicsController) AjaxDelTopics() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	ids := this.GetString("ids")
	if ids == "" {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsTopicsConfig(ids)
	if err != nil {
		ajaxData.Msg = "话题配置删除错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 获取单挑消息
func (this *TopicsController) AjaxInfoTopics() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	id, _ := this.GetInt32("id", 0)
	if id == 0 {
		ajaxData.Msg = "参数错误"
		this.AjaxReturn(ajaxData)
	}
	info, err := models.GetOneTopicsInfo(id)
	if err != nil {
		ajaxData.Msg = "获取会话信息错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功", Data: info}
	this.AjaxReturn(ajaxData)
}

// 保存信息
func (this *TopicsController) AjaxUpTopics() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 接收参数
	id, _ := this.GetInt32("id", 0)
	topics := this.GetString("topics")
	channel := this.GetString("channel")
	if topics == "" || channel == "" || id == 0 {
		ajaxData.Msg = "话题和通道不能为空,id不能为0"
		this.AjaxReturn(ajaxData)
	}
	// 话题models对象
	topicsConfig := new(models.TopicsConfig)
	topicsConfig.Topics = topics
	topicsConfig.Channel = channel
	// 其它参数
	topicsConfig.ChannelCount, _ = this.GetInt32("channel_count", 0)
	topicsConfig.DataType = this.GetString("data_type")
	topicsConfig.Enable, _ = this.GetInt32("enable", 0)
	// 通道数不能为0
	if topicsConfig.ChannelCount == 0 {
		ajaxData.Msg = "通道数不能为0"
		this.AjaxReturn(ajaxData)
	}
	// 执行修改
	_, err := models.UpdateIdTopicsInfo(id, topicsConfig)
	if err != nil {
		ajaxData.Msg = "话题配置修改错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "成功"}
	this.AjaxReturn(ajaxData)
}

// 重启serverd服务
func (this *TopicsController) AjaxRestartServerd() {
	ajaxData := &AjaxData{State: 1, Msg: "服务器错误"}
	// 读取调度器配置信息
	address := beego.AppConfig.DefaultString("DispatchdAddress", "127.0.0.1")
	port := beego.AppConfig.DefaultString("DispatchdPort", "3200")

	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure(), grpc.WithTimeout(30*time.Second))
	if err != nil {
		internal.LogFile.E("调度器服务未开启" + err.Error())
		ajaxData.Msg = "调度器服务未开启"
		this.AjaxReturn(ajaxData)
	}
	defer conn.Close()
	c := proto.NewReportServerdServiceClient(conn)
	// 发起请求
	r, err := c.SayRestartServerd(context.Background(), &proto.RestartRequest{})
	if err != nil {
		internal.LogFile.E("同步话题配置，发起请求错误：" + err.Error())
		ajaxData.Msg = "同步话题配置，发起请求错误"
		this.AjaxReturn(ajaxData)
	}
	if r.Code != "0" {
		internal.LogFile.E("同步话题配置,错误码：" + r.Code)
		ajaxData.Msg = "同步话题配置错误"
		this.AjaxReturn(ajaxData)
	}
	ajaxData = &AjaxData{State: 0, Msg: "同步话题配置成功:" + r.Message}
	this.AjaxReturn(ajaxData)
}

// nsq所有话题列表（调用nsq接口）
func (this *TopicsController) NsqTopics() {

}
