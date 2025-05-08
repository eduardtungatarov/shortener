package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
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
		start := time.Now()

		uri := req.RequestURI
		method := req.Method

		responseData := &responseData {
			status: 0,
			size: 0,
		}
		lres := &logResponseWriter{
			ResponseWriter: res,
			responseData: responseData,
		}

		handler.ServeHTTP(lres, req)

		duration := time.Since(start)

		m.log.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
}
