package middleware

import (
	"compress/gzip"
	"net/http"
)

// GzipResponseWriter сжиматель.
type GzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

// NewGzipResponseWriter конструктор сжимателя.
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

// Write сжимает.
func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

// Close закрытие сжимателя.
func (w *GzipResponseWriter) Close() error {
	return w.gzipWriter.Close()
}
