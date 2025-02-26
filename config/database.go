package config

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func InitPostgres() {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// defer DB.Close()

	if err = DB.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		fullname VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		peer_id INTEGER UNIQUE NOT NULL
	);`
	if _, err = DB.Exec(createTableQuery); err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}
}
