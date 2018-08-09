package serverd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"53it.net/zues/elasticsearch"
	"53it.net/zues/influxdb"
	"53it.net/zues/internal"
	"53it.net/zues/mongo"
	"53it.net/zues/proto"
	"53it.net/zues/redis"

	"github.com/bitly/go-nsq"
)

// serverd 对象
type SERVERD struct {
	*Message                                    // 消息体-message.go中定义
	NSQLookupdAddress string                    // nsq消息队列的lookupd地址
	NSQLookupdPort    string                    // nsq消息队列的lookupd端口
	Topics            string                    // 话题－代表这一类上传数据
	Channel           string                    // 通道－处理上传数据的协程
	ChannelCount      int                       // 通道数量
	TopicsRuleList    []*proto.TopicsConfigRule // 解析规则列表

	consumer *nsq.Consumer // 存储当前消费者
}

// 初始化要做的事
func init() {
	internal.NewLog("serverd")
}

// 启动处理方法
func (s *SERVERD) Run() {
	// 消费者对象
	consumer, err := nsq.NewConsumer(s.Topics, s.Channel, nsq.NewConfig())
	if err != nil {
		internal.LogFile.E("创建nsq消费者失败，错误：" + err.Error())
		panic(err)
	}
	// 初始化管道
	s.msgchan = make(chan *nsq.Message, 2048)
	consumer.AddHandler(nsq.HandlerFunc(s.HandleMessage))
	// 建立链接
	err = consumer.ConnectToNSQLookupd(s.NSQLookupdAddress + ":" + s.NSQLookupdPort)
	if err != nil {
		internal.LogFile.E("与NSQLookupd建立链接失败，错误：" + err.Error())
		panic(err)
	}
	// 保存消费者对象
	s.consumer = consumer
	// 开始处理消息-此处应该没有压力，因为nsq消息消费速度很快＊＊＊
	s.Process()
}

// 停止获取数据
func (s *SERVERD) StopRun() {
	// close(s.msgchan) // 这里重启会出问题，不关闭了
	s.msgchan = nil
	// s.Stop = true
	s.consumer.Stop()
}

// 取数据
func (s *SERVERD) Process() {
	for {
		select {
		case message := <-s.msgchan:
			// 解析消息体元数据
			var magBody *MessageBody
			var err error
			if message != nil {
				magBody, err = s.UnMessageShare(message.Body)
				if err == nil {
					// 判断是否存在：和“
					hostinfos := magBody.Ip + magBody.Hostname + magBody.Group
					if strings.Index(hostinfos, ":") == -1 || strings.Index(hostinfos, "\"") == -1 {
						// 去空格
						magBody.Ip = strings.TrimSpace(magBody.Ip)
						magBody.Hostname = strings.TrimSpace(magBody.Hostname)
						magBody.Group = strings.TrimSpace(magBody.Group)
						magBody.Tag = strings.TrimSpace(magBody.Tag)
						if magBody.DeviceType = strings.TrimSpace(magBody.DeviceType); magBody.DeviceType == "" {
							magBody.DeviceType = "default"
						}
						// magBody.DeviceType = strings.TrimSpace(magBody.DeviceType)
						if magBody.Ip == "" && magBody.Hostname == "" && magBody.Group == "" {
							internal.LogFile.W("设备信息为空，无法定位存储数据", s.Topics, s.Channel)
						} else {
							// 处理数据
							switch s.DataType {
							case "json":
								s.handleJson(magBody)
								break
							case "text":
								s.handleText(magBody)
								break
							}
						}
					} else {
						internal.LogFile.W("主机信息包含非法字符:", hostinfos, s.Topics, s.Channel)
					}
				} else {
					// log.Println(time.Now().Format("2006-01-02 15:04:05"), ":", err.Error())
					internal.LogFile.W("解包错误，采集端上传包格式不符合标准", s.Topics, s.Channel)
				}
			}
		case <-time.After(time.Second):
			if s.Stop {
				close(s.msgchan)
				return
			}
		}
	}
}

// 处理json数据
func (s *SERVERD) handleJson(msg *MessageBody) {
	if len(msg.Data) < 1 {
		return
	}
	// 循环获取多条结果集
	for _, v := range msg.Data {
		// 转换数据格式
		if vv, ok := v.(map[string]interface{}); ok == true {
			// 寻找匹配的字段影射
			for _, rv := range s.TopicsRuleList {
				if msg.Tag != "" && msg.Tag != rv.Tag {
					continue
				}
				// 数据下标影射key
				keys := s.getMappedRuleStr(rv.Mapped)
				// 判断下标对应
				if ok := internal.VerifyMappedJsonKeys(vv, keys); ok == true {
					// 定义存储的map
					rawData := make(map[string]interface{})
					// 基本数据
					rawData["ip"] = msg.Ip
					rawData["group"] = msg.Group
					rawData["hostname"] = msg.Hostname
					rawData["device_type"] = msg.DeviceType

					// 拆包数据
					for dk, dv := range keys {
						if dk == "date" {
							rawData["datetime"] = internal.DateStrToint64(fmt.Sprint(vv[dv]), rv.DateFormat)
						}
						rawData[dk] = vv[dv]
					}
					// 执行插入
					s.saveOneRawData(&rawData, rv.Appname)
					break
				}
			}
		} else {
			internal.LogFile.W("json数据转换格式失败：" + fmt.Sprint(v))
		}
	}
}

