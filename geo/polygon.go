package geo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// Polygon 一个地图多边形 有序的顶点列表
type Polygon []Point

// String 转换为string类型
func (p Polygon) String() string {
	if p == nil {
		return ""
	}
	resp,err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(resp)
}

// Short 格式化为短字符串
func (p Polygon) Short() string {
	if p == nil {
		return ""
	}
	tmp := make([]string,len(p))
	for _,point:=range p {
		tmp = append(tmp,point.Short())
	}
	return strings.Join(tmp,"|")
}

// ParsePolygon 解析字符串为多边形
func ParsePolygon(src string) (Polygon,error) {
	var p = make(Polygon,0)
	// json 格式
	if strings.HasPrefix(src,"[") {
		err := json.Unmarshal([]byte(src),&p)
		if err != nil {
			return nil,err
		}
		return p,nil
	}
	// 短字符串格式
	list := strings.Split(src,"|")
	for _,item := range list {
		point,err := ParsePoint(item)
		if err != nil {
			return nil,err
		}
		p = append(p,*point)
	}
	return p,nil
}

// MarshalJSON 转换为json类型
func (p Polygon) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"),nil
	}
	return json.Marshal(p)
}

// UnmarshalJSON 不做处理
func (p *Polygon) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, p)
}

// Scan implements the Scanner interface.
func (p *Polygon) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	tmp,ok := src.(string)
	if !ok {
		return errors.New("read geo polygon data from DB failed")
	}
	polygon,err := ParsePolygon(tmp)
	if err != nil {
		return err
	}
	*p = polygon
	return nil
}

// Value implements the driver Valuer interface.
// 存储格式 维度,经度
func (p Polygon) Value() (driver.Value, error) {
	return p.Short(), nil
}
