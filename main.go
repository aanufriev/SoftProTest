package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aanufriev/SoftProTest/storage"
)

const (
	lineProvider = "http://localhost:8000/api/v1/lines/"
)

type Handler struct {
	storage storage.StorageInterface
}

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

	err := h.storage.Ping()
	if err != nil {
		// вывести ошибку
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("service is ready")
}

func main() {
	storage := &storage.PostgresStorage{}
	err := storage.Open("user=testuser password=test_password dbname=softpro sslmode=disable")
	if err != nil {
		log.Fatal("can't open database connection: ", err)
	}

	handler := &Handler{
		storage,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ready", handler.checkReady)

	http.ListenAndServe(":9000", mux)
}
