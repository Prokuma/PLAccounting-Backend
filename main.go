package main

import (
	"fmt"
	"time"

	docs "github.com/Prokuma/PLAccounting-Backend/docs"
	"github.com/gin-contrib/cors"

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

	v1 := r.Group("/api/v1")
	{
		// Server Health
		v1.GET("/ping", endpoint.Ping)

		// Authentications
		v1.POST("/user", endpoint.CreateUser)
		v1.POST("/login", endpoint.Login)

		// Books
		v1.GET("/book", endpoint.GetAllBooks)
		v1.POST("/book", endpoint.CreateBook)
		v1.GET("/book/:bid", endpoint.GetBook)
		v1.PATCH("/book/:bid", endpoint.UpdateBook)
		v1.DELETE("/book/:bid", endpoint.DeleteBook)

		// Account Titles
		v1.GET("/book/:bid/accountTitle", endpoint.GetAllAccountTitles)
		v1.POST("/book/:bid/accountTitle", endpoint.CreateAccountTitle)
		v1.GET("/book/:bid/accountTitle/:tid", endpoint.GetAccountTitle)
		v1.PATCH("/book/:bid/accountTitle/:tid", endpoint.UpdateAccountTitle)
		v1.DELETE("/book/:bid/accountTitle/:tid", endpoint.DeleteAccountTitle)

		// Transactions
		v1.GET("/book/:bid/transaction", endpoint.GetTransactions)
		v1.POST("/book/:bid/transaction", endpoint.CreateTransaction)
		v1.GET("/book/:bid/transaction/:tid", endpoint.GetTransaction)
		v1.PATCH("/book/:bid/transaction/:tid", endpoint.UpdateTransaction)
		v1.DELETE("/book/:bid/transaction/:tid", endpoint.DeleteTransaction)
		v1.GET("/book/:bid/transaction/page/:pid", endpoint.GetTransactionsWithPage)
		v1.GET("/book/:bid/accountTitle/:tid/transactions", endpoint.GetSubTransactionsFromAccountTitle)
		v1.GET("/book/:bid/accountTitle/:tid/transactions/:pid", endpoint.GetSubTransactionsFromAccountTitleWithPage)
	}

	// 本登録
	r.GET("/createUser", endpoint.CreateUserAtDatabase)

	// Swgger
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "PLAccounting API"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// CORS (Not for development)
	if gin.Mode() == gin.ReleaseMode {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost:3000",
				"https://example.com",
			},
			AllowMethods: []string{
				"GET",
				"POST",
				"PUT",
				"PATCH",
				"DELETE",
				"HEAD",
				"OPTIONS",
			},
			AllowHeaders: []string{
				"Access-Control-Allow-Headers",
				"Access-Control-Allow-Credentials",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"Autorization",
			},
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
		}))
	}

	// Execution
	r.Run()
}
