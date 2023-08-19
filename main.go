package main

import (
	"fmt"

	docs "github.com/Prokuma/PLAccounting-Backend/docs"

	crud "github.com/Prokuma/PLAccounting-Backend/crud"
	endpoint "github.com/Prokuma/PLAccounting-Backend/endpoints"
	util "github.com/Prokuma/PLAccounting-Backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PLAccounting API
// @version v1
// @license AGPLv3
// @description This is a PLAccounting API Server.
func main() {
	// Load .env File
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env not found: ", err)
	}

	// DB Initilization
	crud.InitDB()
	util.InitRedis()

	// HTTP Endpoints Initilization
	r := gin.Default()

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "PLAccounting API"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/")
		{

			// Server Health
			eg.GET("/ping", endpoint.Ping)

			// Authentications
			eg.POST("/user", endpoint.CreateUser)
			eg.POST("/login", endpoint.Login)

			// Books
			eg.POST("/book", endpoint.CreateBook)
			eg.GET("/book/:bid", endpoint.GetBook)
			eg.PATCH("/book/:bid", endpoint.UpdateBook)
			eg.DELETE("/book/:bid", endpoint.DeleteBook)
			eg.POST("/book/:bid/accountTitle", endpoint.CreateAccountTitle)
			eg.GET("/book/:bid/accountTitle/:tid", endpoint.GetAccountTitle)
			eg.PATCH("/book/:bid/accountTitle/:tid", endpoint.UpdateAccountTitle)
			eg.DELETE("/book/:bid/accountTitle/:tid", endpoint.DeleteAccountTitle)

			// Transactions
			eg.GET("/book/:bid/transaction", endpoint.GetTransactions)
			eg.POST("/book/:bid/transaction", endpoint.CreateTransaction)
			eg.GET("/book/:bid/transaction/:tid", endpoint.GetTransaction)
			eg.PATCH("/book/:bid/transaction/:tid", endpoint.UpdateTransaction)
			eg.DELETE("/book/:bid/transaction/:tid", endpoint.DeleteTransaction)
			eg.GET("/book/:bid/transaction/page/:pid", endpoint.GetTransactionsWithPage)
			eg.GET("/:bid", endpoint.GetBook)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Execution
	r.Run()
}
