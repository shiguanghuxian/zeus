package statisd

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"fmt"

	"53it.net/zues/internal"
	"53it.net/zues/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service struct {
	Lis        net.Listener // tcp监听
	grpcServer *grpc.Server // grpc server
}

// 远程调用重启
func (s *Service) SayReStartStatisd(ctx context.Context, in *proto.StatisdRequest) (*proto.StatisdReply, error) {
	msg := &proto.StatisdReply{Code: "1", Message: "服务错误"}
	// 停止
	mySTATISD.Stop()
	// 重新赋值告警配置列表---rpc调用
	var err error
	mySTATISD.EventSetingList, err = GetEventConfig()
	if err != nil {
		internal.LogFile.E("读取远程配置错误" + err.Error())
		msg.Message = "读取远程配置错误"
		return msg, errors.New("重启失败")
	}
	// 重置appnames列表
	dAddress, err := internal.CFG.String("statisd", "dispatchdaddress")
	if err != nil {
		dAddress = "127.0.0.1"
	}
	dPort, err := internal.CFG.String("statisd", "dispatchdport")
	if err != nil {
		dPort = "3200"
	}
	err = internal.InitAppNameList(dAddress + ":" + dPort)
	if err != nil {
		internal.LogFile.W(err.Error())
	}
	// 启动
	mySTATISD.Run()
	// 成功提示
	msg = &proto.StatisdReply{Code: "0", Message: "服务已重启"}
	return msg, nil
}

// 停止grpc service
func (s *Service) StopService() {
	s.grpcServer.Stop()
	s.Lis.Close()
}

// grpc服务
func (s *Service) RunService() error {
	rpcAddress, err := internal.CFG.String("statisd", "rpcaddress")
	rpcPort, err1 := internal.CFG.String("statisd", "rpcport")
	// 判断启动服务的地址和端口不能为空
	if err != nil || err1 != nil {
		return errors.New("读取重启statisd服务配置失败" + err.Error() + err1.Error())
	}
	// 监听端口
	lis, err := net.Listen("tcp", rpcAddress+":"+rpcPort)
	if err != nil {
		return errors.New("创建tcp监听失败:" + err.Error())
	}
	ser := grpc.NewServer()
	// 注册服务似乎可以放一个里边。。。
	proto.RegisterStatisdServiceServer(ser, &Service{})
	// 存储grpc
	s.Lis = lis
	s.grpcServer = ser

	go ser.Serve(lis)
	return nil
}

// 从调度器获取配置
func GetEventConfig() ([]*proto.EventSeting, error) {
	dAddress, err := internal.CFG.String("statisd", "dispatchdaddress")
	if err != nil {
		return nil, err
	}
	dPort, err := internal.CFG.String("statisd", "dispatchdport")
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%s", dAddress, dPort)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return nil, errors.New("调度器dispatchd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewDispatchConfigServiceClient(conn)
	// 发起请求
	r, err := c.SayEventSetingsConfig(context.Background(), &proto.EventConfigRequest{})
	if err != nil {
		return nil, errors.New("从调度器dispatchd获取配置错误：" + err.Error())
	}
	if r.Code != "0" {
		return nil, errors.New("获取配置失败,错误码：" + r.Code)
	}
	// 输出方便调试
	fltB, _ := json.Marshal(r.Data)
	fmt.Println(string(fltB))
	return r.Data, nil
}
