package apis

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/proto"
)

// 话题设置
type Topics struct {
	Apis
}

// 获取话题列表
func (this *Topics) TopicsList(r *http.Request, args *map[string]interface{}, response *Response) error {
	var totalRows int64 // 总行数
	var err error

	keyword := this.ToString((*args)["keyword"]) // 关键次参数
	if keyword != "" {
		totalRows, err = models.GetKeywordTopicsCount(keyword)
	} else {
		totalRows, err = models.GetKeywordTopicsCount("")
	}
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
	list, err := models.GetKeywordTopicsList(page, listRows, keyword)
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

// 修改话题状态
func (this *Topics) ChangeEnable(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, err := this.ToInt32((*args)["id"], 0)
	enable, err := this.ToInt32((*args)["enable"], 0)
	if err != nil {
		return errors.New("参数错误")
	}
	if enable == 0 {
		enable = 1
	} else {
		enable = 0
	}
	// 调用修改
	_, err = models.UpdateTopicsEnable(id, enable)
	if err != nil {
		return errors.New("修改话题状态错误")
	}
	return nil
}

// 新增话题
func (this *Topics) AddTopics(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	topics := this.ToString((*args)["topics"])
	channel := this.ToString((*args)["channel"])
	if topics == "" || channel == "" {
		return errors.New("话题和通道不能为空")
	}
	// 话题models对象
	topicsConfig := new(models.TopicsConfig)
	topicsConfig.Topics = topics
	topicsConfig.Channel = channel
	// 其它参数
	topicsConfig.ChannelCount, _ = this.ToInt32((*args)["channel_count"], 0)
	topicsConfig.DataType = this.ToString((*args)["data_type"])
	topicsConfig.Enable, _ = this.ToInt32((*args)["enable"], 0)
	// 执行插入
	_, err := models.AddOneTopicsConfig(topicsConfig)
	if err != nil {
		return errors.New("话题配置添加错误")
	}
	return nil
}

// 删除
func (this *Topics) DelTopics(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsTopicsConfig(ids)
	if err != nil {
		return errors.New("话题配置删除错误")
	}
	return nil
}

// 获取单挑消息
func (this *Topics) InfoTopics(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt32((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	info, err := models.GetOneTopicsInfo(id)
	if err != nil {
		return errors.New("获取会话信息错误")
	}
	*response = info
	return nil
}

// 保存信息
func (this *Topics) UpTopics(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt32((*args)["id"], 0)
	topics := this.ToString((*args)["topics"])
	channel := this.ToString((*args)["channel"])
	if topics == "" || channel == "" || id == 0 {
		return errors.New("话题和通道不能为空,id不能为0")
	}
	// 话题models对象
	topicsConfig := new(models.TopicsConfig)
	topicsConfig.Topics = topics
	topicsConfig.Channel = channel
	// 其它参数
	topicsConfig.ChannelCount, _ = this.ToInt32((*args)["channel_count"], 0)
	topicsConfig.DataType = this.ToString((*args)["data_type"])
	topicsConfig.Enable, _ = this.ToInt32((*args)["enable"], 0)
	// 通道数不能为0
	if topicsConfig.ChannelCount == 0 {
		return errors.New("通道数不能为0")
	}
	// 执行修改
	_, err := models.UpdateIdTopicsInfo(id, topicsConfig)
	if err != nil {
		return errors.New("话题配置修改错误")
	}
	return nil
}

// 重启serverd服务
func (this *Topics) RestartServerd(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 读取调度器配置信息
	address, err1 := internal.CFG.String("apis", "dispatchd_address")
	port, err2 := internal.CFG.String("apis", "dispatchd_port")
	if err1 != nil || err2 != nil {
		return errors.New("读取调度器配置信息错误" + err1.Error())
	}
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure(), grpc.WithTimeout(30*time.Second))
	if err != nil {
		internal.LogFile.E("调度器服务未开启" + err.Error())
		return err
	}
	defer conn.Close()
	c := proto.NewReportServerdServiceClient(conn)
	// 发起请求
	req, err := c.SayRestartServerd(context.Background(), &proto.RestartRequest{})
	if err != nil {
		internal.LogFile.E("同步话题配置，发起请求错误：" + err.Error())
		return err
	}
	if req.Code != "0" {
		internal.LogFile.E("同步话题配置,错误码：" + req.Code)
		return errors.New("同步话题配置错误")
	}

	*response = map[string]string{
		"message": "同步话题配置成功:" + req.Message,
	}
	return nil
}

// nsq所有话题列表（调用nsq接口）
func (this *Topics) NsqTopics(r *http.Request, args *map[string]interface{}, response *Response) error {
	return nil
}
