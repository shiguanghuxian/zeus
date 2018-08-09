package internal

import (
	"encoding/json"

	"github.com/bitly/go-nsq"
)

type NSQMessage struct {
	Ip       string                 `json:"ip"`
	Group    string                 `json:"group"`
	HostName string                 `json:"hostname"`
	Data     map[string]interface{} `json:"data"`
}

// 自监控追加数据
func SendNsqMessage(msg NSQMessage) error {
	// 读取配置文件
	address, _ := CFG.String("nsq_prod", "address")
	port, _ := CFG.String("nsq_prod", "port")
	// 创建生产者对象
	producer, err := nsq.NewProducer(address+":"+port, nsq.NewConfig())
	defer producer.Stop()
	if err != nil {
		LogFile.W("创建nsq生产者失败:" + err.Error())
		return err
	}
	// 转json
	jsonByte, err := json.Marshal(msg)
	if err != nil {
		LogFile.W("转json格式错误:" + err.Error())
		return err
	}
	err = producer.Publish("self_monitor", jsonByte)
	if err != nil {
		LogFile.W("nsq生产者添加数据错误:" + err.Error())
		return err
	}
	return nil
}
