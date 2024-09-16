package datastore

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func CreateNewDB() (*sql.DB, error) {
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")

	switch {
	case username == "":
		return nil, fmt.Errorf("переменная окружения username не установлена")

	case password == "":
		return nil, fmt.Errorf("переменная окружения password не установлена")

	case host == "":
		return nil, fmt.Errorf("переменная окружения host не установлена")

	case port == "":
		return nil, fmt.Errorf("переменная окружения port не установлена")

	case database == "":
		return nil, fmt.Errorf("переменная окружения database не установлена")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с базой данных: %v", err)
	}
	return dtb, nil
}
