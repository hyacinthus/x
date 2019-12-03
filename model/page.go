package model

import (
	"strconv"

	"github.com/hyacinthus/x/xerr"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// Paginator 分页器
type Paginator struct {
	Page    string `query:"page"`
	PerPage string `query:"per_page"`
}

// Apply 将翻页信息应用的数据库查询，并在 echo Context 中添加相应 Header
func (p Paginator) Apply(tx *gorm.DB, c echo.Context) (*gorm.DB, error) {
	var page, pageSize int
	var err error
	// 获得页码
	if p.Page == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(p.Page)
		if err != nil {
			return nil, xerr.New(400, "InvalidPage", "请在URL中提供合法的页码")
		}
	}
	// 获得每页条数
	if p.PerPage == "" {
		// 这里没办法获取配置信息，那么就不让在项目级配置默认条数了，就20吧。
		// 如果某个 API 需要修改默认尺寸，在调用这个方法之前修改 Page 变量就好。
		pageSize = 20
	} else {
		pageSize, err = strconv.Atoi(p.PerPage)
		if err != nil {
			return nil, xerr.New(400, "InvalidPage", "请在URL中提供合法的每页条数")
		}
	}
	// 设置返回的Header
	c.Response().Header().Set("X-Page-Num", strconv.Itoa(page))
	c.Response().Header().Set("X-Page-Size", strconv.Itoa(pageSize))
	// 返回加过页码的数据库查询
	return tx.Offset((page - 1) * pageSize).Limit(pageSize), nil
}
