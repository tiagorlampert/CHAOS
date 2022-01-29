package services

import "context"

type Url interface {
	OpenUrl(ctx context.Context, address string, rawUrl string) error
}
