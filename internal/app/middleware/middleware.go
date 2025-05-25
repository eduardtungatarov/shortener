package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
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

func (m *Middleware) WithLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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

func (m *Middleware) WithGzipResp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		oRes := res

		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			gRes := NewGzipResponseWriter(oRes)
			defer gRes.Close()
			oRes.Header().Set("Content-Encoding", "gzip")

			oRes = gRes
		}

		next.ServeHTTP(oRes, req)
	})
}

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


func (m *Middleware) WithJsonReqCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			res.WriteHeader(http.StatusBadRequest)
			return;
		}

		next.ServeHTTP(res, req)
	})
}

