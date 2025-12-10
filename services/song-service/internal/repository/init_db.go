package repository

import (
	"fmt"
	"song-service/api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB() *gorm.DB {
	// Implementation for initializing PostgreSQL DB connection

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.POSTGRES_HOST,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
		config.POSTGRES_PORT,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}
