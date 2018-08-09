package serverd

import (
	"encoding/json"
	"strings"

	"53it.net/zues/internal"

	"github.com/bitly/go-nsq"
)

// 消息对象
type Message struct {
	msgchan  chan *nsq.Message // nsq消息对象
	Stop     bool              // 停止从通道获取数据
	DataType string            // 数据类型－当前只有json和文本(单行)
}

// 消息body
type MessageBody struct {
	Ip         string        `json:"ip"`          // 客户机ip
	Group      string        `json:"group"`       // 客户机分组－默认default
	Hostname   string        `json:"hostname"`    // 主机名
	Tag        string        `json:"tag"`         // 数据标签，此处可以用语选择appname（不为空时）
	DeviceType string        `json:"device_type"` // 设备类型
	Data       []interface{} `json:"data"`        // 具体数据部分
}

// 处理消息体方法
func (m *Message) HandleMessage(message *nsq.Message) error {
	if !m.Stop {
		m.msgchan <- message
	}
	return nil
}

// 初步解析json-最外层共用结构部分
func (m *Message) UnMessageShare(jsonByte []byte) (*MessageBody, error) {
	// 去除空白
	jsonByte = []byte(strings.TrimSpace(string(jsonByte)))
	//	fmt.Println(string(jsonByte))
	magBody := new(MessageBody)
	if err := json.Unmarshal(jsonByte, magBody); err != nil {
		internal.LogFile.W("初步解析json失败" + err.Error())
		return magBody, err
	}
	return magBody, nil
}
