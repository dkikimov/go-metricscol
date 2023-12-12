package server

import (
	"net/http"
)

type Middleware func(http.HandlerFunc, *Config) http.HandlerFunc

func Conveyor(config *Config, h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h, config)
	}
	return h
}
