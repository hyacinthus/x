package xtype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/url"
)

// TODO: 将项目特有信息排除
var (
	imgprefix  = mustParse("https://image.xuebaox.com")
	fileprefix = mustParse("https://static.xuebaox.com")
)

// ImageURL ===============图片链接===============
type ImageURL string

// String 转换为string类型
func (f ImageURL) String() string {
	if f == "" {
		return ""
	}
	u, err := url.Parse(string(f))
	if err != nil {
		return ""
	}
	// 如果已经是网址了，直接返回
	if u.IsAbs() {
		return u.String()
	}
	// 如果只是路径，加上网址
	return imgprefix.ResolveReference(u).String()
}

// IsEmpty 是否为空
func (f ImageURL) IsEmpty() bool {
	s := string(f)
	return s == ""
}

// MarshalJSON 转换为json类型 加域名
func (f ImageURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON 前端提交的如果是URL，那只保留path
func (f *ImageURL) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	u, err := url.Parse(tmp)
	if err != nil {
		return err
	}
	if u.Hostname() == imgprefix.Hostname() {
		*f = ImageURL(u.EscapedPath())
	} else {
		*f = ImageURL(tmp)
	}
	return nil
}

// Scan implements the Scanner interface.
func (f *ImageURL) Scan(src interface{}) error {
	if src == nil {
		*f = ""
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read file url data from DB failed")
	}
	*f = ImageURL(tmp)
	return nil
}

// Value implements the driver Valuer interface.
func (f ImageURL) Value() (driver.Value, error) {
	return string(f), nil
}

// FileURL ================文件链接================
type FileURL string

// String 转换为string类型
func (f FileURL) String() string {
	if f == "" {
		return ""
	}
	u, err := url.Parse(string(f))
	if err != nil {
		return ""
	}
	// 如果已经是网址了，直接返回
	if u.IsAbs() {
		return u.String()
	}
	// 如果只是路径，加上网址
	return fileprefix.ResolveReference(u).String()
}

// IsEmpty 是否为空
func (f FileURL) IsEmpty() bool {
	s := string(f)
	return s == ""
}

// MarshalJSON 转换为json类型 加域名
func (f FileURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON 不做处理
func (f *FileURL) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	u, err := url.Parse(tmp)
	if err != nil {
		return err
	}
	if u.Hostname() == fileprefix.Hostname() {
		*f = FileURL(u.EscapedPath())
	} else {
		*f = FileURL(tmp)
	}
	return nil
}

// Scan implements the Scanner interface.
func (f *FileURL) Scan(src interface{}) error {
	if src == nil {
		*f = ""
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read file url data from DB failed")
	}
	*f = FileURL(tmp)
	return nil
}

// Value implements the driver Valuer interface.
func (f FileURL) Value() (driver.Value, error) {
	return string(f), nil
}

// =============== 内部方法 ================
func mustParse(src string) *url.URL {
	u, err := url.Parse(src)
	if err != nil {
		panic(err)
	}
	return u
}
