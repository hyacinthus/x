// Package xobj 是为了创建一个统一的对象存储接口，以便切换存储服务之后可以不改业务代码
// 只需要修改配置内容和创建部分变量即可
package xobj

import (
	"io"
)

// Provider 云服务提供商
type Provider int

// all providers
const (
	ProviderCOS Provider = iota
)

// Config 腾讯云 COS 配置
type Config struct {
	AppID     string
	Region    string
	SecretID  string
	SecretKey string
}

// Client xobj client
type Client interface {
	// 获取一个文件的 ReadCloser ，记得关闭它
	Reader(key string) (io.ReadCloser, error)
	// 获取文件
	Get(key string) ([]byte, error)
	// 存储文件，不会进行重复检查，注意 reader 必须会自己 EOF
	Put(key string, f io.Reader) error
	// 存储带类型的文件
	PutFile(key, name, contentType string, f io.Reader, contentLength int) error
	// 删除文件
	Delete(key string) error
}

// New 新建存储客户端，为了混用不同的基础施舍，供应商和bucket在调用时填写，不放在设置中。
func New(provider Provider, bucket string, config Config) Client {
	switch provider {
	case ProviderCOS:
		return newCosClient(bucket, config)
	default:
		panic("invalid provider")
	}
}
