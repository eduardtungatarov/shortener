package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

type OriginalURL struct{
	CorrelationId string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type ShortURL struct{
	CorrelationId string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
	Key string `json:"-"`
	OriginalURL string `json:"-"`
}

type Storage interface {
	Set(ctx context.Context, key, value string) error
	SetBatch(ctx context.Context, keyValues map[string]string) error
	Get(ctx context.Context, key string) (value string, ok bool)
	Ping(ctx context.Context) error
}

type Handler struct {
	ctx context.Context
	storage Storage
	baseURL string
	log *zap.SugaredLogger
}

func MakeHandler(ctx context.Context, storage Storage, baseURL string, log *zap.SugaredLogger) *Handler {
	return &Handler{
		ctx: ctx,
		storage: storage,
		baseURL: baseURL,
		log: log,
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
	err = h.storage.Set(h.ctx, key, string(body))
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

	url, ok := h.storage.Get(h.ctx, shortURL)
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
	err := h.storage.Set(h.ctx, key, reqStr.URL)
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
	err := h.storage.Ping(h.ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleShortenBatch(res http.ResponseWriter, req *http.Request) {
	var batch []OriginalURL

	defer req.Body.Close()
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&batch); err != nil {
		log.Printf("decode batch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}

	shortURLBatch := h.getShortURLBatch(batch)
	keyValueBatch := h.getKeyBatch(shortURLBatch)

	err := h.storage.SetBatch(h.ctx, keyValueBatch)
	if err != nil {
		log.Printf("storage SetBatch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}

	resp, err := json.Marshal(&shortURLBatch)
	if err != nil {
		log.Printf("marshal to shortURLBatch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(resp)
	if err != nil {
		log.Printf("response write: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return;
	}
}

func (h *Handler) getKeyBatch(batch []ShortURL) map[string]string {
	res := map[string]string{}
	for _, b := range batch {
		res[b.Key] = b.OriginalURL
	}
	return res
}

func (h *Handler) getShortURLBatch(batch []OriginalURL) []ShortURL {
	var res []ShortURL

	for _, b := range batch {
		key := h.getKey([]byte(b.OriginalURL))
		res = append(res, ShortURL{
			CorrelationId: b.CorrelationId,
			ShortURL: h.baseURL+"/"+key,
			Key: key,
			OriginalURL: b.OriginalURL,
		})
	}

	return res
}

func (h *Handler) getKey(url []byte) string {
	hash := md5.Sum(url)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	return key
}