package services

import "context"

type URL interface {
	OpenURL(ctx context.Context, address string, rawURL string) error
}
