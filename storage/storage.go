package storage

import (
	"database/sql"
	"io/ioutil"
)

const (
	postgres = "postgres"
)

type StorageInterface interface {
	Ping() error
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

	ps.initDatabase()

	return nil
}

func (ps PostgresStorage) initDatabase() error {
	file, err := ioutil.ReadFile("./storage/init.sql")
	if err != nil {
		return err
	}

	_, err = ps.db.Query(string(file))
	if err != nil {
		return err
	}
	return nil
}

func (ps PostgresStorage) Ping() error {
	err := ps.db.Ping()
	return err
}
