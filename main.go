package main

import (
	"fmt"

	crud "github.com/Prokuma/PLAccounting-Backend/crud"
	endpoint "github.com/Prokuma/PLAccounting-Backend/endpoints"
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

	// Books
	r.POST("/book", endpoint.CreateBook())
	r.GET("/book/:bid", endpoint.GetBook())
	r.PATCH("/book/:bid", endpoint.UpdateBook())
	r.DELETE("/book/:bid", endpoint.DeleteBook())
	r.POST("/book/:bid/accountTitle", endpoint.CreateAccountTitle())
	r.GET("/book/:bid/accountTitle/:tid", endpoint.GetAccountTitle())
	r.PATCH("/book/:bid/accountTitle/:tid", endpoint.UpdateAccountTitle())
	r.DELETE("/book/:bid/accountTitle/:tid", endpoint.DeleteAccountTitle())

	// Transactions
	r.GET("/book/:bid/transaction", endpoint.GetTransactions(false))
	r.POST("/book/:bid/transaction", endpoint.CreateTransaction())
	r.GET("/book/:bid/transaction/:tid", endpoint.GetTransaction())
	r.PATCH("/book/:bid/transaction/:tid", endpoint.UpdateTransaction())
	r.DELETE("/book/:bid/transaction/:tid", endpoint.DeleteTransaction())
	r.GET("/book/:bid/transaction/page/:pid", endpoint.GetTransactions(true))

	// Execution
	r.Run()
}
