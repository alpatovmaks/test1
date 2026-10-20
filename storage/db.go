package storage

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite: один писатель — избегаем "database is locked"

	// Читаем и применяем миграцию
	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return nil, err
	}
	return db, nil
}