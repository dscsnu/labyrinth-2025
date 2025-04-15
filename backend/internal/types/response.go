package types

import "encoding/json"

type ApiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
