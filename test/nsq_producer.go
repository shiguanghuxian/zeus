package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
)

func main() {
	// 创建生产者对象
	producer, err := nsq.NewProducer("192.168.1.155:4150", nsq.NewConfig())
	defer producer.Stop()
	if err != nil {
		fmt.Println(err.Error())
	}

	str0 := `{
	    "ip":"101.200.174.134",
	    "group":"default",
	    "hostname":"http-server3",
	    "data":[
			{"kpiid":"1002","value":"#","date":"&&&&","instance":"cpu1"},
			{"kpiid":"1006","value":"#","date":"&&&&","instance":"cpu1","code":"this is code"}
	    ]
	}`
	//,"code":"304"
	str3 := `{
		    "ip":"101.200.174.134",
		    "group":"default",
		    "hostname":"http-server2",
		    "data":[
				"1003|#|2016-05-19 19:38:20|cpu|this is code1",
				"1007|#|2016-05-19 19:38:30|cpu|301"
		    ]
		}`

	str1 := `{
	    "ip":"101.200.174.134",
	    "group":"default",
	    "hostname":"http-server1",
	    "data":[
			"kpiid&1001&jjsd|value:#|date%2016-06-19 19:38:08|hehe*cpu@dsdsd",
			"kpiid&1005&jjsd|value:#|date%2016-07-17 11:34:36|hehe*cpu@dsdsd"
	    ]
	}`

	str1 = `{"ip": "192.168.1.2", "hostname": "shiguanghuxian", "group": "default", "device_type": "linux", "data": [{"shortterm": 0.51, "midterm": 0.27, "plugin": "load", "interval": 10, "longterm": 0.21, "datetime": 1486353704, "date": "2017-02-06 12:01:44", "plugin_instance": "", "type_instance": "", "type": "load"}]}`

	for {
		startDate := time.Now().Format("2006-01-02 15:04:05")

		for i := 0; i < 200; i++ {
			str6 := strings.Replace(str0, "#", strconv.Itoa(rand.Intn(1000000)), -1)
			str2 := strings.Replace(str1, "#", strconv.Itoa(rand.Intn(100000)), -1)
			str4 := strings.Replace(str3, "#", strconv.Itoa(rand.Intn(10000000)), -1)

			myDate := time.Now().Format("2006-01-02 15:04:05")
			log.Println(myDate)
			str7 := strings.Replace(str6, "&&&&", myDate, -1)

			log.Println(str6)

			err = producer.Publish("test0", []byte(str7))
			producer.Publish("test", []byte(str4))
			err = producer.Publish("test1", []byte(str2))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		endDate := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("时间：", startDate, "至", endDate) // 输出每次耗时

		time.Sleep(time.Second * 6)
	}

	// select {}
}
