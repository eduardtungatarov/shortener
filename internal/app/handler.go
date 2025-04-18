package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var store map[string]string = make(map[string]string)

func MainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.RequestURI == `/` {
		handlePost(res, req)
		return;
	}

	if req.Method == http.MethodGet && len(req.RequestURI) > 1 {
		handleGet(res, req)
		return;
	}
	res.WriteHeader(http.StatusBadRequest)
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get(`Content-Type`), `text/plain`) {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	body, _ := io.ReadAll(req.Body)
	_ = req.Body.Close()

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	hash := md5.Sum(body)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	store[key] = string(body)

	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte(`http://localhost:8080/`+key))
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	shortUrl := req.RequestURI[1:]

	if strings.Contains(shortUrl, "/") {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	url, ok := store[shortUrl]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	res.Header().Add(`Location`, url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}