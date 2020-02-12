// Package xcc 缓存客户端
package xcc

import (
	"time"

	"github.com/hyacinthus/x/xtype"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/hyacinthus/x/xlog"
	"github.com/vmihailenco/msgpack/v4"
)

var log = xlog.Get()

// Client is a cache client interface
type Client interface {
	// Set 写缓存
	Set(key string, object interface{}, exp time.Duration)
	// MustSet 写缓存,并检查错误
	MustSet(key string, object interface{}, exp time.Duration) error
	// Get 读缓存
	Get(key string, pointer interface{}) error
	// SetString 写 string 缓存
	SetString(key string, s string, exp time.Duration) error
	// GetString 读 string 缓存
	GetString(key string) (string, error)
	// Exists 是否存在
	Exists(key string) bool
	// Expire 刷新过期时间
	Expire(key string, ex time.Duration) error
	// Delete 清缓存
	Delete(key string)
	// Clean 批量清除一类缓存
	Clean(cate string)
	// LGet 获取列表
	LGet(key string) (xtype.Strings, error)
	// LPush 为列表右侧增加一个元素
	LPush(key, item string) error
	// LRemove 删除列表中的指定元素
	LRemove(key, item string) error
	// SGet 获取集合
	SGet(key string) (xtype.Strings, error)
	// SAdd 为集合增加一个元素
	SAdd(key, item string) error
	// SAdd 为集合增加一个元素，并刷新过期时间
	SAddEx(key, item string, ex time.Duration) error
	// SRemove 删除集合中的指定元素
	SRemove(key, item string) error
}

// client 缓存客户端
type client struct {
	kv    *redis.Client
	codec *cache.Codec
}

// New 初始化客户端
func New(kv *redis.Client) Client {
	return &client{
		kv: kv,
		codec: &cache.Codec{
			Redis: kv,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
}

// Set 写缓存
func (c *client) Set(key string, object interface{}, exp time.Duration) {
	err := c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
	if err != nil {
		log.WithError(err).WithField("key", key).Error("set cache failed")
	}
}

// MustSet 写缓存,检查并返回错误
func (c *client) MustSet(key string, object interface{}, exp time.Duration) error {
	return c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
}

// Get 读缓存
func (c *client) Get(key string, pointer interface{}) error {
	return c.codec.Get(key, pointer)
}

// Set 写 string 缓存
func (c *client) SetString(key string, s string, exp time.Duration) error {
	return c.kv.Set(key, s, exp).Err()
}

// Get 读 string 缓存
func (c *client) GetString(key string) (string, error) {
	return c.kv.Get(key).Result()
}

// Exists 是否存在
func (c *client) Exists(key string) bool {
	return c.kv.Exists(key).Val() != 0
}

// Delete 清缓存
func (c *client) Delete(key string) {
	err := c.codec.Delete(key)
	if err == cache.ErrCacheMiss {
		return
	} else if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache failed")
	}
}

// Expire 刷新过期时间
func (c *client) Expire(key string, ex time.Duration) error {
	return c.kv.Expire(key, ex).Err()
}

// Clean 批量清除一类缓存
func (c *client) Clean(cate string) {
	if cate == "" {
		log.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range c.kv.Keys(cate + "*").Val() {
		err := c.codec.Delete(key)
		if err != nil {
			log.WithError(err).WithField("key", key).Error("delete cache failed,stop batch delete")
			break
		}
		i++
	}
	log.Infof("delete %d %s cache", i, cate)
}
