package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 发起post请求（body-json）
func PostUrlJsonBody(urlstr string, request []byte) (result []byte, err error) {
	// 请求体
	body := bytes.NewBuffer([]byte(request))
	// 发起请求
	res, err := http.Post(urlstr, "application/json;charset=utf-8", body)
	if err != nil {
		LogFile.E("发起http post请求失败:" + err.Error())
		return nil, err
	}
	// 读取请求结果
	result, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		LogFile.E("读取http post 请求返回数据失败:" + err.Error())
	}
	return
}

// 发起get请求获取数据(参数请自己带好)
func GetUrlBody(urlstr string) (result []byte, err error) {
	res, err := http.Get(urlstr)
	if err != nil {
		LogFile.E("发起http get 请求失败:" + err.Error())
		return nil, err
	}
	// 请求失败重新发起请求
	if res.StatusCode != 200 {
		LogFile.E(fmt.Sprintf("请求失败:%s", urlstr))
		res.Body.Close()
		err = errors.New("请求失败:" + urlstr)
	} else {
		result, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			LogFile.E("读取http get 请求返回数据失败:" + err.Error())
		}
	}
	return
}
