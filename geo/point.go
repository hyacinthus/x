package geo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Point 一个地图坐标点
type Point struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
}

// String 转换为string类型
func (p Point) String() string {
	resp, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(resp)
}

// Short 格式化为短字符串
func (p Point) Short() string {
	return fmt.Sprintf("%.6f,%.6f", p.Latitude, p.Longitude)
}

// ParsePoint 从字符串解析
// 两种格式均可 json格式和短字符串
func ParsePoint(src string) (*Point, error) {
	// json 格式
	if strings.HasPrefix(src, "{") {
		p := new(Point)
		err := json.Unmarshal([]byte(src), p)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	// 短字符串格式
	pair := strings.Split(src, ",")
	if len(pair) != 2 {
		return nil, fmt.Errorf("parse geo point data from DB failed:%s", src)
	}
	la, err := strconv.ParseFloat(pair[0], 64)
	if err != nil {
		return nil, fmt.Errorf("parse geo point data from DB failed:%s", src)
	}
	lo, err := strconv.ParseFloat(pair[1], 64)
	if err != nil {
		return nil, fmt.Errorf("parse geo point data from DB failed:%s", src)
	}
	return &Point{
		Latitude:  la,
		Longitude: lo,
	}, nil
}

// Scan implements the Scanner interface.
func (p *Point) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read geo point data from DB failed")
	}
	point, err := ParsePoint(string(tmp))
	if err != nil {
		return err
	}
	*p = *point
	return nil
}

// Value implements the driver Valuer interface.
// 存储格式 维度,经度
func (p Point) Value() (driver.Value, error) {
	return p.Short(), nil
}
