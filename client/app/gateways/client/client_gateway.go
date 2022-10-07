package client

import (
	"bytes"
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/gateways"
	"io"
	"net/http"
)

type Gateway struct {
	Configuration *environment.Configuration
	HttpClient    *http.Client
}

func NewGateway(configuration *environment.Configuration, httpClient *http.Client) gateways.Gateway {
	return &Gateway{
		Configuration: configuration,
		HttpClient:    httpClient,
	}
}

func (c Gateway) NewRequest(method string, url string, body []byte) (*gateways.HttpResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &gateways.HttpResponse{
		ResponseBody: bodyBytes,
		StatusCode:   res.StatusCode,
	}, nil
}
