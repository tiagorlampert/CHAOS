package url

import "context"

type Service interface {
	OpenUrl(ctx context.Context, address string, rawUrl string) error
}
