package services

import (
	"context"
	"fmt"
	"net/url"
)

type urlService struct {
	Client
}

func NewURLService(clientService Client) URL {
	return &urlService{clientService}
}

func (u urlService) OpenURL(ctx context.Context, address string, rawURL string) error {
	uri, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}
	if _, err := u.SendCommand(ctx, SendCommandInput{
		MacAddress: address,
		Request:    fmt.Sprintf("open-url %s", uri.String()),
	}); err != nil {
		return err
	}
	return nil
}
