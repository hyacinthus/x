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
