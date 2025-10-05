package middleware

import (
	"compress/gzip"
	"net/http"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func NewGzipResponseWriter(res http.ResponseWriter, level int) (*GzipResponseWriter, error) {
	w, err := gzip.NewWriterLevel(res, level)
	if err != nil {
		return nil, err
	}

	return &GzipResponseWriter{
		ResponseWriter: res,
		gzipWriter:     w,
	}, nil
}

func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

func (w *GzipResponseWriter) Close() error {
	return w.gzipWriter.Close()
}
