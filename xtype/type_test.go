package xtype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	var src = Strings{"a", "b", "c"}
	var s1 = Strings{"c", "d", "e"}
	var s2 = Strings{"d", "e", "f"}

	assert.Equal(t, src.Contains("a"), true)
	assert.Equal(t, src.Contains("d"), false)
	assert.Equal(t, src.Intersectant(s1), true)
	assert.Equal(t, src.Intersectant(s2), false)
	assert.Equal(t, src.SAdd("a"), Strings{"a", "b", "c"})
	assert.Equal(t, src.SAdd("d"), Strings{"a", "b", "c", "d"})
	assert.Equal(t, src.Union(s1), Strings{"a", "b", "c", "d", "e"})
	assert.Equal(t, src.Union(s2), Strings{"a", "b", "c", "d", "e", "f"})
}
