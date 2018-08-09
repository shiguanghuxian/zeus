package internal

import (
	"time"

	"github.com/hanguofeng/gocaptcha"
)

var MyCaptcha *gocaptcha.Captcha

func GetCaptcha() (*gocaptcha.Captcha, error) {
	if MyCaptcha == nil {
		wordDict, captchaConfig, imageConfig, filterConfig, storeConfig := loadConfig()
		wordmgr, err := gocaptcha.CreateWordManagerFromDataFile(wordDict)
		MyCaptcha, err = gocaptcha.CreateCaptcha(wordmgr, captchaConfig, imageConfig, filterConfig, storeConfig)
		return MyCaptcha, err
	}
	return MyCaptcha, nil
}

func loadConfig() (string, *gocaptcha.CaptchaConfig, *gocaptcha.ImageConfig, *gocaptcha.FilterConfig, *gocaptcha.StoreConfig) {
	// 配置文件路径-字符集
	wordDict := "./conf/gocaptcha/en_char"

	captchaConfig := new(gocaptcha.CaptchaConfig)
	captchaConfig.LifeTime = 120 * time.Second

	imageConfig := new(gocaptcha.ImageConfig)
	imageConfig.FontFiles = []string{"./conf/gocaptcha/zpix.ttf"} // 字体
	imageConfig.FontSize = 28
	imageConfig.Height = 40
	imageConfig.Width = 120

	filterConfig := new(gocaptcha.FilterConfig)
	filterConfig.Init()
	filterConfig.Filters = []string{"ImageFilterNoiseLine", "ImageFilterNoisePoint", "ImageFilterStrike"}

	var filterConfigGroup *gocaptcha.FilterConfigGroup
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "5")
	filterConfig.SetGroup("ImageFilterNoiseLine", filterConfigGroup)
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "10")
	filterConfig.SetGroup("ImageFilterNoisePoint", filterConfigGroup)
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "1")
	filterConfig.SetGroup("ImageFilterStrike", filterConfigGroup)

	storeConfig := new(gocaptcha.StoreConfig)
	storeConfig.Engine = gocaptcha.STORE_ENGINE_BUILDIN
	storeConfig.GcDivisor = 100
	storeConfig.GcProbability = 1

	return wordDict, captchaConfig, imageConfig, filterConfig, storeConfig
}
