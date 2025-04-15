package types

type ApiResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Payload map[string]interface{} `json:"payload"`
}
