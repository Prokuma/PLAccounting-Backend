package main

import (
	"fmt"

	crud "github.com/Prokuma/ProkumaLabAccount-Backend/crud"
	endpoint "github.com/Prokuma/ProkumaLabAccount-Backend/endpoints"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env File
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env not found: ", err)
	}

	// DB Initilization
	crud.InitDB()

	// HTTP Endpoints Initilization
	r := gin.Default()

	// Server Health
	r.GET("/ping", endpoint.Ping())

	// Authentications
	r.POST("/user", endpoint.CreateUser())
	r.POST("/login", endpoint.Login())

	// Transactions
	r.POST("/book", endpoint.CreateBook())
	r.GET("/book/:bid", endpoint.GetBook())
	r.DELETE("/book/:bid", endpoint.DeleteBook())
	r.POST("/book/:bid/accountTitle", endpoint.CreateAccountTitle())
	r.GET("/book/:bid/accountTitle/:tid", endpoint.GetAccountTitle())
	r.DELETE("/book/:bid/accountTitle/:tid", endpoint.DeleteAccountTitle())

	// Execution
	r.Run()
}
