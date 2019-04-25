// Package xlog 会
// 返回一个固定的 logrus loggger
// 可以自己用 也可以传给别的应用
// 用来统一日志格式
package xlog

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Get 取得默认 logger 没有则创建一个
func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

// Debug 将logger设置为debug模式
func Debug() {
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "06-01-02 15:04:05",
	})
	logger.SetLevel(logrus.DebugLevel)
}
