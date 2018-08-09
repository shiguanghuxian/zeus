package internal

import (
	"errors"
	"fmt"
	"strconv"

	"53it.net/zues/proto"

	lua "github.com/yuin/gopher-lua"
)

// 系统单位值转换
func ConvertUnit(myUnit, value string, conversionInfo *proto.SystemUnitConversion) (val string, err error) {
	value1, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", errors.New("value数值转换错误")
	}
	if conversionInfo.OriginalUnit == myUnit {
		if conversionInfo.Type == 0 {
			val = fmt.Sprintf("%.2f", value1/conversionInfo.Multiple)
		} else {
			return ConvertLuaCode(conversionInfo.LuaCode, value, 0)
		}
	} else if conversionInfo.AfterUnit == myUnit {
		if conversionInfo.Type == 0 {
			val = fmt.Sprintf("%.2f", value1*conversionInfo.Multiple)
		} else {
			return ConvertLuaCode(conversionInfo.LuaCode, value, 1)
		}
	}
	return
}

/*
		--[[
            original: 要转换的原单位数
            aspect: 转换方向标识，0从左到右，1从右到左
        ]]--
        if (aspect == 0)
        then
            return original / 1024
        else
            return original * 1024
        end
*/
// 系统单位转换，使用到了lua代码
func ConvertLuaCode(luaCode, original string, aspect int) (string, error) {
	if luaCode == "" {
		return "", errors.New("lua代码不能为空")
	}
	if original == "" {
		return "", errors.New("原单位数不能为空")
	}
	ch := make(chan lua.LValue, 1)
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("ch", lua.LChannel(ch))
	if err := L.DoString(fmt.Sprintf(`
    function convert(original, aspect)
		%s
    end
    ch:send(convert(%s, %d))
  `, luaCode, original, aspect)); err != nil {
		return "", err
	}
	var outVal lua.LValue
	select {
	case outVal = <-ch:
		return outVal.String(), nil
	}
}
