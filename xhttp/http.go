package xhttp

import (
	"net/http"
	"time"
)

// Client 自定义 http 客户端
var Client *http.Client

// 全局只会执行一次
func init() {
	Client = &http.Client{
		Timeout: time.Second * 5,
	}
}

// SetTimeout 设置全局 http 客户端的超时时间
func SetTimeout(seconds int) {
	Client.Timeout = time.Second * time.Duration(seconds)
}
