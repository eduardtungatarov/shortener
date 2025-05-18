package middleware

import (
	"compress/gzip"
	"net/http"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func NewGzipResponseWriter(res http.ResponseWriter) *GzipResponseWriter {
	w := gzip.NewWriter(res)
	return &GzipResponseWriter{
		ResponseWriter: res,
		gzipWriter: w,
	}
}

func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

func (w *GzipResponseWriter) Close() error {
	return w.gzipWriter.Close()
}