package tyc

import (
	"net/http"
	"time"

	"github.com/hyacinthus/x/xerr"

	"github.com/levigross/grequests"
)

// 错误
var (
	ErrorEmptyName = xerr.New(400, "EmptyName", "企业名称不能为空")
	ErrorNotFound  = xerr.New(400, "NotFound", "在天眼查没有找到这个企业名称")
)

// Config 配置
type Config struct {
	Token   string
	Timeout int `default:"5"`
}

// Client 维持一个持久化的 http client ，避免每次都重建
type Client struct {
	httpc *http.Client
}

// NewClient create a 天眼查 client
func NewClient(config Config) *Client {
	return &Client{httpc: grequests.BuildHTTPClient(grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": config.Token,
		},
		RequestTimeout: time.Second * time.Duration(config.Timeout),
	})}
}
