package model

import (
	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/x/xtype"
	"github.com/jinzhu/gorm"
)

// Sorter 排序器
type Sorter struct {
	Sort      string `query:"sort"`      // 排序字段的 gorm 名称
	Direction string `query:"direction"` // 方向 asc/desc
}

// Apply 将排序信息应用的数据库查询
// d 为缺省值 例如 "created_at desc" 不会被此函数检查 写错了会造成数据库查询出错
// limits 是 sort 字段的范围，防止攻击
func (s Sorter) Apply(tx *gorm.DB, d string, limits ...string) (*gorm.DB, error) {
	var sum string
	// 获得排序字段
	if s.Sort == "" {
		if d == "" {
			// 查询和默认值都为空 不排序
			return tx, nil
		}
		sum = d
	} else {
		if !xtype.Strings(limits).Contains(s.Sort) {
			return nil, xerr.New(400, "InvalidSort", "请求包含非法的排序字段")
		}
		sum = s.Sort
		// 只有存在 sort 参数，才检查方向，只有 desc 时有效
		if s.Direction == "desc" {
			sum += " desc"
		}
	}
	// 返回加过排序的数据库查询
	return tx.Order(sum), nil
}
