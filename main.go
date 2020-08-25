package main

import (
	"io/ioutil"

	"github.com/aanufriev/SoftProTest/config"
	grpcserver "github.com/aanufriev/SoftProTest/grpcServer"
	httpserver "github.com/aanufriev/SoftProTest/httpServer"
	"github.com/aanufriev/SoftProTest/storage"
	workerspool "github.com/aanufriev/SoftProTest/workersPool"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func setLogrusLevel(level string) {
	levels := map[string]logrus.Level{
		"fatal": logrus.FatalLevel,
		"debug": logrus.DebugLevel,
		"panic": logrus.PanicLevel,
		"info":  logrus.InfoLevel,
		"error": logrus.ErrorLevel,
		"trace": logrus.TraceLevel,
		"warn":  logrus.WarnLevel,
	}

	logrusLevel, ok := levels[level]
	if !ok {
		logrusLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logrusLevel)
}

func main() {
	config := config.Config{}
	configData, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		logrus.Fatal("can't read config: ", err)
	}
	config.UnmarshalJSON(configData)

	setLogrusLevel(config.LogLevel)

	logrus.WithFields(logrus.Fields{
		"config": config,
	}).Info("Got config")

	storage := &storage.PostgresStorage{}
	err = storage.Open(config.DBDataSource)
	if err != nil {
		logrus.Fatal("can't open database connection: ", err)
	}

	sports := config.LinesProvider.Sports
	storage.InitDatabase(sports)

	intervals := config.LinesProvider.Intervals
	linesProviderURL := config.LinesProvider.URL

	go func() {
		pool := workerspool.NewWorkersPool(1, sports, intervals, storage)
		pool.Start(linesProviderURL)
	}()

	go httpserver.StartHTTPServer(config.HTTPPort, storage, sports, linesProviderURL)

	grpcserver.StartSubServer(config.GrpcPort, storage)
}
