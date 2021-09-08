package entities

type Payload struct {
	MacAddress string `json:"mac_address,omitempty"`
	Request    string `json:"request,omitempty"`
	Response   []byte `json:"response,omitempty"`
	HasError   bool   `json:"has_error,omitempty"`
}
