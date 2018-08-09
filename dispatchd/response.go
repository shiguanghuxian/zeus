package dispatchd

import (
	"errors"
	"time"

	"53it.net/zues/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*调度器请求其它组件*/

// 重启serverd服务，grpc
func (d *Dispatchd) RestartServerd(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return errors.New("重启serverd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewServerdServiceClient(conn)
	// 发起请求
	r, err := c.SayRestart(context.Background(), &proto.ServerdRequest{})
	if err != nil {
		return errors.New("重启serverd获取配置错误：" + err.Error())
	}
	if r.Code != "0" {
		return errors.New("重启serverd,错误码：" + r.Code)
	}
	return nil
}

// 重启statisd服务，grpc
func (d *Dispatchd) RestartStatisd(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return errors.New("重启statisd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewStatisdServiceClient(conn)
	// 发起请求
	r, err := c.SayReStartStatisd(context.Background(), &proto.StatisdRequest{})
	if err != nil {
		return errors.New("重启statisd获取配置错误：" + err.Error())
	}
	if r.Code != "0" {
		return errors.New("重启statisd,错误码：" + r.Code)
	}
	return nil
}
