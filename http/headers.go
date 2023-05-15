package http

import "github.com/google/uuid"

func GetDefaultHeaders(correlationID uuid.UUID, step string) map[string][]string {
	return map[string][]string{
		"Content-Type":   {"application/json"},
		"Correlation_id": {correlationID.String()}, // Do not change the syntax
		"Step":           {step},
	}
}
