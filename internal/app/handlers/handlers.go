package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strings"
)

type Storage interface {
	Set(key, value string)
	Get(key string) (value string, ok bool)
}

type Handler struct {
	storage Storage
	baseURL string
}

func MakeHandler(storage Storage, baseURL string) *Handler {
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

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса: %v", err)
		return;
	}
	defer req.Body.Close()

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	hash := md5.Sum(body)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	h.storage.Set(key, string(body))

	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(h.baseURL+"/"+key))
	if err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
		return;
	}
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