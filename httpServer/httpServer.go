package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aanufriev/SoftProTest/storage"
	"github.com/sirupsen/logrus"
)

// InfiniteCheckLinesProvider checks availability endlessly
// to start GRPC server when the resource becomes available
func InfiniteCheckLinesProvider(url string, sports []string, timeout int) bool {
	for {
		err := checklinesProvider(url, sports)
		if err == nil {
			return true
		}

		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

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

func checkStorage(storage storage.DatabaseInterface) error {
	err := storage.Ping()
	if err != nil {
		logrus.WithError(err).Info("storage")
		return fmt.Errorf("storage is not available, error: %v", err)
	}

	return nil
}

type handler struct {
	storage          storage.DatabaseInterface
	linesProviderURL string
	sports           []string
}

func (h handler) writeSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"data": "service is available"}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logrus.WithError(err).Info("writeSuccess")
	}
}

func (h handler) writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	resp := map[string]string{"error": err.Error()}
	funcErr := json.NewEncoder(w).Encode(resp)
	if funcErr != nil {
		logrus.WithError(funcErr).Info("writeError")
	}
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

// StartHTTPServer creates a http server
// to handle request on /ready
func StartHTTPServer(port string, storage storage.DatabaseInterface, sports []string, linesProviderURL string) {
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

	err := http.ListenAndServe(port, mux)
	if err != nil {
		logrus.WithError(err).Info("HTTP server error")
	}
}
