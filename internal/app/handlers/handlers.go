package handlers

import (
	"crypto/md5"
	"encoding/json"
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

func (h *Handler) HandlePost() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get(`Content-Type`), `text/plain`) && !strings.Contains(req.Header.Get(`Content-Type`), `application/x-gzip`) {
			fmt.Println(req.Header.Get(`Content-Type`))
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

		key := h.getKey(body)
		h.storage.Set(key, string(body))

		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(h.baseURL+"/"+key))
		if err != nil {
			log.Printf("Ошибка при записи ответа: %v", err)
			return;
		}
	}
}

func (h *Handler) HandleGet() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		shortURL := chi.URLParam(req, "shortUrl")

		url, ok := h.storage.Get(shortURL)
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return;
		}

		res.Header().Add(`Location`, url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *Handler) HandleShorten() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Обрабатываем запрос.
		if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			res.WriteHeader(http.StatusBadRequest)
			return;
		}

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
		h.storage.Set(key, reqStr.URL)

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
}

func (h *Handler) getKey(url []byte) string {
	hash := md5.Sum(url)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	return key
}