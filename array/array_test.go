package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	assert.Equal(t, Contains([]string{"abcd", "test", "go"}, "test"), true)
	assert.Equal(t, Contains([]int{1, 2, 3, 4, 5}, 3), true)
	assert.Equal(t, Contains([]bool{false, false}, true), false)
}
