package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

type handler struct {
	storage storage.Storage
	host string
}

func MakeHandler(storage storage.Storage, host string) *handler {
	return &handler{
		storage: storage,
		host: host,
	}
}

func (h *handler) MainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.RequestURI == `/` {
		h.handlePost(res, req)
		return;
	}

	if req.Method == http.MethodGet && len(req.RequestURI) > 1 {
		h.handleGet(res, req)
		return;
	}
	res.WriteHeader(http.StatusBadRequest)
}

func (h *handler) handlePost(res http.ResponseWriter, req *http.Request) {
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
	h.storage.Set(key, string(body))

	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte("http://"+h.host+"/"+key))
}

func (h *handler) handleGet(res http.ResponseWriter, req *http.Request) {
	shortURL := req.RequestURI[1:]

	if strings.Contains(shortURL, "/") {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	url, ok := h.storage.Get(shortURL)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	res.Header().Add(`Location`, url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}