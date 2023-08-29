package crud

import (
	"errors"
	"fmt"
	"os"

	model "github.com/Prokuma/PLAccounting-Backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NoAuthorizationError = errors.New("No Authorization")

func InitDB() {
	// Load Environment Variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")

	// DB Connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Migration
	db.AutoMigrate(
		&model.User{}, &model.Application{}, &model.Permit{},

		&model.Book{}, &model.AccountTitle{}, &model.BookAuthorization{},
		&model.Transaction{}, &model.SubTransaction{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("db connected: ", &db)
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}
