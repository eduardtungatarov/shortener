package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type Storage interface {
	Set(key, value string) error
	Get(key string) (value string, ok bool)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	ctx context.Context
	storage Storage
	baseURL string
	pinger Pinger
}

func MakeHandler(ctx context.Context, storage Storage, baseURL string, pinger Pinger) *Handler {
	return &Handler{
		ctx: ctx,
		storage: storage,
		baseURL: baseURL,
		pinger: pinger,
	}
}

func (h *Handler) HandlePost(res http.ResponseWriter, req *http.Request)  {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Не удалось прочитать тело запроса: %v", err)
		return;
	}
	defer req.Body.Close()

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	key := h.getKey(body)
	err = h.storage.Set(key, string(body))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Не удалось сохранить url: %v", err)
		return;
	}

	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(h.baseURL+"/"+key))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
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

func (h *Handler) HandleShorten(res http.ResponseWriter, req *http.Request) {
	reqStr := struct {
		URL string `json:"url"`
	}{}

	defer req.Body.Close()
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqStr); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}

	if reqStr.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return;
	}

	// Сохраняем url.
	key := h.getKey([]byte(reqStr.URL))
	err := h.storage.Set(key, reqStr.URL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}

	// Формируем ответ.
	respStr := struct {
		Result string `json:"result"`
	}{}
	respStr.Result = h.baseURL+"/"+key

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	if err := enc.Encode(respStr); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}
}

func (h *Handler) HandleGetPing(res http.ResponseWriter, req *http.Request) {
	err := h.pinger.Ping(h.ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *Handler) getKey(url []byte) string {
	hash := md5.Sum(url)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	return key
}