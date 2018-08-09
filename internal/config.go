package internal

import (
	"io/ioutil"
	"os"

	"github.com/robfig/config"
)

var configFile string = "./conf/cfg.ini"
var CFG *config.Config

func init() {
	configFile = GetRootDir() + "/conf/cfg.ini"
	if _, err := os.Stat(configFile); err != nil {
		err1 := ioutil.WriteFile(configFile, []byte(""), 0666)
		if err1 != nil {
			LogFile.E("防止崩溃写文件cfg.ini错误:" + err1.Error())
		}
	}
	initConfig()
}

// 初始化配置文件
func initConfig() error {
	// 设置配置文件
	cfg, err := config.ReadDefault(configFile)
	if err != nil {
		LogFile.E("配置文件读取失败:" + err.Error())
		return err
	}
	CFG = cfg
	return nil
}

// 重置config
func ResetConfig() error {
	return initConfig()
}

// 设置配置文件路径
func SetConfigName(name string) {
	configFile = GetRootDir() + "/config/" + name + ".ini"
	initConfig()
}
