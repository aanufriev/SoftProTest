package main

import (
	"fmt"
	"net/http"
)

const (
	lineProvider = "http://localhost:8000/api/v1/lines/"
)

type Handler struct{}

func (h Handler) checkReady(w http.ResponseWriter, r *http.Request) {
	sports := []string{"baseball", "football", "soccer"}

	for _, sport := range sports {
		_, err := http.Get(lineProvider + sport)
		if err != nil {
			// вывести ошибку
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("service is ready")
}

func main() {
	handler := Handler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/ready", handler.checkReady)

	http.ListenAndServe(":9000", mux)
}
