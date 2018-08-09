package internal

import (
	"errors"
	"time"

	"53it.net/zues/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// appname
type AppNames struct {
	Id     int32
	Name   string             // appname名
	Fields map[string]*Fields // 字段列表
}

// 字段
type Fields struct {
	Field string // 字段名
	Type  string // 数据类型
	Unit  string // 单位
	Index int32  // 是否索引
}

// 保存appname列表
var AppnameList map[string]*AppNames

func InitAppNameList(address string) error {
	AppnameList = make(map[string]*AppNames, 0)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(24*time.Second))
	if err != nil {
		return errors.New("调度器dispatchd未开启服务：" + err.Error())
	}
	defer conn.Close()
	c := proto.NewDispatchConfigServiceClient(conn)
	// 发起请求
	r, err := c.SayGetAppNameList(context.Background(), &proto.AppNameRequest{})
	if err != nil {
		return errors.New("从调度器dispatchd获取appname列表错误：" + err.Error())
	}
	if r.Code != "0" {
		return errors.New("获取失败,错误码：" + r.Code)
	}
	// js, _ := json.Marshal(r.Data)
	// log.Println(string(js))
	// 数据格式转换
	for _, v := range r.Data {
		if AppnameList[v.AppName] == nil {
			AppnameList[v.AppName] = &AppNames{
				Id:   v.Id,
				Name: v.AppName,
			}
		}
		if len(AppnameList[v.AppName].Fields) < 1 {
			AppnameList[v.AppName].Fields = make(map[string]*Fields, 0)
		}
		AppnameList[v.AppName].Fields[v.Field] = &Fields{
			Field: v.Field,
			Type:  v.Type,
			Unit:  v.Unit,
			Index: v.Index,
		}
	}
	return nil
}
