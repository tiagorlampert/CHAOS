package services

import "context"

type URL interface {
	OpenUrl(ctx context.Context, address string, rawUrl string) error
}
