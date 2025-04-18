package app

import "net/http"

func MainHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost && req.RequestURI == `/` {
		contentType := req.Header.Get(`Content-Type`)
		if contentType != `text/plain` {
			res.WriteHeader(http.StatusBadRequest)
			return;
		}
		res.WriteHeader(http.StatusCreated)
		_, _ = res.Write([]byte(`http://localhost:8080/EwHXdJfB`))
		return;
	}

	if req.Method == http.MethodGet && len(req.RequestURI) > 1 {
		res.Header().Add(`Location`, `https://practicum.yandex.ru/`)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return;
	}

	res.WriteHeader(http.StatusBadRequest)
}