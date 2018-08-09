package dispatchd

import (
	"log"
	"net"

	"53it.net/zues/proto"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*请求调度器功能*/

// 开启一个服务
func (d *Dispatchd) Run() {
	// 监听地址
	listenAddress := d.Address + ":" + d.Port
	// 监听端口
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		internal.LogFile.E("创建tcp监听失败:" + err.Error())
		log.Println("创建tcp监听失败:" + err.Error())
		panic(err)
	}
	s := grpc.NewServer()
	// 注册服务似乎可以放一个里边。。。
	proto.RegisterDispatchConfigServiceServer(s, d)
	proto.RegisterReportServerdServiceServer(s, d)
	go s.Serve(lis)
}

// rpc获取会话配置
func (d *Dispatchd) SayDispatchConfig(ctx context.Context, in *proto.DispatchRequest) (*proto.DispatchReply, error) {
	// 查询
	dispatchReply := &proto.DispatchReply{Code: "1", Message: "未查询到数据"}
	var list []*models.TopicsConfig
	var err error
	if in.Topics == "" {
		list, err = models.GetAllTopicsConfig()
	} else {
		list, err = models.GetWhereTopicsConfig(in.Topics)
	}

	if err != nil {
		dispatchReply.Message = "执行查询出现错误：" + err.Error()
		return dispatchReply, err
	}
	// 如果配置为空
	if len(list) < 1 {
		dispatchReply.Message = "未查询到话题信息"
		return dispatchReply, nil
	}
	// 涉及格式转换问题
	var list1 []*proto.TopicsConfig
	for _, v := range list {
		// 查询解析列表
		ruleList, err := models.GetTCIdEnableTopicsConfigRuleList(int(v.Id))
		if err != nil {
			dispatchReply.Message = "查询解析规则列表错误"
			return dispatchReply, err
		}
		// 转换格式
		var ruleListProto []*proto.TopicsConfigRule
		for _, vv := range ruleList {
			ruleListProto = append(ruleListProto, &proto.TopicsConfigRule{
				Id:         int32(vv.Id),
				Mapped:     vv.Mapped,
				TextUnType: vv.TextUnType,
				TextUnRule: vv.TextUnRule,
				DateFormat: vv.DateFormat,
				Appname:    vv.AppName,
			})
		}
		// 组合最终数据
		list1 = append(list1, &proto.TopicsConfig{
			Id:           v.Id,
			Topics:       v.Topics,
			Channel:      v.Channel,
			ChannelCount: v.ChannelCount,
			DataType:     v.DataType,
			RuleList:     ruleListProto,
		})
	}
	dispatchReply = &proto.DispatchReply{Code: "0", Message: "配置列表获取成功", Data: list1}
	return dispatchReply, nil
}

// serverd启动上报用于重启start
func (d *Dispatchd) SayReportServerd(ctx context.Context, in *proto.ServerdConfigRequest) (*proto.DispatchReply, error) {
	dispatchReply := &proto.DispatchReply{Code: "1", Message: "服务错误"}
	if in.ServerId == "" || in.IpAddress == "" || in.Port == "" {
		dispatchReply.Message = "参数错误"
		return dispatchReply, nil
	}
	serverdConfig[in.ServerId] = in
	// ok
	dispatchReply = &proto.DispatchReply{Code: "0", Message: "成功"}
	return dispatchReply, nil
}

// serverd停止上报用于重启stop
func (d *Dispatchd) SayReportServerdStop(ctx context.Context, in *proto.ServerdConfigRequest) (*proto.DispatchReply, error) {
	dispatchReply := &proto.DispatchReply{Code: "1", Message: "服务错误"}
	if in.ServerId == "" || in.IpAddress == "" || in.Port == "" {
		dispatchReply.Message = "参数错误"
		return dispatchReply, nil
	}
	delete(serverdConfig, in.ServerId)
	// ok
	dispatchReply = &proto.DispatchReply{Code: "0", Message: "成功"}
	return dispatchReply, nil
}

// 重启所有serverd服务
func (d *Dispatchd) SayRestartServerd(ctx context.Context, in *proto.RestartRequest) (*proto.DispatchReply, error) {
	dispatchReply := &proto.DispatchReply{Code: "1", Message: "服务错误"}
	// 循环重启服务
	var err error
	serverids := ""
	for _, v := range serverdConfig {
		err = d.RestartServerd(v.IpAddress + ":" + v.Port)
		serverids += ";" + v.ServerId + ":" + v.IpAddress
		if err != nil {
			break
		}
	}
	if err != nil {
		dispatchReply.Message = "重启服务出现错误：" + err.Error()
	} else {
		dispatchReply = &proto.DispatchReply{Code: "0", Message: "重启成功：" + serverids}
	}
	return dispatchReply, nil
}

// 重启statisd服务
func (d *Dispatchd) SayRestartStatisd(ctx context.Context, in *proto.RestartRequest) (*proto.DispatchReply, error) {
	dispatchReply := &proto.DispatchReply{Code: "1", Message: "服务错误"}
	// 读取配置文件
	rpcAddress, err := internal.CFG.String("dispatchd", "statisd_address")
	if err != nil {
		dispatchReply.Message = "重启服务出现错误：" + err.Error()
		return dispatchReply, err
	}
	rpcPort, err := internal.CFG.String("dispatchd", "statisd_port")
	if err != nil {
		dispatchReply.Message = "重启服务出现错误：" + err.Error()
		return dispatchReply, err
	}
	// 重启服务
	err = d.RestartStatisd(rpcAddress + ":" + rpcPort)
	if err != nil {
		dispatchReply.Message = "重启服务出现错误：" + err.Error()
	} else {
		dispatchReply = &proto.DispatchReply{Code: "0", Message: "重启成功"}
	}
	return dispatchReply, nil
}

