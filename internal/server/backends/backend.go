package backends

import "context"

type BackendType int

const (
	GRPCBackend BackendType = iota
	HTTPBackend
)

type Backend interface {
	ListenAndServe() error
	GracefulShutdown(ctx context.Context) error
}
