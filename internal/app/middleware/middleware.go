// Package middleware http middleware для запросов.
package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Middleware посредник.
type Middleware struct {
	log *zap.SugaredLogger
}

// MakeMiddleware конструктор посредника.
func MakeMiddleware(log *zap.SugaredLogger) *Middleware {
	return &Middleware{
		log: log,
	}
}

// WithLog с логированием запросов и ответов.
func (m *Middleware) WithLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		uri := req.RequestURI
		method := req.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lres := &logResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}

		next.ServeHTTP(lres, req)

		duration := time.Since(start)

		m.log.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	})
}

// WithGzipResp с сжатием ответа.
func (m *Middleware) WithGzipResp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		oRes := res

		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			gRes, err := NewGzipResponseWriter(oRes, gzip.BestSpeed)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gRes.Close()
			oRes.Header().Set("Content-Encoding", "gzip")

			oRes = gRes
		}

		next.ServeHTTP(oRes, req)
	})
}

// WithGzipReq с разжатием респонса.
func (m *Middleware) WithGzipReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			gzipR, err := NewGzipReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gzipR.Close()

			req.Body = gzipR
		}

		next.ServeHTTP(res, req)
	})
}

// WithJSONReqCheck с проверкой на json content-type request'a.
func (m *Middleware) WithJSONReqCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			m.log.Infoln("Ожидался json тип запроса")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)
	})
}
