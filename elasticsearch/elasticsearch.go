package elasticsearch

import (
	"log"
	"os"
	"strings"
	"time"

	"53it.net/elastic_pool"
	"53it.net/zues/internal"
)

// var eClient *elastic.Client
var elasticPool *elastic_pool.ElasticPool

//var mutex sync.Mutex

func init() {
	newElastic()
}

// 创建连接
func newElastic() {
	defer func() {
		if e := recover(); e != nil {
			internal.LogFile.E("elastic致命错误", e)
			log.Println("elastic致命错误", e)
		}
	}()
	if internal.CFG == nil {
		return
	}
	// 是否启用
	enable, _ := internal.CFG.Bool("elasticsearch", "enable")
	if !enable {
		return
	}
	// 获取配置
	address, _ := internal.CFG.String("elasticsearch", "address")
	port, _ := internal.CFG.String("elasticsearch", "port")
	username, _ := internal.CFG.String("elasticsearch", "username")
	password, _ := internal.CFG.String("elasticsearch", "password")
	maxIdle, _ := internal.CFG.Int("elasticsearch", "max_idle")
	// maxPoolSize, _ := internal.CFG.Int("elasticsearch", "max_pool_size")
	if maxIdle == 0 {
		maxIdle = 64
	}
	// 创建一个客户端
	var err error
	elasticPool, err = elastic_pool.NewElasticPool(elastic_pool.SetUrl(address, port),
		elastic_pool.SetAuth(username, password),
		elastic_pool.SetShowLog(true),
		elastic_pool.SetMaxPoolIdle(maxIdle))
	if err != nil {
		internal.LogFile.E("elasticsearch创建连接池失败:" + err.Error())
		log.Println("elasticsearch创建连接池失败:" + err.Error())
		os.Exit(2)
	}
	log.Println("elasticsearch创建连接池成功")
	// 创建索引
	err = newIndex(getIndexName())
	if err != nil {
		internal.LogFile.E("elasticsearch创建索引失败:" + err.Error())
		log.Println("elasticsearch创建索引失败:" + err.Error())
	}
}

// 获取连接
func getClient() (*elastic_pool.IdleConn, error) {
	return elasticPool.Get()
}

// 归还连接
func releaseClient(idleConn *elastic_pool.IdleConn) {
	// idleConn.C.CloseIndex(getIndexName()).Do()
	elasticPool.Release(idleConn)
}

// 创建索引
func newIndex(idx string) (err error) {
	idleConn, err := getClient()
	defer releaseClient(idleConn)
	if err != nil {
		return err
	}
	// 验证是否存在索引
	if ok, err := idleConn.C.IndexExists(idx).Do(); ok && err == nil {
		return nil
	}
	// 请求体
	body := `{
	   "settings" : {
	      "number_of_shards" : ##,
	      "number_of_replicas" : &&
	   }
	}`
	// 读取索引配置
	shards, _ := internal.CFG.String("elasticsearch", "shards")
	replicas, _ := internal.CFG.String("elasticsearch", "replicas")
	if shards == "" {
		shards = "10"
	}
	if replicas == "" {
		replicas = "1"
	}
	// 替换字符串
	body = strings.Replace(body, "##", shards, -1)
	body = strings.Replace(body, "&&", replicas, -1)
	// 创建一个索引
	_, err = idleConn.C.CreateIndex(idx).BodyString(body).Do()
	if err != nil {
		log.Println("创建索引失败:" + err.Error())
		return err
	}
	return nil
}

// 获取索引名
func getIndexName() string {
	dbname, err := internal.CFG.String("elasticsearch", "dbname")
	subdbname, _ := internal.CFG.String("elasticsearch", "subdbname")
	if err != nil {
		dbname = "zues"
	}
	if subdbname == "" {
		subdbname = "2006"
	}
	subDname := time.Now().Format(subdbname)
	return dbname + "_" + subDname
}
