package middleware

import (
	"compress/gzip"
	"io"
)

type GzipReader struct {
	io.ReadCloser
	gzip *gzip.Reader
}

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

func (r *GzipReader) Read(b []byte) (n int, err error) {
	return r.gzip.Read(b)
}

func (r *GzipReader) Close() error {
	if err := r.ReadCloser.Close(); err != nil {
		return err
	}
	return r.gzip.Close()
}
