package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

type OriginalURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	Key           string `json:"-"`
	OriginalURL   string `json:"-"`
}

type DeleteRequest struct {
	UserID string
	Urls []string
}


type Storage interface {
	Set(ctx context.Context, key, value string) error
	SetBatch(ctx context.Context, keyValues map[string]string) error
	DeleteBatch(ctx context.Context, keys []string, userID string) error
	Get(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
	GetByUserID(ctx context.Context) ([]map[string]string, error)
}

type Handler struct {
	storage Storage
	baseURL string
	log     *zap.SugaredLogger
	deleteCh chan DeleteRequest
}

func MakeHandler(storage Storage, baseURL string, log *zap.SugaredLogger) *Handler {
	return &Handler{
		storage: storage,
		baseURL: baseURL,
		log:     log,
		deleteCh: make(chan DeleteRequest, 1024),
	}
}

func (h *Handler) HandlePost(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Не удалось прочитать тело запроса: %v", err)
		return
	}
	defer req.Body.Close()

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	key := h.getKey(body)
	err = h.storage.Set(req.Context(), key, string(body))
	isConflict := errors.Is(err, storage.ErrConflict)
	if err != nil && !isConflict {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Не удалось сохранить url: %v", err)
		return
	}

	if isConflict {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	_, err = res.Write([]byte(h.baseURL + "/" + key))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Ошибка при записи ответа: %v", err)
		return
	}
}

func (h *Handler) HandleGet(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "shortUrl")

	url, err := h.storage.Get(req.Context(), shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrDeleted) {
			res.WriteHeader(http.StatusGone)
			return
		}

		log.Printf("Ошибка при получении ссылки: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
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
		return
	}

	if reqStr.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// Сохраняем url.
	key := h.getKey([]byte(reqStr.URL))
	err := h.storage.Set(req.Context(), key, reqStr.URL)
	isConflict := errors.Is(err, storage.ErrConflict)
	if err != nil && !isConflict {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Формируем ответ.
	respStr := struct {
		Result string `json:"result"`
	}{}
	respStr.Result = h.baseURL + "/" + key

	res.Header().Set("Content-Type", "application/json")

	if isConflict {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(respStr); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetPing(res http.ResponseWriter, req *http.Request) {
	err := h.storage.Ping(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusServiceUnavailable)
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
		return
	}

	shortURLBatch := h.getShortURLBatch(batch)
	keyValueBatch := h.getKeyBatch(shortURLBatch)

	err := h.storage.SetBatch(req.Context(), keyValueBatch)
	if err != nil {
		log.Printf("storage SetBatch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&shortURLBatch)
	if err != nil {
		log.Printf("marshal to shortURLBatch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(resp)
	if err != nil {
		log.Printf("response write: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetUserUrls(res http.ResponseWriter, req *http.Request) {
	urls, err := h.storage.GetByUserID(req.Context())
	if err != nil {
		log.Printf("storage GetByUserId: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if len(urls) == 0 {
		res.WriteHeader(http.StatusNoContent)
	}

	type respStr struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var respStrSlice []respStr

	for _, v := range urls {
		respStrSlice = append(respStrSlice, respStr{
			OriginalURL: v["original_url"],
			ShortURL:    h.baseURL + "/" + v["short_url"],
		})
	}

	body, err := json.Marshal(&respStrSlice)
	if err != nil {
		log.Printf("marshal to shortURLBatch: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = res.Write(body)
	if err != nil {
		log.Printf("response write: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDeleteUserUrls(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	respStr := []string{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&respStr)
	if err != nil {
		log.Printf("unmarshal body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	userID, ok := ctx.Value(config.UserIDKeyName).(string)
	if !ok  {
		log.Printf("userID not found: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
	}

	h.deleteCh <- DeleteRequest{
		UserID: userID,
		Urls: respStr,
	}

	res.WriteHeader(http.StatusAccepted)
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
			CorrelationID: b.CorrelationID,
			ShortURL:      h.baseURL + "/" + key,
			Key:           key,
			OriginalURL:   b.OriginalURL,
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
