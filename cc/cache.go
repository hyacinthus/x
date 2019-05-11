// Package cc 将要废弃，使用 xcache 返回客户端的版本
package cc

import (
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/hyacinthus/x/xlog"
	"github.com/vmihailenco/msgpack"
)

var (
	codec *cache.Codec
	rdb   *redis.Client
	log   = xlog.Get()
)

// Init 用一个 redis 连接初始化缓存
func Init(client *redis.Client) {
	rdb = client
	codec = &cache.Codec{
		Redis: client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

// Set 写缓存
func Set(key string, object interface{}, exp time.Duration) {
	if codec == nil {
		log.Panic("init cache with cc.Init(rdb) first")
	}
	err := codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
	if err != nil {
		log.WithError(err).WithField("key", key).Error("set cache faild")
	}
}

// Get 读缓存
func Get(key string, pointer interface{}) error {
	if codec == nil {
		panic("init cache with cc.Init(rdb) first")
	}
	return codec.Get(key, pointer)
}

// Exists 是否存在
func Exists(key string) bool {
	return rdb.Exists(key).Val() != 0
}

// Delete 清缓存
func Delete(key string) {
	if codec == nil {
		panic("init cache with cc.Init(rdb) first")
	}
	err := codec.Delete(key)
	if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache faild")
	}
}

// Clean 批量清除一类缓存
func Clean(cate string) {
	if codec == nil {
		panic("init cache with cc.Init(rdb) first")
	}
	if cate == "" {
		log.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range rdb.Keys(cate + "*").Val() {
		err := codec.Delete(key)
		if err != nil {
			log.WithError(err).WithField("key", key).Error("delete cache faild,stop batch delete")
			break
		}
		i++
	}
	log.Infof("delete %d %s cache", i, cate)
}
