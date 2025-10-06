package middleware

import (
	"compress/gzip"
	"net/http/httptest"
	"testing"
)

func BenchmarkWriteNoCompression(b *testing.B) {
	runWithLevel(b, gzip.NoCompression)
}

func BenchmarkWriteBestSpeed(b *testing.B) {
	runWithLevel(b, gzip.BestSpeed)
}

func BenchmarkWriteBestCompression(b *testing.B) {
	runWithLevel(b, gzip.BestCompression)
}

func BenchmarkWriteDefaultCompression(b *testing.B) {
	runWithLevel(b, gzip.DefaultCompression)
}

func BenchmarkWriteHuffmanOnly(b *testing.B) {
	runWithLevel(b, gzip.HuffmanOnly)
}

func runWithLevel(b *testing.B, level int) {
	w := httptest.NewRecorder()
	g, _ := NewGzipResponseWriter(w, level)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		_, _ = g.Write([]byte(`[{"short_url":"http://localhost:8080/f6ca88e","original_url":"http://bujdyzmn.biz/p4atoih/fmk4d4w/gbux5"},{"short_url":"http://localhost:8080/2cd23b0","original_url":"http://scm9jddxmi1y.net/dl11ky/agl7mqpqtb"},{"short_url":"http://localhost:8080/f6ca88e","original_url":"http://bujdyzmn.biz/p4atoih/fmk4d4w/gbux5"},{"short_url":"http://localhost:8080/2cd23b0","original_url":"http://scm9jddxmi1y.net/dl11ky/agl7mqpqtb"},{"short_url":"http://localhost:8080/f6ca88e","original_url":"http://bujdyzmn.biz/p4atoih/fmk4d4w/gbux5"},{"short_url":"http://localhost:8080/2cd23b0","original_url":"http://scm9jddxmi1y.net/dl11ky/agl7mqpqtb"},{"short_url":"http://localhost:8080/e9fb5de","original_url":"http://ydaddddasddsdddd.ru"}]`))
	}
}
