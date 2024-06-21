package db

import (
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

func InitDB() (*gorm.DB, error) {
	pgConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  pgConn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.UserToken{},
		&models.Chat{},
		&models.Message{},
	)
	if err != nil {
		panic(err)
	}

	return database, nil
}
