package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func WithAuth() Middleware {
	return func(next http.Handler) http.Handler {
		return &AuthMiddleware{Next: next}
	}
}

func WithTimeout() Middleware {
	return func(next http.Handler) http.Handler {
		return TimeoutMiddleware{Next: next}
	}
}
