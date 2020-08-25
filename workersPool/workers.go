package workersPool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aanufriev/SoftProTest/storage"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	url     string
	storage storage.DatabaseInterface
}

func (w Worker) GetLine(sport string) {
	resp, err := http.Get(w.url + sport)
	if err != nil {
		logrus.WithError(err).Info("GetLine error")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Info("GetLine error")
		return
	}

	line, err := parseLine(body, sport)
	if err != nil {
		return
	}

	err = w.storage.Save(sport, line)
	if err != nil {
		return
	}
}

func parseLine(body []byte, sport string) (string, error) {
	var lines map[string]interface{}

	err := json.Unmarshal(body, &lines)
	if err != nil {
		return "", err
	}

	lines, ok := lines["lines"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("takeLine error")
	}

	line, ok := lines[strings.ToUpper(sport)].(string)
	if !ok {
		return "", fmt.Errorf("takeLine error")
	}

	return line, nil
}
