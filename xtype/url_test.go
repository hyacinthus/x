package xtype

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用结构体
type ts struct {
	URL ImageURL `json:"url"`
}

func TestURL(t *testing.T) {
	var target = "https://image.xuebaox.com/company/cover/a.jpg"
	var s1 = ImageURL("company/cover/a.jpg")
	var s2 = ImageURL("/company/cover/a.jpg")
	var s3 = ImageURL("./company/cover/a.jpg")
	var s4 = ImageURL("company/cover/a.jpg?b=c")

	assert.Equal(t, target, s1.String())
	assert.Equal(t, target, s2.String())
	assert.Equal(t, target, s3.String())
	assert.Equal(t, target+"?b=c", s4.String())
	var j = new(ts)
	err := json.Unmarshal([]byte(fmt.Sprintf(`{"url":"%s"}`, target)), j)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, s2, j.URL)
}
