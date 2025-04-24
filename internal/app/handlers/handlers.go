package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	storage storage.Storage
	baseURL string
}

func MakeHandler(storage storage.Storage, baseURL string) *Handler {
	return &Handler{
		storage: storage,
		baseURL: baseURL,
	}
}

func (h *Handler) HandlePost(res http.ResponseWriter, req *http.Request) {
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
	_, _ = res.Write([]byte(h.baseURL+"/"+key))
}

func (h *Handler) HandleGet(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "shortUrl")

	url, ok := h.storage.Get(shortURL)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	res.Header().Add(`Location`, url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}