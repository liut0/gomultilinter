package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairs(t *testing.T) {
	expected := map[string]interface{}{"A": "B", "C": "D"}
	result := pairs("A", "B", "C", "D")
	assert.Equal(t, expected, result)
}

func TestPairsEmpty(t *testing.T) {
	assert.Equal(t, map[string]interface{}{}, pairs())
}

func TestPairsOdd(t *testing.T) {
	assert.Panics(t, func() {
		pairs("A", "B", "C")
	})
}
