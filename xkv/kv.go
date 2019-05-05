// Package xkv 获得一个 redis 的 client，
// 虽然 redis 客户端会自动重连，但是首次若连接不成功会报错，这里会持续尝试连接。
package xkv

import (
	"github.com/go-redis/redis"
	"github.com/hyacinthus/x/xlog"
)

var log = xlog.Get()

// Config 数据库配置，可以被主配置直接引用
type Config struct {
	Host     string `default:"redis"`
	Port     string `default:"6379"`
	Password string
	DB       int `default:"0"`
}

// New 用配置生成一个 redis 数据库 client,若目标数据库未启动会一直等待
func New(config Config) *redis.Client {
	var kv = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	// TODO: 这里要ping到通为止

	log.Info("Redis connect successful.")

	return kv
}
