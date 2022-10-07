package request

type SendCommandRequestForm struct {
	Address   string `form:"address" binding:"required"`
	Command   string `form:"command" binding:"required"`
	Parameter string `form:"parameter"`
}

type RespondCommandRequestBody struct {
	ClientID string `json:"client_id,omitempty"`
	Response []byte `json:"response"`
	HasError bool   `json:"has_error,omitempty"`
}
