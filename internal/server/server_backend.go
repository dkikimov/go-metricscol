package server

import "net/http"

type BackendType int

const (
	GRPC BackendType = iota
	HTTP
)

type Backend interface {
	GetRouter() http.Handler
}
