package src

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	conf1 := GetConfig()
	conf2 := GetConfig()
	assert.Same(t, conf1, conf2)
}
