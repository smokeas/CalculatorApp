package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB инициализирует подключение к БД
// Используется PostgreSQL, строка подключения берётся из DATABASE_DSN
func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
