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

type CreateTransactionRequest struct {
	Description     string                 `json:"description" binding:"required"`
	OccuredAt       time.Time              `json:"occured_at" binding:"required"`
	SubTransactions []model.SubTransaction `json:"sub_transactions" binding:"required"`
}

// CreateTransaction godoc
// @Summary Create Transaction
// @Tags Transaction
// @Description Create Transaction
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param transaction body CreateTransactionRequest true "Create Transaction"
// @Success 200 {string} string	"Created Transaction"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction [post]
func CreateTransaction(c *gin.Context) {
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

	var createTransaction CreateTransactionRequest
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

type UpdateTransactionRequest struct {
	Description     *string                 `json:"description"`
	OccuredAt       *time.Time              `json:"occured_at"`
	SubTransactions *[]model.SubTransaction `json:"sub_transactions"`
}

// UpdateTransaction godoc
// @Summary Update Transaction
// @Tags Transaction
// @Description Update Transaction
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Transaction ID"
// @Param transaction body UpdateTransactionRequest true "Update Transaction"
// @Success 200 {string} string	"Updated Transaction"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction/{tid} [patch]
func UpdateTransaction(c *gin.Context) {
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

	var updateTransaction UpdateTransactionRequest
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

// GetTransaction godoc
// @Summary Get Transaction
// @Tags Transaction
// @Description Get Transaction
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Transaction ID"
// @Success 200 {string} string	"Transaction was found"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction/{tid} [get]
func GetTransaction(c *gin.Context) {
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

// GetTransactions godoc
// @Summary Get Transactions
// @Tags Transaction
// @Description Get Transactions
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Transactions was found"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction [get]
func GetTransactions(c *gin.Context) {
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

	transactions, err := crud.GetTransactions(&book, 20, 0)
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

// GetTransactionsWithPage godoc
// @Summary Get Transactions with Page
// @Tags Transaction
// @Description Get Transactions with Page
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param pid path string true "Page ID"
// @Success 200 {string} string	"Transactions was found"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction/page/{pid} [get]
func GetTransactionsWithPage(c *gin.Context) {
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

	page, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.String(http.StatusBadRequest, "Page ID is invalid")
		c.Abort()
		return
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

// DeleteTransaction godoc
// @Summary Delete Transaction
// @Tags Transaction
// @Description Delete Transaction
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Transaction ID"
// @Success 200 {string} string	"Transaction was deleted"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/transaction/{tid} [delete]
func DeleteTransaction(c *gin.Context) {
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
