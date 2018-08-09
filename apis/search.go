package apis

import (
	"errors"
	"net/http"
	"time"

	"53it.net/zues/elasticsearch"
	"53it.net/zues/influxdb"
	"53it.net/zues/internal"
	"53it.net/zues/models"
	"53it.net/zues/mongo"
)

type Search struct {
	Apis
}

// 查询api
func (this *Search) ZqlQueryV1(r *http.Request, args *map[string]interface{}, response *Response) error {
	zqlStr := this.ToString((*args)["zql"])
	if zqlStr == "" {
		return errors.New("查询语句不能为空")
	}
	// 获取数据源类型
	dataType := this.ToString((*args)["data_source"])
	// 区分数据源
	var list interface{}
	var err error
	if dataType == "" {
		dataType = "influxdb"
	}
	if dataType == "influxdb" {
		list, err = influxdb.ZqlQueryCmd(zqlStr)
	} else if dataType == "mongodb" {
		list, err = mongo.GetZqlList(zqlStr)
	} else if dataType == "elastic" {
		list, err = elasticsearch.GetZqlList(zqlStr)
	}
	if err != nil {
		return err
	}

	*response = list
	return nil
}

// 当前可用数据源
func (this *Search) DataSource(r *http.Request, args *map[string]interface{}, response *Response) error {
	list := make([]map[string]string, 0)
	// 读配置组织数据
	ismgo, err := internal.CFG.Bool("mongodb", "enable")
	if err != nil {
		return errors.New("读区配置文件错误")
	}
	if ismgo == true {
		list = append(list, map[string]string{
			"type": "mongodb",
			"name": "mongodb",
		})
	}
	isinf, _ := internal.CFG.Bool("influxdb", "enable")
	if isinf == true {
		list = append(list, map[string]string{
			"type": "influxdb",
			"name": "influxdb",
		})
	}
	isela, _ := internal.CFG.Bool("elasticsearch", "enable")
	if isela == true {
		list = append(list, map[string]string{
			"type": "elastic",
			"name": "elastic",
		})
	}

	*response = list
	return nil
}

// 当前可用数据源
func (this *Search) Autocomplete(r *http.Request, args *map[string]interface{}, response *Response) error {
	name := this.ToString((*args)["name"])
	list, err := models.GetAllAutocompleteList(name)
	if err != nil {
		return err
	}

	*response = list
	return nil
}

// 添加搜索zql
func (this *Search) AddAutocomplete(r *http.Request, args *map[string]interface{}, response *Response) error {
	name := this.ToString((*args)["name"])
	desc := this.ToString((*args)["desc"])
	if name == "" || len(desc) < 6 {
		return errors.New("参数错误")
	}
	autocomplete := new(models.Autocomplete)
	autocomplete.Name = name
	autocomplete.Desc = desc
	autocomplete.Addtime = time.Now().Unix()
	autocomplete.Count = 0
	_, err := models.AddOneAutocomplete(autocomplete)
	if err != nil {
		return err
	}

	return nil
}

// 搜索次数自增
func (this *Search) IncrementAutocomplete(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"])
	count, _ := this.ToInt((*args)["count"])
	if id == 0 {
		return errors.New("参数错误")
	}
	autocomplete := new(models.Autocomplete)
	autocomplete.Count = count
	_, err := models.UpdateIdAutocomplete(id, autocomplete)
	if err != nil {
		return err
	}

	return nil
}
