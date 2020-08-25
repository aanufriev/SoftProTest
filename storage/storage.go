package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	postgres = "postgres"
)

type StorageInterface interface {
	Ping() error
	Save(sport string, value string) error
	GetLastLine(sport string) (float32, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (ps *PostgresStorage) Open(dataSourceName string) error {
	db, err := sql.Open(postgres, dataSourceName)
	if err != nil {
		return err
	}
	ps.db = db

	err = db.Ping()
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"database":         postgres,
		"data source name": dataSourceName,
	}).Info("Open db connection")

	return nil
}

func (ps PostgresStorage) InitDatabase(sports []string) error {
	for _, sport := range sports {
		_, err := ps.db.Exec(
			`CREATE TABLE IF NOT EXISTS ` + sport + ` (
			id SERIAL NOT NULL PRIMARY KEY,
			line REAL NOT NULL,
			get_at_time TIMESTAMP NOT NULL
			);`)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"sports": sports,
	}).Info("Init database")

	return nil
}

func (ps PostgresStorage) Ping() error {
	err := ps.db.Ping()
	return err
}

func (ps PostgresStorage) Save(sport string, line string) error {
	_, err := ps.db.Exec(
		"INSERT INTO "+sport+" (line, get_at_time) VALUES ($1, $2)",
		line, time.Now(),
	)

	logrus.WithFields(logrus.Fields{
		"sport": sport,
		"line":  line,
	}).Info("Save line in database")
	return err
}

func (ps PostgresStorage) GetLastLine(sport string) (float32, error) {
	var line float32

	err := ps.db.QueryRow(
		"SELECT line FROM " + sport + " ORDER BY get_at_time DESC LIMIT 1",
	).Scan(&line)
	if err != nil {
		return 0, err
	}

	return line, nil
}
