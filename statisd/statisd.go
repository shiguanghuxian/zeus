package statisd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"53it.net/zues/elasticsearch"
	"53it.net/zues/influxdb"
	"53it.net/zues/internal"
	"53it.net/zues/mongo"
	"53it.net/zues/proto"
	"53it.net/zues/redis"
	"github.com/robfig/cron"
)

type STATISD struct {
	EventSetingList []*proto.EventSeting
	crontab         *cron.Cron // 定时器
}

var mySTATISD *STATISD
var dataSource string
var Debug bool

// 初始化要做的事
func init() {
	internal.NewLog("statisd")
}

func NewSTATISD() *STATISD {
	if mySTATISD == nil {
		mySTATISD = new(STATISD)
	}
	return mySTATISD
}

// 启动处理方法
func (s *STATISD) Run() {
	// 创建定时对象
	s.crontab = cron.New()
	var err error
	for _, v := range s.EventSetingList {
		v := v
		err = s.crontab.AddFunc(v.CycleTime, func() {
			s.runTask(time.Now(), v)
		})
		if err != nil {
			internal.LogFile.E("创建定时任务错误" + err.Error())
		}
	}
	s.crontab.Start() // 启动定时任务
}

// 结束定时器
func (s *STATISD) Stop() {
	internal.LogFile.W("停止定时任务")
	s.crontab.Stop()
}

