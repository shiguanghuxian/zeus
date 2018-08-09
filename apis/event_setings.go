package apis

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/proto"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
)

type Setings struct {
	Apis
}

// 告警列表
func (this *Setings) EventList(r *http.Request, args *map[string]interface{}, response *Response) error {
	var totalRows int64 // 总行数
	var err error

	keyword := this.ToString((*args)["keyword"]) // 关键次参数
	totalRows, err = models.GetKeywordEventSetingCount(keyword)
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
	list, err := models.GetKeywordEventSetingList(page, listRows, keyword)
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

// 修改启用状态
func (this *Setings) ChangeEnable(r *http.Request, args *map[string]interface{}, response *Response) error {
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
	_, err = models.UpdateEventSetingEnable(id, enable)
	if err != nil {
		return errors.New("修改告警配置状态错误")
	}
	return nil
}

// 添加告警配置
func (this *Setings) AddSetingEvent(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	name := this.ToString((*args)["name"])
	appName := this.ToString((*args)["app_name"])
	field := this.ToString((*args)["field"])
	if name == "" || appName == "" || field == "" {
		return errors.New("告警标题、APPNAM和字段名 不能为空")
	}
	// 话题models对象
	eventSeting := new(models.EventSeting)
	eventSeting.Name = name
	eventSeting.AppName = appName
	eventSeting.Field = field
	eventSeting.ValueType = this.ToString((*args)["value_type"])
	// // 数值型判断
	// continuedCount, err := this.ToInt((*args)["continued_count"], 0)
	// if err != nil || continuedCount == 0 {
	// 	return errors.New("出现次数输入错误")
	// }
	// eventSeting.ContinuedCount = continuedCount
	// 步长时间
	continuedTime, err := this.ToInt((*args)["continued_time"], 0)
	if err != nil || continuedTime == 0 {
		return errors.New("步长时间输入错误")
	}
	eventSeting.ContinuedTime = continuedTime
	// 执行周期
	cycleTime := this.ToString((*args)["cycle_time"])
	_, err = cron.Parse(cycleTime)
	if err != nil {
		return err
	}
	eventSeting.CycleTime = cycleTime
	//  启用情况
	enable, _ := this.ToInt32((*args)["enable"], 0)
	eventSeting.Enable = enable
	// 描述
	eventSeting.Describe = this.ToString((*args)["describe"])

	// 执行插入
	_, err = models.AddOneEventSeting(eventSeting)
	if err != nil {
		return errors.New("告警配置添加错误")
	}
	return nil
}

// 删除告警设置
func (this *Setings) DelSetingsEvent(r *http.Request, args *map[string]interface{}, response *Response) error {
	ids := this.ToString((*args)["ids"])
	if ids == "" {
		return errors.New("参数错误")
	}
	ids = strings.Trim(ids, ",")
	_, err := models.DelIdsEventSeting(ids)
	if err != nil {
		return errors.New("告警配置删除错误")
	}
	return nil
}

// 获取单条消息
func (this *Setings) InfoSetingsEvent(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	info, err := models.GetOneSetingsEventInfo(id)
	if err != nil {
		return errors.New("获取setings_event信息错误")
	}
	*response = info
	return nil
}

// 保存编辑告警
func (this *Setings) UpSetingEvent(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	name := this.ToString((*args)["name"])
	appName := this.ToString((*args)["app_name"])
	field := this.ToString((*args)["field"])
	if name == "" || appName == "" || field == "" {
		return errors.New("告警标题、APPNAM和字段名 不能为空")
	}
	// 话题models对象
	eventSeting := new(models.EventSeting)
	eventSeting.Name = name
	eventSeting.AppName = appName
	eventSeting.Field = field
	eventSeting.ValueType = this.ToString((*args)["value_type"])
	// // 数值型判断
	// continuedCount, err := this.ToInt((*args)["continued_count"], 0)
	// if err != nil || continuedCount == 0 {
	// 	return errors.New("出现次数输入错误")
	// }
	// eventSeting.ContinuedCount = continuedCount
	// 步长时间
	continuedTime, err := this.ToInt((*args)["continued_time"], 0)
	if continuedTime == 0 {
		return errors.New("步长时间输入错误")
	}
	eventSeting.ContinuedTime = continuedTime
	// 执行周期，定时执行
	cycleTime := this.ToString((*args)["cycle_time"])
	// cron 包的验证有出入，当出现0时，会报错
	if strings.Index(cycleTime+" ", "*/0 ") >= 0 {
		return errors.New("执行周期输入错误:不可以使用*/0")
	}
	if strings.Index(cycleTime+" ", "*/00 ") >= 0 {
		return errors.New("执行周期输入错误:不可以使用*/00")
	}
	_, err = cron.Parse(cycleTime)
	if err != nil {
		return errors.New("执行周期输入错误:" + err.Error())
	}
	eventSeting.CycleTime = cycleTime

	//  启用情况
	enable, _ := this.ToInt32((*args)["enable"], 0)
	eventSeting.Enable = enable
	// 描述
	eventSeting.Describe = this.ToString((*args)["describe"])
	// log.Println(eventSeting)
	// 执行插入
	_, err = models.UpdateIdSetingsEventInfo(id, eventSeting)
	if err != nil {
		return errors.New("告警配置编辑错误")
	}
	return nil
}

// 同步配置
func (this *Setings) SynchroConfigure(r *http.Request, args *map[string]interface{}, response *Response) error {
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
	req, err := c.SayRestartStatisd(context.Background(), &proto.RestartRequest{})
	if err != nil {
		internal.LogFile.E("同步告警配置，发起请求错误：" + err.Error())
		return err
	}
	if req.Code != "0" {
		internal.LogFile.E("同步告警配置,错误码：" + req.Code)
		return errors.New("同步告警配置错误")
	}

	*response = map[string]string{
		"message": "同步告警配置成功:" + req.Message,
	}
	return nil
}
