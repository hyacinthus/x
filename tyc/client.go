package tyc

import (
	"net/http"

	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/x/xhttp"
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
	httpc  *http.Client
	config Config
}

// NewClient create a 天眼查 client
// 注意： 调用这个方法 会将 xhttp 包的全局 http client 超时调整为这里指定的值
func NewClient(config Config) *Client {
	xhttp.SetTimeout(config.Timeout)
	return &Client{
		httpc:  xhttp.Client,
		config: config,
	}
}
