package health

import "context"

type UseCase interface {
	Ping(ctx context.Context) error
}
