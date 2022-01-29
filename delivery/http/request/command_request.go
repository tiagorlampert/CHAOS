package request

type SendCommandRequestForm struct {
	Address string `form:"address" binding:"required"`
	Command string `form:"command" binding:"required"`
}

type RespondCommandRequestBody struct {
	MacAddress string `json:"mac_address" binding:"required"`
	Response   []byte `json:"response"`
	HasError   bool   `json:"has_error,omitempty"`
}
