package xtype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"
)

// 新版的列表字段，存往 mysql 数据库的时候，将存储成 json 类型，用以支持 json 查询语法
// JStrings ===========字符串列表===========
type JStrings []string

// String 转换为string类型
func (t JStrings) String() string {
	if t == nil {
		return ""
	}
	tmp, _ := json.Marshal([]string(t))
	return string(tmp)
}

// Contains 是否包含元素
func (t JStrings) Contains(s string) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// Intersectant 是否有交集
func (t JStrings) Intersectant(s JStrings) bool {
	for _, ss := range s {
		if t.Contains(ss) {
			return true
		}
	}
	return false
}

// SAdd 添加不重复的元素
func (t JStrings) SAdd(s string) JStrings {
	if t.Contains(s) {
		return t
	}
	return append(t, s)
}

// Remove 移除所有值为 s 的元素
func (t JStrings) Remove(s string) JStrings {
	tmp := make(JStrings, 0, len(t))
	for _, tt := range t {
		if tt != s {
			tmp = append(tmp, tt)
		}
	}
	return tmp
}

// Union 返回并集，s若有重复不会被去重
func (t JStrings) Union(s JStrings) JStrings {
	for _, tt := range t {
		s = s.SAdd(tt)
	}
	return s
}

// Sub 返回集合 t 减掉 s 后的集合
func (t JStrings) Sub(s Strings) JStrings {
	tmp := t
	for _, ss := range s {
		tmp = tmp.Remove(ss)
	}
	return tmp
}

// MarshalJSON 转换为json类型
func (t JStrings) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(t))
}

// UnmarshalJSON 不做处理
func (t *JStrings) UnmarshalJSON(data []byte) error {
	var tmp []string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = tmp
	return nil
}

// MarshalYAML 转换为json类型
func (t JStrings) MarshalYAML() ([]byte, error) {
	if t == nil {
		return []byte{}, nil
	}
	return yaml.Marshal([]string(t))
}

// UnmarshalYAML 不做处理
func (t *JStrings) UnmarshalYAML(data []byte) error {
	var tmp []string
	if err := yaml.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = tmp
	return nil
}

// Scan implements the Scanner interface.
func (t *JStrings) Scan(src interface{}) error {
	*t = make([]string, 0)
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read json string array from DB failed")
	}
	if len(tmp) == 0 {
		return nil
	}
	return t.UnmarshalJSON(tmp)
}

// Value implements the driver Valuer interface.
func (t JStrings) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.String(), nil
}

// JNumbers ===========数字列表===========
type JNumbers []int

func (t JNumbers) String() string {
	if t == nil {
		return ""
	}
	tmp, _ := json.Marshal([]int(t))
	return string(tmp)
}

// Contains 是否包含元素
func (t JNumbers) Contains(s int) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// MarshalJSON 转换为json类型
func (t JNumbers) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]int(t))
}

// UnmarshalJSON 不做处理
func (t *JNumbers) UnmarshalJSON(data []byte) error {
	var tmp []int
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = tmp
	return nil
}

// Scan implements the Scanner interface.
func (t *JNumbers) Scan(src interface{}) error {
	*t = make([]int, 0)
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read json int array from DB failed")
	}
	if len(tmp) == 0 {
		return nil
	}
	return t.UnmarshalJSON(tmp)
}

// Value implements the driver Valuer interface.
func (t JNumbers) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.String(), nil
}

// Hashes ===========哈希列表===========
type Hashes []uint64

func (t Hashes) String() string {
	if t == nil {
		return ""
	}
	tmp, _ := json.Marshal([]uint64(t))
	return string(tmp)
}

// Contains 是否包含元素
func (t Hashes) Contains(s uint64) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// MarshalJSON 转换为json类型
func (t Hashes) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]uint64(t))
}

// UnmarshalJSON 不做处理
func (t *Hashes) UnmarshalJSON(data []byte) error {
	var tmp []uint64
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = tmp
	return nil
}

// Scan implements the Scanner interface.
func (t *Hashes) Scan(src interface{}) error {
	*t = make([]uint64, 0)
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("read json int array from DB failed")
	}
	if len(tmp) == 0 {
		return nil
	}
	return t.UnmarshalJSON(tmp)
}

// Value implements the driver Valuer interface.
func (t Hashes) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.String(), nil
}
