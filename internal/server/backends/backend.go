package backends

import "context"

type BackendType int

const (
	GRPCType BackendType = iota
	HTTPType
)

type Backend interface {
	ListenAndServe() error
	GracefulShutdown(ctx context.Context) error
}
