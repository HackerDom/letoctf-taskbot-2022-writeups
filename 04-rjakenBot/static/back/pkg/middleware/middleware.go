package middleware

import (
	"github.com/valyala/fasthttp"
)

type Registry struct {
	middlewares []Middleware
}

func NewRegistry() *Registry {
	return &Registry{
		middlewares: nil,
	}
}

func (m *Registry) Register(f Middleware) {
	m.middlewares = append(m.middlewares, f)
}

func (m *Registry) Apply(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	resultHandler := handler
	for _, registered := range m.middlewares {
		resultHandler = registered.Apply(resultHandler)
	}

	return resultHandler
}

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler

func (m Middleware) Apply(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return m(handler)
}
