// Package xcache 可以实例化一个缓存客户端，和 cc 的区别是可以在多个 package 之间传递
package xcache

import (
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/hyacinthus/x/xlog"
	"github.com/vmihailenco/msgpack/v4"
)

var log = xlog.Get()

// Client 缓存客户端
type Client struct {
	kv    *redis.Client
	codec *cache.Codec
}

// NewClient 初始化客户端
func NewClient(kv *redis.Client) *Client {
	return &Client{
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
func (c *Client) Set(key string, object interface{}, exp time.Duration) {
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
func (c *Client) Get(key string, pointer interface{}) error {
	return c.codec.Get(key, pointer)
}

// Exists 是否存在
func (c *Client) Exists(key string) bool {
	return c.kv.Exists(key).Val() != 0
}

// Delete 清缓存
func (c *Client) Delete(key string) {
	err := c.codec.Delete(key)
	if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache faild")
	}
}

// Clean 批量清除一类缓存
func (c *Client) Clean(cate string) {
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
