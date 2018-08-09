package apis

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/redis"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	rpc "github.com/gorilla/rpc/v2"
	json "github.com/gorilla/rpc/v2/json2"
)

// 响应数据--请求数据都使用map[string]interface{}
type Response interface{}

type Apis struct {
	muxRouter *mux.Router
}

func NewApis(mux *mux.Router) *Apis {
	return &Apis{muxRouter: mux}
}

func (this *Apis) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	bodyRead, err := this.readAndAssignResponseBody(req)
	if err != nil {
		log.Println("请求验证处理错误" + err.Error())
		internal.LogFile.W("请求验证处理错误:" + err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(bodyRead)
	postBody, err := simplejson.NewJson(bytes)
	if err != nil {
		log.Println("请求体json处理错误" + err.Error())
		internal.LogFile.W("请求体json处理错误:" + err.Error())
		return
	}
	// 获取信息
	method, _ := postBody.Get("method").String()
	methodInfo := strings.Split(method, ".")
	if len(methodInfo) == 2 && methodInfo[0] != "Public" {
		token, _ := postBody.Get("id").String()
		// 验证登录
		u, _, err := redis.GetSessionAdmin(token)
		if err != nil || u == nil {
			wBody := `{"jsonrpc": "2.0", "error": {"code": -32000, "message": "not login"}, "id": "` + token + `"}`
			w.Write([]byte(wBody))
			return
		}
	}
	// 验证接口权限
	log.Println(method)
	// 在这里验证用户登录信息，和刷新redis中的session--登录和退出不需要验证权限
	// 设置头信息
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("X-Powered-By", "PHP/7.0.0")
	w.Header().Add("Server", "nginx/1.6.2")
	this.muxRouter.ServeHTTP(w, req)
}

// 读取body内容，并且重新赋值到Request
func (this *Apis) readAndAssignResponseBody(req *http.Request) (io.Reader, error) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return bytes.NewReader(buf), nil
}

/* 格式转换 */
// 转string
func (req *Apis) ToString(val interface{}, def ...string) (str string) {
	if val != nil {
		str = fmt.Sprint(val)
	} else {
		if len(def) > 0 {
			str = def[0]
		} else {
			str = ""
		}
		str = ""
	}
	return str
}

// 转int
func (req *Apis) ToInt(val interface{}, def ...int) (rval int, err error) {
	if val != nil {
		str := fmt.Sprint(val)
		rval, err = strconv.Atoi(str)
		if err != nil {
			if str == "true" {
				rval = 1
			} else {
				rval = 0
			}
		}
	} else {
		err = errors.New("格式转换错误 int")
		if len(def) > 0 {
			rval = def[0]
		} else {
			rval = 0
		}
	}
	return
}

// 转int32
func (req *Apis) ToInt32(val interface{}, def ...int32) (rval int32, err error) {
	if val != nil {
		str := fmt.Sprint(val)
		rval1, err := strconv.Atoi(str)
		if err != nil {
			if str == "true" {
				rval = 1
			} else {
				rval = 0
			}
		}
		rval = int32(rval1)
	} else {
		err = errors.New("格式转换错误 int")
		if len(def) > 0 {
			rval = def[0]
		} else {
			rval = 0
		}
	}
	return
}

// 转bool
func (req *Apis) ToBool(val interface{}, def ...bool) (rval bool, err error) {
	if val != nil {
		str := fmt.Sprint(val)
		if str == "true" {
			rval = true
		} else {
			rval = false
		}
	} else {
		err = errors.New("格式转换错误 bool")
		if len(def) > 0 {
			rval = def[0]
		} else {
			rval = false
		}
	}
	return
}

// 验证token--权限严重，每个方法添加此方法调用
func (req *Apis) ChkToken() error {
	// 在redis中查询登录信息
	return errors.New("not login")
}

// 获取用户信息
func (req *Apis) getUserInfo(token string) models.User {
	return models.User{Id: 1}
}

// 更新redis中登录用户信息
func (req *Apis) setSession(token string, u *models.User) error {
	return nil
}

var RpcServer *rpc.Server

func Run(route ...string) {
	// 服务监听
	rpcRoute := ""
	if len(route) > 0 {
		rpcRoute = route[0]
	} else {
		rpcRoute = "/v1/jsonrpc"
	}
	log.Println(fmt.Sprintf("rpc:服务监听地址[%s]", rpcRoute))
	internal.LogFile.I(fmt.Sprintf("rpc:服务监听地址[%s]", rpcRoute))
	http.Handle(rpcRoute, NewApis(initRpcServer(rpcRoute)))
}

func initRpcServer(route string) *mux.Router {
	// 初始化rpc对象
	RpcServer = rpc.NewServer()
	RpcServer.RegisterCodec(json.NewCodec(), "application/json")
	// 远程方法
	RpcServer.RegisterService(new(AppName), "")
	RpcServer.RegisterService(new(Public), "")
	RpcServer.RegisterService(new(Device), "")
	RpcServer.RegisterService(new(DeviceGroup), "")
	RpcServer.RegisterService(new(DeviceGroupGroup), "")
	RpcServer.RegisterService(new(Search), "")
	RpcServer.RegisterService(new(EventLevel), "")
	RpcServer.RegisterService(new(EventPush), "")
	RpcServer.RegisterService(new(EventRule), "")
	RpcServer.RegisterService(new(EventDevice), "")
	RpcServer.RegisterService(new(Setings), "")
	RpcServer.RegisterService(new(SetingsTemplate), "")
	RpcServer.RegisterService(new(Topics), "")
	RpcServer.RegisterService(new(TopicsRule), "")
	RpcServer.RegisterService(new(User), "")
	RpcServer.RegisterService(new(System), "")

	// 设置路由
	rpcRouter := mux.NewRouter()
	rpcRouter.Handle(route, RpcServer)
	return rpcRouter
}
