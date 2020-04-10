// Package xobj 是为了创建一个统一的对象存储接口，以便切换存储服务之后可以不改业务代码
// 只需要修改配置内容和创建部分变量即可
package xobj

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/gommon/log"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/hyacinthus/x/xerr"
)

// cosClient 封装的 cos client
type cosClient struct {
	config Config
	client *cos.Client
}

// New 新建 cos 客户端
func newCosClient(bucket string, config Config) Client {
	u, _ := url.Parse(fmt.Sprintf("http://%s-%s.cos.%s.myqcloud.com",
		bucket, config.AppID, config.Region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	return &cosClient{
		config: config,
		client: c,
	}
}

func (c *cosClient) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, xerr.ErrNotFound
	}
	reader, err := c.Reader(key)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	return file, nil
}

func (c *cosClient) Reader(key string) (io.ReadCloser, error) {
	if key == "" {
		return nil, xerr.ErrNotFound
	}
	resp, err := c.client.Object.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// 上传只是程序内部读取的文件
func (c *cosClient) Put(key string, f io.Reader) error {
	return c.PutFile(key, "", "", f, 0)
}

// 上传对下载友好的文件
func (c *cosClient) PutFile(key, name, contentType string, f io.Reader, contentLength int) error {
	if key == "" {
		return xerr.New(400, "EmptyKey", "empty key")
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, name),
			ContentType:        contentType,
			ContentLength:      contentLength,
		},
	}
	_, err := c.client.Object.Put(context.Background(), key, f, opt)
	if err != nil {
		return err
	}
	return nil
}

func (c *cosClient) Delete(key string) error {
	if key == "" {
		return xerr.New(400, "EmptyKey", "empty key")
	}
	_, err := c.client.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}

// Exists 腾讯云的 SDK 暂时不支持将不存在的情况分拣出来，目前出错了也不返回错误。
func (c *cosClient) Exists(key string) (bool, error) {
	if key == "" {
		return false, xerr.New(400, "EmptyKey", "empty key")
	}
	resp, err := c.client.Object.Head(context.Background(), key, nil)
	if resp != nil && resp.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		log.Errorf("COS发生了预期外的错误，请注意。Key： %s， Error： %s", key, err)
		return false, err
	}
	return true, nil
}
