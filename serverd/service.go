package serverd

import (
	"errors"
	"io/ioutil"
	"net"
	"time"

	"53it.net/zues/internal"
	"53it.net/zues/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// 每个处理对象一个实例
var nsqServerList []*SERVERD // nsq 处理列表
var rpcServer Service        // grpc server

type Service struct {
	Lis        net.Listener // tcp监听
	grpcServer *grpc.Server // grpc server
}

func init() {
	nsqServerList = make([]*SERVERD, 0)
}

// 获取nsq服务列表
func GetNsqServerList() []*SERVERD {
	return nsqServerList
}

// 远程调用重启nsq数据处理
func (s *Service) SayRestart(ctx context.Context, in *proto.ServerdRequest) (*proto.ServerdReply, error) {
	msg := &proto.ServerdReply{Code: "1", Message: "服务错误"}
	// 停止
	stopServerd()
	// 清空
	nsqServerList = make([]*SERVERD, 0)
	time.Sleep(1 * time.Second)
	// 启动
	RunServerd()
	if len(nsqServerList) > 0 {
		msg = &proto.ServerdReply{Code: "0", Message: "服务已重启"}
	} else {
		msg.Message = "重启失败"
		return msg, errors.New("重启失败")
	}
	return msg, nil
}

// 停止grpc service
func stopService() {
	rpcServer.grpcServer.Stop()
	rpcServer.Lis.Close()
}

// grpc服务
func RunService() error {
	rpcAddress, err := internal.CFG.String("serverd", "rpcaddress")
	rpcPort, err1 := internal.CFG.String("serverd", "rpcport")
	// 判断启动服务的地址和端口不能为空
	if err != nil || err1 != nil {
		return errors.New("读取重启serverd服务配置失败" + err.Error() + err1.Error())
	}
	// 监听端口
	lis, err := net.Listen("tcp", rpcAddress+":"+rpcPort)
	if err != nil {
		return errors.New("创建tcp监听失败:" + err.Error())
	}
	s := grpc.NewServer()
	// 注册服务似乎可以放一个里边。。。
	proto.RegisterServerdServiceServer(s, &Service{})
	// 存储grpc
	rpcServer.Lis = lis
	rpcServer.grpcServer = s

	go s.Serve(lis)
	return nil
}

// 停止serverd服务
func stopServerd() {
	// 关闭所有nsq消费之(serverd)
	for _, v := range nsqServerList {
		v.StopRun()
	}
}

// 运行serverd服务
func RunServerd() []*SERVERD {
	// 重新获取配置文件
	errconf := internal.ResetConfig()
	if errconf != nil {
		internal.LogFile.W("配置文件重新读取错误" + errconf.Error())
	}
	// 读取nsq配置
	nsqAddress, err1 := internal.CFG.String("nsq", "nsqaddress")
	nsqPort, err2 := internal.CFG.String("nsq", "nsqport")
	if err1 != nil || err2 != nil {
		internal.LogFile.E("读取nsq配置错误：" + err1.Error())
		panic(err1)
	}
	// 从调度器获取话题和通道配置
	dAddress, err3 := internal.CFG.String("serverd", "dispatchdaddress")
	dPort, err4 := internal.CFG.String("serverd", "dispatchdport")
	if err3 != nil || err4 != nil {
		internal.LogFile.E("读取dispatchd配置错误：" + err3.Error())
		panic(err3)
	}
	// 上报serverd配置
	err := reportRpcConfig(dAddress + ":" + dPort)
	if err != nil {
		internal.LogFile.E(err.Error())
		panic(err)
	}
	// 获取appname列表
	err = internal.InitAppNameList(dAddress + ":" + dPort)
	if err != nil {
		internal.LogFile.W(err.Error())
	}
	// 获取配置
	cfg, err := getConfig(dAddress + ":" + dPort)
	if err != nil {
		internal.LogFile.E(err.Error())
		panic(err)
	}
	// 根据配置创建对应通道处理数据
	for _, v := range cfg {
		// 每个话题几每个serverd运行几个通道
		for i := 0; i < int(v.ChannelCount); i++ {
			nsqServer := &SERVERD{
				NSQLookupdAddress: nsqAddress,
				NSQLookupdPort:    nsqPort,
				Topics:            v.Topics,
				Channel:           v.Channel,
				ChannelCount:      int(v.ChannelCount),
				TopicsRuleList:    v.RuleList, // 直接用了proto中数据类型
				Message: &Message{
					DataType: v.DataType,
					Stop:     false,
				},
			}
			// 保存serverd对象，便于停止服务
			nsqServerList = append(nsqServerList, nsqServer)
			go nsqServer.Run()
		}
	}

	return nsqServerList
}

// 从调度器获取配置
func getConfig(address string) ([]*proto.TopicsConfig, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return nil, errors.New("调度器dispatchd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewDispatchConfigServiceClient(conn)
	// 发起请求
	r, err := c.SayDispatchConfig(context.Background(), &proto.DispatchRequest{})
	if err != nil {
		return nil, errors.New("从调度器dispatchd获取配置错误：" + err.Error())
	}
	if r.Code != "0" {
		return nil, errors.New("获取配置失败,错误码：" + r.Code)
	}
	return r.Data, nil
}

// 上报配置信息，便于重启服务
func reportRpcConfig(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return errors.New("调度器dispatchd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewReportServerdServiceClient(conn)
	// serverd重启服务
	rpcAddress, _ := internal.CFG.String("serverd", "rpcaddress")
	rpcPort, _ := internal.CFG.String("serverd", "rpcport")
	serverId := getServerdId()
	// fmt.Println(serverId)
	// 发起请求
	r, err := c.SayReportServerd(context.Background(), &proto.ServerdConfigRequest{ServerId: serverId, IpAddress: rpcAddress, Port: rpcPort})
	if err != nil {
		return errors.New("上传配置到调度器dispatchd错误：" + err.Error())
	}
	if r.Code != "0" {
		return errors.New("上传配置到调度器dispatchd错误,错误码：" + r.Code)
	}
	return nil
}

// 上报配置信息，停止，防止调度出问题
func ReportStopRpcConfig(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return errors.New("调度器dispatchd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewReportServerdServiceClient(conn)
	// serverd重启服务
	rpcAddress, _ := internal.CFG.String("serverd", "rpcaddress")
	rpcPort, _ := internal.CFG.String("serverd", "rpcport")
	serverId := getServerdId()
	// fmt.Println(serverId)
	// 发起请求
	r, err := c.SayReportServerdStop(context.Background(), &proto.ServerdConfigRequest{ServerId: serverId, IpAddress: rpcAddress, Port: rpcPort})
	if err != nil {
		return errors.New("上传配置到调度器dispatchd错误：" + err.Error())
	}
	if r.Code != "0" {
		return errors.New("上传配置到调度器dispatchd错误,错误码：" + r.Code)
	}
	return nil
}

// 读取serverd id
func getServerdId() string {
	rootDir := internal.GetRootDir()
	idPath := rootDir + "/serverd.id"
	var serverdid string
	// 读文件
	strByte, err := ioutil.ReadFile(idPath)
	if err != nil {
		serverdid = internal.Rand().Hex()
		err1 := ioutil.WriteFile(idPath, []byte(serverdid), 0666)
		if err1 != nil {
			internal.LogFile.E("写serverd id到文件失败" + err1.Error())
		}
	} else {
		serverdid = string(strByte)
	}
	return serverdid
}
