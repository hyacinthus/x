package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// Entity 实体共用字段
type Entity struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" gorm:"type:varchar(20);primary_key"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 软删除
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// BeforeCreate GORM hook
func (*Entity) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", xid.New().String())
	return nil
}

// Log 日志共用字段
type Log struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" gorm:"type:varchar(20);primary_key"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate GORM hook
func (*Log) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", xid.New().String())
	return nil
}
