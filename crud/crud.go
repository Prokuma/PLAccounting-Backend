package crud

import (
	"errors"
	"fmt"
	"os"

	model "github.com/Prokuma/ProkumaLabAccount-Backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NoAuthorizationError = errors.New("No Authorization")

func InitDB() {
	// Load Environment Variables
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	dbname := os.Getenv("POSTGRESQL_DBNAME")
	sslmode := os.Getenv("POSTGRESQL_SSLMODE")
	timezone := os.Getenv("POSTGRESQL_TIMEZONE")

	// DB Connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Migration
	db.AutoMigrate(&model.SubTransaction{})
	db.AutoMigrate(&model.Permit{})
	db.AutoMigrate(&model.BookAuthorization{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Application{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Book{})
	db.AutoMigrate(&model.AccountTitle{})

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
