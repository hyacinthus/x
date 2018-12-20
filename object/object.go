package object

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// Config 腾讯云 COS 配置
type Config struct {
	AppID     string
	Region    string
	Bucket    string
	SecretID  string
	SecretKey string
}

// Client 封装的 client
type Client struct {
	config *Config
	client *cos.Client
}

// New 新建 cos 客户端
func New(config *Config) *Client {
	u, _ := url.Parse(fmt.Sprintf("http://%s-%s.cos.%s.myqcloud.com",
		config.Bucket, config.AppID, config.Region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	return &Client{
		config: config,
		client: c,
	}
}

// Get 获取 cos 对象
func (c *Client) Get(key string) ([]byte, error) {
	resp, err := c.client.Object.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bs, nil
}

// Put 写文件
func (c *Client) Put(key string, f io.Reader) error {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/markdown",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			//XCosACL: "public-read",
			XCosACL: "private",
		},
	}
	_, err := c.client.Object.Put(context.Background(), key, f, opt)
	if err != nil {
		return err
	}
	return nil
}
