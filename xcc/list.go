package xcc

import (
	"github.com/hyacinthus/x/xtype"
)

// LGet 获取列表
func (c *client) LGet(key string) (xtype.Strings, error) {
	list, err := c.kv.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return xtype.Strings(list), nil
}

// LPush 为列表右侧增加一个元素
func (c *client) LPush(key, item string) error {
	return c.kv.RPush(key, item).Err()
}

// LRemove 删除列表中的指定元素
func (c *client) LRemove(key, item string) error {
	return c.kv.LRem(key, 0, item).Err()
}
