package internal

import (
	"github.com/aiwuTech/fileLogger"
)

var (
	LogFile *fileLogger.FileLogger
)

// 创建日志对象
func NewLog(name string) {
	LogFile = fileLogger.NewDefaultLogger(GetRootDir()+"/logs", name+".log")
	LogFile.SetLogLevel(fileLogger.INFO)
}

func init() {
	NewLog("error")
}
