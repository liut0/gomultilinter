package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNiceRegexp_FindNamedStringSubmatch(t *testing.T) {
	t.Parallel()

	rgx := MustCompile(`(?P<GRP1>[A-Za-z]*)(?P<GRP2>[0-9]+)`)

	result := rgx.FindNamedStringSubmatch("noMatch")
	assert.True(t, result == nil)

	result = rgx.FindNamedStringSubmatch("123")
	_, ok := result["GRP1"]
	assert.False(t, ok)
	assert.Equal(t, "123", result["GRP2"])

	result = rgx.FindNamedStringSubmatch("foo123")
	assert.Equal(t, "foo", result["GRP1"])
	assert.Equal(t, "123", result["GRP2"])
}