// 获取appname列表
func (d *Dispatchd) SayGetAppNameList(ctx context.Context, in *proto.AppNameRequest) (*proto.AppNameReply, error) {
	// 查询
	reply := &proto.AppNameReply{Code: "1", Message: "未查询到数据"}
	var list []*models.AppnameFieldType
	var err error
	if in.Name == "" {
		list, err = models.GetNotDeleteAppnameFieldTypeList()
	} else {
		list, err = models.GetNameNotDeleteAppnameFieldTypeList(in.Name)
	}
	if err != nil {
		reply.Message = "执行查询出现错误：" + err.Error()
		return reply, err
	}
	// 如果配置为空
	if len(list) < 1 {
		reply.Message = "未查询到话题信息"
		return reply, nil
	}
	var list1 []*proto.AppNameInfo
	for _, v := range list {
		list1 = append(list1, &proto.AppNameInfo{
			Id:      v.Id,
			AppName: v.AppName,
			Field:   v.Field,
			Type:    v.Type,
			Unit:    v.Unit,
			Index:   v.Index,
		})
	}
	reply = &proto.AppNameReply{Code: "0", Message: "appname列表获取成功", Data: list1}
	return reply, nil
}

// SayEventSetingsConfig 获取告警配置文件
func (d *Dispatchd) SayEventSetingsConfig(ctx context.Context, in *proto.EventConfigRequest) (*proto.EventConfigReply, error) {
	// 定义返回数据
	reply := &proto.EventConfigReply{Code: "1", Message: "未查询到数据"}
	// 查询数据
	// var list []models.EventSetingTemplate
	// var err error
	list, err := models.GetEventSetingAll(in.Field)
	if err != nil {
		return reply, err
	}
	// 格式转换
	var dataList []*proto.EventSeting
	for _, v := range list {
		one := &proto.EventSeting{
			Id:              int32(v.Id),
			Name:            v.Name,
			AppName:         v.AppName,
			Field:           v.Field,
			ValueType:       v.ValueType,
			Describe:        v.Describe,
			ContinuedTime:   int32(v.ContinuedTime),
			CycleTime:       v.CycleTime,
			TemplateName:    v.TemplateName,
			TemplateContent: v.TemplateContent,
		}
		// 查询规则列表
		var eventRuleList []*proto.EventRuleLevel
		list1, err1 := models.GetSetingIdEventRuleLevelList(v.Id)
		if err1 == nil {
			for _, vv := range list1 {
				// 一个规则
				oneEventRuleLe := &proto.EventRuleLevel{
					Id:            int32(vv.Id),
					EventLevelId:  int32(vv.EventLevelId),
					EventSetingId: int32(vv.EventSetingId),
					Value:         vv.Value,
					Expression:    vv.Expression,
					Sort:          int32(vv.Sort),
					LevelName:     vv.Name,
					Level:         int32(vv.Level),
					Unit:          vv.Unit,
				}
				/* 这里查询出单位转换算法 */
				if vv.Unit != "" {
					//查询字段是否设置单位
					appnameInfo, err := models.GetAppnameFieldAppNameInfo(v.AppName, v.Field)
					if err == nil {
						//查询单位转换
						systemUnitConversionInfo, err := models.GetDoubleUnitSystemUnitConversion(vv.Unit, appnameInfo.Unit)
						if err == nil {
							oneEventRuleLe.SystemUnitConversionInfo = &proto.SystemUnitConversion{
								Id:           int32(systemUnitConversionInfo.Id),
								OriginalUnit: systemUnitConversionInfo.OriginalUnit,
								AfterUnit:    systemUnitConversionInfo.AfterUnit,
								Multiple:     systemUnitConversionInfo.Multiple,
								LuaCode:      systemUnitConversionInfo.LuaCode,
								Type:         int32(systemUnitConversionInfo.Type),
							}
						}
					}
				}
				// 告警规则追加信息
				eventRuleList = append(eventRuleList, oneEventRuleLe)
			}
		}
		// 添加规则列表
		one.EventRuleList = eventRuleList
		// 推送列表
		var eventDeviceList []*proto.EventDevice
		list2, err := models.GetRpcESIDEventDeviceAll(v.Id)
		if err == nil {
			for _, vvv := range list2 {
				eventDeviceList = append(eventDeviceList, &proto.EventDevice{
					Id:            int32(vvv.Id),
					EventSetingId: int32(vvv.EventSetingId),
					DeviceId:      int32(vvv.DeviceId),
					HostName:      vvv.HostName,
					Ip:            vvv.Ip,
					DeviceType:    vvv.DeviceType,
					GroupName:     vvv.GroupName,
				})
			}
		}
		one.EventDeviceList = eventDeviceList
		// 查询推送列表
		var eventPushList []*proto.EventPush
		list3, err := models.GetAllEventPushList(v.Id)
		if err == nil {
			for _, vvvv := range list3 {
				eventPushList = append(eventPushList, &proto.EventPush{
					Id:            int32(vvvv.Id),
					EventSetingId: int32(vvvv.EventSetingId),
					Url:           vvvv.Url,
					Name:          vvvv.Name,
					DataType:      int32(vvvv.DataType),
				})
			}
		}
		one.EventPushList = eventPushList
		// 添加最后数据
		dataList = append(dataList, one)
	}

	reply = &proto.EventConfigReply{Code: "0", Message: "成功", Data: dataList}
	return reply, nil
}
