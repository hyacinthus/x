package xcc

import (
	"github.com/hyacinthus/x/xtype"
)

// SGet 获取集合
func (c *client) SGet(key string) (xtype.Strings, error) {
	set, err := c.kv.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	return xtype.Strings(set), nil
}

// SAdd 为集合增加一个元素
func (c *client) SAdd(key, item string) error {
	return c.kv.SAdd(key, item).Err()
}

// LRemove 删除集合中的指定元素
func (c *client) SRemove(key, item string) error {
	return c.kv.SRem(key, item).Err()
}