// 定时任务执行--具体任务部分
func (s *STATISD) runTask(t time.Time, eventeting *proto.EventSeting) {
	if Debug == true {
		log.Println(fmt.Sprint(" 时间：", t))
	}
	// 循环设备
	for _, v := range eventeting.EventDeviceList {
		for _, vv := range eventeting.EventRuleList {
			// 单位转换
			if vv.SystemUnitConversionInfo != nil && vv.Unit != "" {
				var err error
				vv.Value, err = internal.ConvertUnit(vv.Unit, vv.Value, vv.SystemUnitConversionInfo)
				if err != nil {
					internal.LogFile.W("告警阈值设置单位转换错误：" + err.Error())
				}
			}
			// 组织zql
			zqlPublic := "select %s from %s where ((group = '%s') and (hostname = '%s') and (ip = '%s')) %s "
			var zql = ""                    // 最终执行的sql
			compareVal := make([]string, 1) // 查询出的结果
			inexistence := false            // 数据是否存在
			limit := 1                      // 查询行数
			var selectField string          // 查询的字段
			switch eventeting.ValueType {
			case "当前值": // 从redis中查询数据
				if vv.Expression == "diff" {
					selectField = eventeting.Field
					limit = 2
				} else {
					selectField = ""
				}
				break
			case "统计值":
				selectField = fmt.Sprintf("count(%s) as compare_val", eventeting.Field)
				break
			case "平均值":
				if s.GetDataSource() == "influxdb" {
					selectField = fmt.Sprintf("mean(%s) as compare_val", eventeting.Field)
				} else {
					selectField = fmt.Sprintf("avg(%s) as compare_val", eventeting.Field)
				}
				break
			case "最大值":
				selectField = fmt.Sprintf("max(%s) as compare_val", eventeting.Field)
				break
			case "最小值":
				selectField = fmt.Sprintf("min(%s) as compare_val", eventeting.Field)
				break
			}
			// 这里具体执行zql查询
			if selectField != "" {
				switch s.GetDataSource() {
				case "influxdb":
					// 计算时间
					startTime := time.Unix(t.Unix()-int64(eventeting.ContinuedTime), 0).Format("2006-01-02 03:04:05")
					endTime := t.Format("2006-01-02 03:04:05")
					// 组织zql
					zqlOwn := fmt.Sprintf(" and ((time > '%s') and (time < '%s')) ", startTime, endTime)
					zql = fmt.Sprintf(zqlPublic, selectField, eventeting.AppName, v.GroupName, v.HostName, v.Ip, zqlOwn)
					zql = fmt.Sprintf("%s order by time desc limit %d ", zql, limit)
					if Debug == true {
						log.Println(zql)
					}
					list, err := influxdb.ZqlQueryCmd(zql)
					if err != nil {
						internal.LogFile.E("定时告警任务查询错误:[influxdb]->", zql)
					}
					if len(list) == 0 {
						inexistence = true
					} else {
						if vv.Expression == "diff" {
							compareVal[0] = fmt.Sprint(list[0][eventeting.Field])
							if len(list) > 1 {
								compareVal = append(compareVal, fmt.Sprint(list[1][eventeting.Field]))
							} else {
								internal.LogFile.W("差值比较，只查询到了一条数据:[influxdb]->", zql)
								continue
							}
						} else {
							compareVal[0] = fmt.Sprint(list[0]["compare_val"])
						}
					}
					break
				case "mongodb":
					// 计算时间
					startTime := time.Unix(t.Unix()-int64(eventeting.ContinuedTime), 0).Format("2006-01-02 03:04:05")
					endTime := t.Format("2006-01-02 03:04:05")
					// 组织zql
					zqlOwn := fmt.Sprintf(" and ((date > date('%s')) and (date < date('%s'))) ", startTime, endTime)
					zql = fmt.Sprintf(zqlPublic, selectField, eventeting.AppName, v.GroupName, v.HostName, v.Ip, zqlOwn)
					zql = fmt.Sprintf("%s group by time(%ds) order by datetime desc limit %d ", zql, eventeting.ContinuedTime, limit)
					if Debug == true {
						log.Println(zql)
					}
					list, err := mongo.GetZqlList(zql)
					if err != nil {
						internal.LogFile.E("定时告警任务查询错误:[mongodb]->", zql)
					}
					if len(list) == 0 {
						inexistence = true
					} else {
						if vv.Expression == "diff" {
							compareVal[0] = fmt.Sprint(list[0][eventeting.Field])
							if len(list) > 1 {
								compareVal = append(compareVal, fmt.Sprint(list[1][eventeting.Field]))
							} else {
								internal.LogFile.W("差值比较，只查询到了一条数据:[mongodb]->", zql)
								continue
							}
						} else {
							compareVal[0] = fmt.Sprint(list[0]["compare_val"])
						}
					}
					break
				case "elasticsearch":
					// 计算时间
					startTime := time.Unix(t.Unix()-int64(eventeting.ContinuedTime), 0).Format("2006-01-02T03:04:05")
					endTime := t.Format("2006-01-02T03:04:05")
					// 组织zql
					zqlOwn := fmt.Sprintf(" and ((date > '%s') and (date < '%s')) ", startTime, endTime)
					zql = fmt.Sprintf(zqlPublic, selectField, eventeting.AppName, v.GroupName, v.HostName, v.Ip, zqlOwn)
					zql = fmt.Sprintf("%s group by time(%ds) order by datetime desc limit %d ", zql, eventeting.ContinuedTime, limit)
					if Debug == true {
						log.Println(zql)
					}
					list, err := elasticsearch.GetZqlList(zql)
					if err != nil {
						internal.LogFile.E("定时告警任务查询错误:[elasticsearch]->", zql)
					}
					if len(list) == 0 {
						inexistence = true
					} else {
						if vv.Expression == "diff" {
							compareVal[0] = fmt.Sprint(list[0][eventeting.Field])
							if len(list) > 1 {
								compareVal = append(compareVal, fmt.Sprint(list[1][eventeting.Field]))
							} else {
								internal.LogFile.W("差值比较，只查询到了一条数据:[elasticsearch]->", zql)
								continue
							}
						} else {
							compareVal[0] = fmt.Sprint(list[0]["compare_val"])
						}
					}
					break
				default:
					internal.LogFile.E("没有可用的数据源")
					return
				}
			} else { // 当前值情况
				rawData := map[string]string{
					"group":    v.GroupName,
					"hostname": v.HostName,
					"ip":       v.Ip,
				}
				var err error
				// 先获取时间
				writeTimeStr, err := redis.GetOneNewestData(rawData, eventeting.AppName, "datetime")
				if err == nil {
					writeTime, _ := strconv.ParseInt(writeTimeStr, 10, 64)
					if t.Unix() < writeTime && (t.Unix()-writeTime) > int64(eventeting.ContinuedTime) {
						if Debug == true {
							log.Println("当前值判断数据不在查询范围内")
						}
						break // 跳出循环
					} else {
						compareVal[0], err = redis.GetOneNewestData(rawData, eventeting.AppName, eventeting.Field)
						if err != nil {
							internal.LogFile.W("redis读取当前值错误：" + err.Error())
							compareVal[0] = ""
						} else {
							inexistence = true
						}
					}
				} else {
					if Debug == true {
						log.Println("获取当前值时间戳错误:" + err.Error())
					}
				}
			}
			// 用户告警模版显示
			compareVal0 := compareVal[0]
			var compareVal1 string
			if len(compareVal) == 2 {
				compareVal1 = compareVal[1]
			}
			// 判断值是否是空--并且表达式不是是否存在情况
			if vv.Expression != "inexistence" {
				if compareVal[0] == "" {
					if Debug == true {
						log.Println("告警未查询到结果，并且表达式不是是否存在情况")
					}
					continue
				}
			} else {
				// 保证可以使用统一行使判断
				if inexistence == true {
					compareVal[0] = "true"
				} else {
					compareVal[0] = "false"
				}
			}
			// 比较值判断
			alarmValue, threshold, err := s.compareThresholdVal(compareVal, vv.Value, vv.Expression)
			if err != nil {
				internal.LogFile.W("阈值比较错误：" + err.Error())
			}
			if threshold == true {
				if Debug == true {
					log.Println("阈值比较结果true")
				}
				// 写入告警历史表
				alarmData := &map[string]interface{}{
					"esid":             eventeting.Id,
					"current_value0":   compareVal0,
					"current_value1":   compareVal1,
					"alarm_value":      alarmValue,
					"date":             t.Format("2006-01-02 15:04:05"),
					"datetime":         t.Unix(),
					"event_level":      vv.Level,
					"event_level_name": vv.LevelName,
					"hostname":         v.HostName,
					"ip":               v.Ip,
					"group":            v.GroupName,
				}
				tplAlarm, err := internal.TplAnalysisToString(eventeting.TemplateContent, alarmData)
				if err != nil {
					internal.LogFile.W("解析告警模版错误：" + err.Error())
				}
				(*alarmData)["message"] = tplAlarm
				err = s.saveAlarmData(alarmData)
				if err != nil {
					if Debug == true {
						log.Println(err.Error())
					}
				}
				// 发送推送告警
				err = s.SendEventToHttp(eventeting.EventPushList, alarmData)
				if err != nil {
					internal.LogFile.E("发送推送告警错误：" + err.Error())
					if Debug == true {
						log.Println("发送推送告警错误：" + err.Error())
					}
				}

				break // 比较出来则跳出--因为规则是有顺序的，只匹配一个即可
			} else {
				if Debug == true {
					log.Println("阈值比较结果false:alarmValue->" + alarmValue + ",threshold->" + vv.Value)
				}
			}
		}
	}
}

