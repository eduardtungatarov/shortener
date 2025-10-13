package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ExampleHandler_HandlePost() {
	// URL для сокращения.
	urlToShort := "https://practicum.yandex.ru/"

	// Строим запрос post /
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		strings.NewReader(urlToShort),
	)

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Обрабатываем успешный ответ.
	if resp.StatusCode == http.StatusCreated {
		// В теле ответа находится сокращенная ссылка.
		shortURL, _ := io.ReadAll(resp.Body)
		fmt.Println(shortURL) // http://localhost:8080/e9fb5de
	}
}

func ExampleHandler_HandleGet() {
	// Сокращенный URL.
	shortURL := "http://localhost:8080/e9fb5de"

	// Строим запрос get /e9fb5de
	req, _ := http.NewRequest(
		http.MethodGet,
		shortURL,
		nil,
	)

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Обрабатываем успешный редирект.
	if resp.StatusCode == http.StatusTemporaryRedirect {
		// В хедере Location будет хранится целевой URL.
		url := resp.Header.Get("Location")
		fmt.Println(url) // https://practicum.yandex.ru/
	}
}

func ExampleHandler_HandleShorten() {
	// URL для сокращения.
	urlToShort := "https://practicum.yandex.ru/"

	// Строим запрос post /
	jsonBody := `{"url":"` + urlToShort + `"}`
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		strings.NewReader(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Обрабатываем успешный json ответ.
	if resp.StatusCode == http.StatusCreated {
		// В теле ответа находится json с сокращенной ссылкой.
		jsonBody, _ := io.ReadAll(resp.Body)

		reqStr := struct {
			URL string `json:"url"`
		}{}

		_ = json.Unmarshal(jsonBody, &reqStr)

		fmt.Println(reqStr.URL) // http://localhost:8080/e9fb5de
	}
}

func ExampleHandler_HandleGetPing() {
	// Строим запрос get /ping
	req, _ := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/ping",
		nil,
	)

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Сервис и его подсистемы доступны.
	if resp.StatusCode == http.StatusOK {
		//
	}
}

func ExampleHandler_HandleShortenBatch() {
	// urls json с ссылками на сокращение.
	urls := `
	[
	  {
		"correlation_id": "1",
		"original_url": "https://ya.ru"
	  },
	  {
		"correlation_id": "2",
		"original_url": "https://r0.ru"
	  }
	]
	`

	// Строим запрос post /api/shorten/batch
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten/batch",
		strings.NewReader(urls),
	)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Обрабатываем успешный json ответ.
	if resp.StatusCode == http.StatusCreated {
		var respStr []struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}

		body, _ := io.ReadAll(resp.Body)

		_ = json.Unmarshal(body, &respStr)

		// [
		//  {
		//    "correlation_id": "1",
		//    "short_url": "http://localhost:8080/f6ca88e"
		//  },
		//  {
		//    "correlation_id": "2",
		//    "short_url": "http://localhost:8080/2cd23b0"
		//  }
		//]
		fmt.Println(respStr)
	}
}

func ExampleHandler_HandleGetUserUrls() {
	// Строим запрос post /api/user/urls
	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/user/urls",
		nil,
	)

	// Отправляем запрос.
	client := http.Client{}
	resp, _ := client.Do(req)

	// Обрабатываем успешный json ответ.
	if resp.StatusCode == http.StatusOK {
		var respStr struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}

		body, _ := io.ReadAll(resp.Body)

		_ = json.Unmarshal(body, &respStr)

		// [
		//  {
		//    "original_url": "http://ya.ru",
		//    "short_url": "e520966"
		//  }
		//]
		fmt.Println(respStr)
	}
}
