package backends

import "context"

type BackendType int

const (
	GRPC BackendType = iota
	HTTP
)

type Backend interface {
	ListenAndServe() error
	GracefulShutdown(ctx context.Context) error
}
