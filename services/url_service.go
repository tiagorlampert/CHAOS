package services

import (
	"context"
	"fmt"
	"net/url"
)

type urlService struct {
	Client
}

func NewUrlService(clientService Client) URL {
	return &urlService{clientService}
}

func (u urlService) OpenUrl(ctx context.Context, address string, rawUrl string) error {
	//TODO add http protocol on url if doesnt contains
	uri, err := url.ParseRequestURI(rawUrl)
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
