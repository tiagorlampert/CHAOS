package services

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type urlService struct {
	Client
}

func NewUrlService(clientService Client) Url {
	return &urlService{clientService}
}

func (u urlService) OpenUrl(ctx context.Context, address string, rawUrl string) error {
	if !strings.Contains(strings.ToLower(rawUrl), "http") {
		rawUrl = fmt.Sprintf("https://%s", rawUrl)
	}

	uri, err := url.Parse(rawUrl)
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
