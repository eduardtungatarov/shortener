package middleware

import (
	"compress/gzip"
	"io"
)

// GzipReader разжиматель.
type GzipReader struct {
	io.ReadCloser
	gzip *gzip.Reader
}

// NewGzipReader конструктор разжимателя.
func NewGzipReader(reader io.ReadCloser) (*GzipReader, error) {
	gzipR, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
			reader,
			gzipR,
		},
		nil
}

// Read разжимает.
func (r *GzipReader) Read(b []byte) (n int, err error) {
	return r.gzip.Read(b)
}

// Close закрытие разжимателя.
func (r *GzipReader) Close() error {
	if err := r.ReadCloser.Close(); err != nil {
		return err
	}
	return r.gzip.Close()
}
