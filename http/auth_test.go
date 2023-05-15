package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBasicAuth(t *testing.T) {
	expected := map[string]string{
		"username": "mock_username",
		"password": "mock_password",
	}

	auth := GetBasicAuth("mock_username", "mock_password")

	assert.Equal(t, expected, auth)
}
