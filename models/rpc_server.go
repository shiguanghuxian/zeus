package models

import "53it.net/zues/internal"

type RpcServer struct {
	Id   int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name string `json:"name" xorm:"VARCHAR(30)"`
	Url  string `json:"url" xorm:"VARCHAR(255)"`
}

func (this *RpcServer) TableName() string {
	return "zn_rpc_server"
}

// 获取全部
func GetAllRpcServerList() ([]RpcServer, error) {
	var list []RpcServer
	err := dbEngine().Find(&list)
	if err != nil {
		internal.LogFile.E("查询rpc服务列表:" + err.Error())
		return list, err
	}
	return list, nil
}
