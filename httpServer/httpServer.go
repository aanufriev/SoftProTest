package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aanufriev/SoftProTest/storage"
	"github.com/sirupsen/logrus"
)

func checklinesProvider(url string, sports []string) error {
	for _, sport := range sports {
		_, err := http.Get(url + sport)
		if err != nil {
			logrus.WithError(err).Info("lines provider")
			return fmt.Errorf("lines provider is not available, error: %v", err)
		}
	}

	return nil
}

func checkStorage(storage storage.StorageInterface) error {
	err := storage.Ping()
	if err != nil {
		logrus.WithError(err).Info("storage")
		return fmt.Errorf("storage is not available, error: %v", err)
	}

	return nil
}

type handler struct {
	storage          storage.StorageInterface
	linesProviderURL string
	sports           []string
}

func (h handler) writeSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"data": "service is available"}
	json.NewEncoder(w).Encode(resp)
}

func (h handler) writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	resp := map[string]string{"error": err.Error()}
	json.NewEncoder(w).Encode(resp)
}

func (h handler) checkReady(w http.ResponseWriter, r *http.Request) {
	err := checklinesProvider(h.linesProviderURL, h.sports)
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = checkStorage(h.storage)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeSuccess(w)
}

func StartHTTPServer(port string, storage storage.StorageInterface, sports []string, linesProviderURL string) {
	handler := &handler{
		storage:          storage,
		linesProviderURL: linesProviderURL,
		sports:           sports,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ready", handler.checkReady)

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting http server")

	http.ListenAndServe(port, mux)
}
