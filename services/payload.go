package services

type PayloadData struct {
	Request     string `json:"request,omitempty"`
	Response    []byte `json:"response,omitempty"`
	HasError    bool   `json:"has_error,omitempty"`
	HasResponse bool   `json:"has_response,omitempty"`
}

type Payload interface {
	Set(string, *PayloadData)
	Get(string) (*PayloadData, bool)
	Remove(string)
}
