// Package xcc 缓存客户端
package xcc

import (
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/hyacinthus/x/xlog"
	"github.com/vmihailenco/msgpack"
)

var log = xlog.Get()

// Client is a cache client interface
type Client interface {
	// Set 写缓存
	Set(key string, object interface{}, exp time.Duration)
	// Get 读缓存
	Get(key string, pointer interface{}) error
	// Exists 是否存在
	Exists(key string) bool
	// Delete 清缓存
	Delete(key string)
	// Clean 批量清除一类缓存
	Clean(cate string)
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
		log.WithError(err).WithField("key", key).Error("set cache faild")
	}
}

// Get 读缓存
func (c *client) Get(key string, pointer interface{}) error {
	return c.codec.Get(key, pointer)
}

// Exists 是否存在
func (c *client) Exists(key string) bool {
	return c.kv.Exists(key).Val() != 0
}

// Delete 清缓存
func (c *client) Delete(key string) {
	err := c.codec.Delete(key)
	if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache faild")
	}
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
			log.WithError(err).WithField("key", key).Error("delete cache faild,stop batch delete")
			break
		}
		i++
	}
	log.Infof("delete %d %s cache", i, cate)
}