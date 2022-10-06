package entities

type Payload struct {
	ClientID string `json:"client_id,omitempty"`
	Request  string `json:"request,omitempty"`
	Response []byte `json:"response,omitempty"`
	HasError bool   `json:"has_error,omitempty"`
}
