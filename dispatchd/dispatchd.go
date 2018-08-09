package dispatchd

import (
	"53it.net/zues/internal"
	"53it.net/zues/proto"
)

// 调度器服务对象
type Dispatchd struct {
	Address string // 调度器sever 地址
	Port    string // 调度器server 端口
}

// 用于存储所有serverd服务配置
var serverdConfig map[string]*proto.ServerdConfigRequest

// 初始化要做的事
func init() {
	internal.NewLog("dispatchd")
	serverdConfig = make(map[string]*proto.ServerdConfigRequest)
}
