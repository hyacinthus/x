// Package xdb 获得一个 gorm 的 db client，
// 虽然 db 客户端会自动重连，但是首次若连接不成功会报错，这里会持续尝试连接。
package xdb

import (
	"time"

	"github.com/hyacinthus/x/xlog"
	"github.com/jinzhu/gorm"
)

var log = xlog.Get()

// Config 数据库配置，可以被主配置直接引用
type Config struct {
	Host     string `default:"mysql"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"root"`
	Name     string `default:"demo"`
	Lifetime int    `default:"3000"`
}

// New 用配置生成一个 gorm mysql 数据库对象,若目标数据库未启动会一直等待
func New(config Config) *gorm.DB {
	var db *gorm.DB
	var err error

	for {
		db, err = gorm.Open("mysql", config.User+":"+config.Password+
			"@tcp("+config.Host+":"+config.Port+")/"+config.Name+
			"?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s")
		if err != nil {
			log.WithError(err).Warn("waiting for connect to db")
			time.Sleep(time.Second * 2)
			continue
		}
		db.DB().SetConnMaxLifetime(time.Duration(config.Lifetime) * time.Second)
		log.Info("Mysql connect successful.")
		break
	}

	return db
}
