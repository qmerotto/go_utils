package http

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDefaultHeaders(t *testing.T) {
	expected := map[string][]string{
		"Content-Type":   {"application/json"},
		"Correlation_id": {"e785b679-b18b-4e0e-bf2a-3a3fb78e9c2e"},
		"Step":           {"0"},
	}

	correlationID, _ := uuid.Parse("e785b679-b18b-4e0e-bf2a-3a3fb78e9c2e")

	headers := GetDefaultHeaders(correlationID, "0")

	assert.Equal(t, expected, headers)
}