// 保存告警信息
func (s *STATISD) saveAlarmData(alarmData *map[string]interface{}) (err error) {
	switch s.GetDataSource() {
	case "influxdb":
		err = influxdb.InsertInto(alarmData, "zn_original_alarm")
		break
	case "mongodb":
		err = mongo.InsertInto(alarmData, "zn_original_alarm")
		break
	case "elasticsearch":
		err = elasticsearch.InsertInto(alarmData, "zn_original_alarm")
		break
	default:
		internal.LogFile.W("未找到数据源")
		err = errors.New("未找到数据源")
	}
	return
}

// 与阈值比较
func (s *STATISD) compareThresholdVal(values []string, threshold, expression string) (value string, bl bool, err error) {
	if threshold == "" || values[0] == "" {
		return "", false, errors.New("阈值和比较值不能为空threshold:" + threshold + ";value=" + values[0])
	}
	if expression == "diff" {
		if len(values) < 2 {
			return "", false, errors.New("差值比较需要最少两个值")
		}
		thresholds := strings.Split(threshold, "|")
		thresholds[0] = strings.TrimSpace(thresholds[0])
		if len(thresholds) != 3 {
			if thresholds[0] != "true" && thresholds[0] != "false" {
				return "", false, errors.New("阈值字段格式错误")
			}
		}
		threshold = thresholds[0] // 阈值取第一个值
		if thresholds[0] == "true" || thresholds[0] == "false" {
			expression = "="
			if values[0] == values[1] {
				value = "true"
			} else {
				value = "false"
			}
		} else {
			thresholds[1] = strings.TrimSpace(thresholds[1])
			thresholds[2] = strings.TrimSpace(thresholds[2])
			expression = thresholds[2]
			// 转换格式
			val1, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				return "", false, err
			}
			val2, err := strconv.ParseFloat(values[1], 64)
			if err != nil {
				return "", false, err
			}
			val3 := val1 - val2
			if thresholds[1] == "0" {
				if val3 < 0 {
					value = fmt.Sprintf("%.2f", val3*-1)
				} else {
					value = fmt.Sprintf("%.2f", val3)
				}
			} else {
				value = fmt.Sprintf("%.2f", val3)
			}
			// 阈值也转换为2位小数
			thresholdDiff, err := strconv.ParseFloat(threshold, 64)
			if err != nil {
				return "", false, err
			}
			threshold = fmt.Sprintf("%.2f", thresholdDiff)
		}
	} else {
		value = values[0]
	}
	// log.Println(value + "||" + threshold + "||" + expression)
	bl = false
	var threshold1 float64
	var value1 float64
	if expression == ">" || expression == ">=" || expression == "<" || expression == "<=" {
		threshold1, err = strconv.ParseFloat(threshold, 64)
		if err != nil {
			return "", false, err
		}
		value1, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return "", false, err
		}
	}
	switch expression {
	case "=":
		bl = (value == threshold)
		break
	case ">":
		if value1 > threshold1 {
			bl = true
		}
		break
	case "<":
		if value1 < threshold1 {
			bl = true
		}
		break
	case ">=":
		if value1 >= threshold1 {
			bl = true
		}
		break
	case "<=":
		if value1 <= threshold1 {
			bl = true
		}
		break
	case "!=":
		bl = (value != threshold)
		break
	case "like":
		if strings.Index(value, threshold) >= 0 {
			bl = true
		}
		break
	case "inexistence":
		bl = (value == threshold)
		break
	}
	return
}

// 获取当前可用数据源alarm
func (s *STATISD) GetDataSource() string {
	if dataSource != "" {
		return dataSource
	}
	// 读取是否配置了告警数据源
	dataSource, _ = internal.GetAlarmDataSource()
	return dataSource
}

// 发送推送
func (s *STATISD) SendEventToHttp(pushList []*proto.EventPush, data *map[string]interface{}) (err error) {
	if len(pushList) == 0 {
		return nil
	}
	for _, v := range pushList {
		body := map[string]interface{}{
			"datetime": time.Now().Format("2006-01-02 15:04:05"),
			"name":     v.Name,
			"type":     v.DataType,
			"result":   nil,
		}
		if v.DataType == 0 {
			body["result"] = (*data)["message"]
		} else {
			body["result"] = data
		}
		bodyByte, err := json.Marshal(body)
		if err != nil {
			return err
		}
		ret, err := internal.PostUrlJsonBody(v.Url, bodyByte)
		if err != nil {
			return err
		}
		if Debug == true {
			log.Println(string(ret))
		}
	}
	return
}
