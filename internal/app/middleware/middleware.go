package middleware

import (
	"go.uber.org/zap"
	"net/http"
)

type Middleware struct {
	log *zap.SugaredLogger
}

func MakeMiddleware(log *zap.SugaredLogger) *Middleware {
	return &Middleware{
		log: log,
	}
}

func (m *Middleware) WithLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req)
	}
}