// 处理text数据
func (s *SERVERD) handleText(msg *MessageBody) {
	if len(msg.Data) < 1 {
		return
	}
	// 处理每个数据包中的数据
	for _, v := range msg.Data {
		// 寻找匹配的字段影射
		for _, rv := range s.TopicsRuleList {
			if msg.Tag != "" && msg.Tag != rv.Tag {
				continue
			}
			// 定义存储的map
			rawData := make(map[string]interface{})
			// 基本数据
			rawData["ip"] = msg.Ip
			rawData["group"] = msg.Group
			rawData["hostname"] = msg.Hostname
			rawData["device_type"] = msg.DeviceType

			// 在这里处理规则，减少处理次数
			switch rv.TextUnType {
			case "char":
				// 字段影射
				var keys map[string]int
				keys = s.getMappedRule(rv.Mapped)
				// 判断是否是空白字符
				var strList []string
				if strings.TrimSpace(rv.TextUnRule) == "" {
					strList = strings.Fields(fmt.Sprint(v))
				} else {
					strList = strings.Split(fmt.Sprint(v), rv.TextUnRule)
				}
				// 识别是否符合规则
				if len(keys) != len(strList) {
					break
				}
				// 拆包数据
				for dk, dv := range keys {
					if dk == "date" {
						rawData["datetime"] = internal.DateStrToint64(fmt.Sprint(strList[dv]), rv.DateFormat)
					}
					rawData[dk] = strList[dv]
				}
				// 执行插入
				s.saveOneRawData(&rawData, rv.Appname)
				break
			case "regular":
				// 字段影射
				var keys1 map[string]string
				keys1 = s.getMappedRuleStr(rv.Mapped)
				if strList, err := internal.RegStrToMap(rv.TextUnRule, fmt.Sprint(v)); err == nil {
					// 判断key是否合法
					if ok := internal.VerifyMappedRegularKeys(strList, keys1); ok == true {
						// 拆包数据
						for dk, dv := range keys1 {
							if dk == "date" {
								rawData["datetime"] = internal.DateStrToint64(fmt.Sprint(strList[dv]), rv.DateFormat)
							}
							rawData[dk] = strList[dv]
						}
					}
				}
				// 执行插入
				s.saveOneRawData(&rawData, rv.Appname)
				break
			}
		}
	}
}

//// 保存数据
//func (s *SERVERD) saveRawData(rawDataAll []map[string]interface{}, appname string) error {
//	if len(rawDataAll) < 1 {
//		return nil
//	}
//	// 插入redis 最新数据
//	go redis.SaveNewestData(rawDataAll, appname)
//	// 掉用插入mongodb
//	enable, _ := internal.CFG.Bool("mongodb", "enable")
//	if enable {
//		go mongo.AddRawDataAll(rawDataAll, appname)
//	}
//	return nil
//}

// 插入单条数据
func (s *SERVERD) saveOneRawData(rawData *map[string]interface{}, appname string) {
	// 保证datetime一定存在
	if (*rawData)["datetime"] == nil {
		(*rawData)["datetime"] = time.Now().Local().Unix()
	}
	// 获取设备id
	rKey := fmt.Sprintf("deviceids:%s:%s:%s", (*rawData)["group"], (*rawData)["hostname"], (*rawData)["ip"])
	rval, err := redis.GetKeyVal(rKey)
	if err != nil || rval == "" || rval == "0" {
		if rval != "0" {
			// 设备发现
			redis.SaveDeviceInfo((*rawData))
		}
		return
	}
	go func(rd *map[string]interface{}) {
		// rd1 := *rd
		// 插入redis 最新数据
		redis.SaveRawDataOne(rd, appname)
		// 掉用插入mongodb
		enable, _ := internal.CFG.Bool("mongodb", "enable")
		if enable {
			mongo.AddRawDataOne(rd, appname)
		}
		// 是否插入elasticsearch
		enable2, _ := internal.CFG.Bool("elasticsearch", "enable")
		if enable2 {
			elasticsearch.AddRawDataOne(rd, appname)
		}
		// 调用插入influxdb
		enable1, _ := internal.CFG.Bool("influxdb", "enable")
		if enable1 {
			influxdb.AddRawDataOne(rd, appname)
		}
	}(rawData)
	// log.Println(rawData)
}

// 当是char时的字段对应
func (s *SERVERD) getMappedRule(mapped string) map[string]int {
	strList := strings.Split(mapped, "|")
	var mappedMap map[string]int
	mappedMap = make(map[string]int)
	for k, v := range strList {
		strL := strings.Split(v, ":")
		key, err := strconv.Atoi(strL[0])
		if err != nil {
			internal.LogFile.W("解析字段配置错误：" + err.Error())
			key = k
		}
		mappedMap[strL[1]] = key
	}
	return mappedMap
}

// 当是regular时的字段对应
func (s *SERVERD) getMappedRuleStr(mapped string) map[string]string {
	strList := strings.Split(mapped, "|")
	var mappedMap map[string]string
	mappedMap = make(map[string]string)
	for _, v := range strList {
		strL := strings.Split(v, ":")
		mappedMap[strL[1]] = strL[0]
	}
	return mappedMap
}
