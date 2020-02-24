package tyc

import (
	"net/url"
	"time"

	"github.com/hyacinthus/x/xerr"

	"github.com/levigross/grequests"
)

// FetchCompany 从天眼查获取数据
func (c *Client) FetchCompany(name string) (*Company, error) {
	// 检查输入
	if len(name) == 0 {
		return nil, ErrorEmptyName
	}
	// 请求
	params := make(url.Values)
	params.Set("name", name)
	data, err := grequests.Get("http://open.api.tianyancha.com/services/open/ic/baseinfo/2.0",
		&grequests.RequestOptions{
			Params:     map[string]string{"name": name},
			HTTPClient: c.httpc,
		})
	if err != nil {
		return nil, err
	}
	var resp = new(CompanyResp)
	err = data.JSON(resp)
	if err != nil {
		return nil, err
	}
	// 300000 errCode是未找到相关公司
	if resp.ErrorCode == 300000 {
		return nil, ErrorNotFound
	}
	if resp.ErrorCode != 0 {
		return nil, xerr.Newf(500, "TycError", "请求[%s]企业天眼查企业数据报错: %d %s", name, resp.ErrorCode, resp.Reason)
	}

	resp.Result.TycUpdatedAt = time.Now()
	return resp.Result, nil
}
