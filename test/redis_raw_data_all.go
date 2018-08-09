package main

import (
	"encoding/json"
	"fmt"

	"53it.net/zues/internal"
	"53it.net/zues/redis"
)

func main() {
	// 读取mongodb配置
	address, _ := internal.CFG.String("redis", "address")
	port, _ := internal.CFG.String("redis", "port")

	fmt.Println(address, port)

	fmt.Println("查询设备全部指标")
	testAll()
	fmt.Println("查询单条指标")
	testOne()
	// select {}
}

func testOne() {
	rawData := make(map[string]string)
	rawData["group"] = "default"
	rawData["hostname"] = "http-server2"
	rawData["ip"] = "101.200.174.134"

	one, err := redis.GetOneNewestData(rawData, "zn_raw_data", "value")
	fmt.Println("错误：", err)
	str, _ := json.Marshal(one)
	fmt.Println(string(str))
}

func testAll() {
	rawData := make(map[string]string)
	rawData["group"] = "default"
	rawData["hostname"] = "http-server2"
	rawData["ip"] = "101.200.174.134"

	list, err := redis.GetAllNewestData(rawData, "zn_raw_data")

	fmt.Println("错误：", err)
	str, _ := json.Marshal(list)
	fmt.Println(string(str))
}
