package endpoint

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Prokuma/PLAccounting-Backend/crud"
	model "github.com/Prokuma/PLAccounting-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		book, err := crud.GetBook(c.Param("bid"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Book was not found")
			c.Abort()
			return
		}

		bookAuthorization, err := crud.GetBookAuthorization(&user, &book)
		if err != nil {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}
		if strings.Index(bookAuthorization.Authority, "write") == -1 {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}

		type CreateTransaction struct {
			Description     string                 `json:"description" binding:"required"`
			OccuredAt       time.Time              `json:"occured_at" binding:"required"`
			SubTransactions []model.SubTransaction `json:"sub_transactions" binding:"required"`
		}
		var createTransaction CreateTransaction
		err = c.BindJSON(&createTransaction)

		if err != nil {
			c.String(http.StatusBadRequest, "Infection Informations")
			c.Abort()
			return
		}

		for idx := range createTransaction.SubTransactions {
			createTransaction.SubTransactions[idx].BookId = book.BookId
		}

		var transaction = model.Transaction{
			BookId:          book.BookId,
			Description:     createTransaction.Description,
			OccuredAt:       createTransaction.OccuredAt,
			SubTransactions: createTransaction.SubTransactions,
		}

		err = crud.CreateTransaction(&transaction)
		if err != nil {
			c.String(http.StatusInternalServerError, "Transactiuon could not created")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transaction": transaction,
			"message":     "Transaction was created",
		})
	}
}

func UpdateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		book, err := crud.GetBook(c.Param("bid"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Book was not found")
			c.Abort()
			return
		}

		bookAuthorization, err := crud.GetBookAuthorization(&user, &book)
		if err != nil {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}
		if strings.Index(bookAuthorization.Authority, "update") == -1 {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}

		type UpdateTransaction struct {
			Description     *string                 `json:"description"`
			OccuredAt       *time.Time              `json:"occured_at"`
			SubTransactions *[]model.SubTransaction `json:"sub_transactions"`
		}
		var updateTransaction UpdateTransaction
		err = c.BindJSON(&updateTransaction)

		if err != nil {
			c.String(http.StatusBadRequest, "Infection Informations")
			c.Abort()
			return
		}

		transactionId, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Transaction ID is invalid")
			c.Abort()
			return
		}

		transaction, err := crud.GetTransaction(&book, transactionId)
		if err != nil {
			c.String(http.StatusBadRequest, "Transaction ID is invalid")
		}

		if updateTransaction.Description != nil {
			transaction.Description = *updateTransaction.Description
		}
		if updateTransaction.OccuredAt != nil {
			transaction.OccuredAt = *updateTransaction.OccuredAt
		}
		if updateTransaction.SubTransactions != nil {
			for idx := range *updateTransaction.SubTransactions {
				(*updateTransaction.SubTransactions)[idx].BookId = book.BookId
			}
			transaction.SubTransactions = *updateTransaction.SubTransactions
		}

		err = crud.UpdateTransaction(&transaction)
		if err != nil {
			c.String(http.StatusInternalServerError, "Transactiuon could not created")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transaction": transaction,
			"message":     "Transaction was created",
		})
	}
}

func GetTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		book, err := crud.GetBook(c.Param("bid"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Book was not found")
			c.Abort()
			return
		}

		bookAuthorization, err := crud.GetBookAuthorization(&user, &book)
		if err != nil {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}
		if strings.Index(bookAuthorization.Authority, "read") == -1 {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}

		transactionId, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Transaction ID is invalid")
			c.Abort()
			return
		}

		transaction, err := crud.GetTransaction(&book, transactionId)
		if err != nil {
			c.String(http.StatusNotFound, "Transaction could not found")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transaction": transaction,
			"message":     "Transaction was found",
		})
	}
}

func GetTransactions(isPageNotZero bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		book, err := crud.GetBook(c.Param("bid"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Book was not found")
			c.Abort()
			return
		}

		bookAuthorization, err := crud.GetBookAuthorization(&user, &book)
		if err != nil {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}
		if strings.Index(bookAuthorization.Authority, "read") == -1 {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}

		page := 0
		if isPageNotZero {
			page, err = strconv.Atoi(c.Param("pid"))
		}

		transactions, err := crud.GetTransactions(&book, 20, page)
		if err != nil {
			c.String(http.StatusNotFound, "Transactions could not found")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transactions": transactions,
			"message":      "Transactions was found",
		})
	}
}

func DeleteTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		book, err := crud.GetBook(c.Param("bid"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Book was not found")
			c.Abort()
			return
		}

		bookAuthorization, err := crud.GetBookAuthorization(&user, &book)
		if err != nil {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}
		if strings.Index(bookAuthorization.Authority, "admin") == -1 {
			c.String(http.StatusUnauthorized, NoAuthorizationError.Error())
			c.Abort()
			return
		}

		transactionId, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Transaction ID is invalid")
			c.Abort()
			return
		}

		err = crud.DeleteTransaction(&book, transactionId)
		if err != nil {
			c.String(http.StatusNotFound, "Transaction could not delete")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction was deleted",
		})
	}
}
