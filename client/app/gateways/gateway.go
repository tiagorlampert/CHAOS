package gateways

type HttpResponse struct {
	ResponseBody []byte
	StatusCode   int
}

type Gateway interface {
	NewRequest(method string, url string, body []byte) (*HttpResponse, error)
}
