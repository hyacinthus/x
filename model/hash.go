package model

import (
	"time"

	"github.com/mitchellh/hashstructure"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// Source 用作抓取数据临时存储，靠哈希值判断是否已存在记录，不可更新和删除
type Source struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" hash:"-" gorm:"type:varchar(20);primary_key"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" hash:"-"`
	// 校验和
	Hash uint64 `json:"-" hash:"-" gorm:"index"`
}

// BeforeCreate GORM hook
func (*Source) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}

// Hash FNV hash64 of v
func Hash(v interface{}) (uint64, error) {
	return hashstructure.Hash(v, nil)
}
