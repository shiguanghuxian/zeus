package apis

import (
	"errors"
	"net/http"
	"strconv"

	"53it.net/zues/internal"
	"53it.net/zues/models"
)

type System struct {
	Apis
}

// 告警列表
func (this *System) SystemUnitConversionList(r *http.Request, args *map[string]interface{}, response *Response) error {
	var totalRows int64 // 总行数
	var err error
	// 关键词
	keyword := this.ToString((*args)["keyword"])
	// 总行数
	totalRows, err = models.GetSystemUnitConversionCount(keyword)
	if err != nil {
		return errors.New("服务端错误 count")
	}
	// 每页行数
	listRows, err := internal.CFG.Int("apis", "pagecount")
	if err != nil {
		listRows = 10
	}
	// 当前页码
	page, err := this.ToInt((*args)["page"], 1)
	if err != nil {
		page = 1
	}
	// 查询列表
	list, err := models.GetSystemUnitConversionList(page, listRows, keyword)
	if err != nil {
		return errors.New("服务端错误 list")
	}
	// 页面数据
	*response = map[string]interface{}{
		"page":       page,
		"total_rows": totalRows,
		"list_rows":  listRows,
		"list":       list,
	}
	return nil
}

// 添加
func (this *System) AddSystemUnitConversion(r *http.Request, args *map[string]interface{}, response *Response) error {
	// 接收参数
	originalUnit := this.ToString((*args)["original_unit"])
	afterUnit := this.ToString((*args)["after_unit"])
	if originalUnit == "" || afterUnit == "" {
		return errors.New("原单位和转换后的单位不能为空")
	}
	// 检查是否存在单位转换
	if ok := models.ChkSystemUnitConversion(originalUnit, afterUnit); ok == true {
		return errors.New("原单位转换已经存在（每组单位转换只能存在一个）")
	}
	multiple, err := strconv.ParseFloat(this.ToString((*args)["multiple"]), 64)
	if err != nil || multiple == 0 {
		return errors.New("倍率输入错误")
	}
	luaCode := this.ToString((*args)["lua_code"])
	Ctype, _ := this.ToInt((*args)["type"], 0)
	// models对象
	systemUnitConversion := new(models.SystemUnitConversion)
	systemUnitConversion.OriginalUnit = originalUnit
	systemUnitConversion.AfterUnit = afterUnit
	systemUnitConversion.Multiple = multiple
	systemUnitConversion.LuaCode = luaCode
	systemUnitConversion.Type = Ctype
	// 执行插入
	_, err = models.AddOneSystemUnitConversion(systemUnitConversion)
	if err != nil {
		return errors.New("单位转换添加错误")
	}
	return nil
}

// 编辑
func (this *System) EditSystemUnitConversion(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数id错误")
	}
	// 接收参数
	originalUnit := this.ToString((*args)["original_unit"])
	afterUnit := this.ToString((*args)["after_unit"])
	multiple, err := strconv.ParseFloat(this.ToString((*args)["multiple"]), 64)
	if err != nil || multiple == 0 {
		return errors.New("倍率输入错误")
	}
	luaCode := this.ToString((*args)["lua_code"])
	Ctype, _ := this.ToInt((*args)["type"], 0)
	if originalUnit == "" || afterUnit == "" {
		return errors.New("原单位和转换后的单位不能为空")
	}
	// models对象
	systemUnitConversion := new(models.SystemUnitConversion)
	systemUnitConversion.OriginalUnit = originalUnit
	systemUnitConversion.AfterUnit = afterUnit
	systemUnitConversion.Multiple = multiple
	systemUnitConversion.LuaCode = luaCode
	systemUnitConversion.Type = Ctype
	// 执行插入
	_, err = models.EditOneSystemUnitConversion(id, systemUnitConversion)
	if err != nil {
		return errors.New("单位转换添加错误")
	}
	return nil
}

// 删除
func (this *System) DelSystemUnitConversion(r *http.Request, args *map[string]interface{}, response *Response) error {
	id, _ := this.ToInt((*args)["id"], 0)
	if id == 0 {
		return errors.New("参数错误")
	}
	_, err := models.DelOneSystemUnitConversion(id)
	if err != nil {
		return errors.New("删除错误")
	}
	return nil
}
