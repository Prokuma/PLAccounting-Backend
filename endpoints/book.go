package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Prokuma/ProkumaLabAccount-Backend/crud"
	model "github.com/Prokuma/ProkumaLabAccount-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		type CreateBook struct {
			Name string `json:"name" binding:"required"`
			Year uint   `json:"year" binding:"required"`
		}
		var createBook CreateBook
		err = c.BindJSON(&createBook)

		if err != nil {
			c.String(http.StatusBadRequest, "Infection Informations")
			c.Abort()
			return
		}

		var book = model.Book{
			Name: createBook.Name,
			Year: createBook.Year,
		}

		err = crud.CreateBook(&user, &book)
		if err != nil {
			c.String(http.StatusInternalServerError, "Book could not created")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"book":    book,
			"message": "Book was created",
		})
	}
}

func GetBook() gin.HandlerFunc {
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

		c.JSON(http.StatusOK, gin.H{
			"book":    book,
			"message": "Created Book Successed",
		})
	}
}

func UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserIdFromJWT(c)
		if err != nil {
			c.Abort()
			return
		}

		type UpdateBook struct {
			Name *string `json:"name"`
			Year *uint   `json:"year"`
		}
		var updateBook UpdateBook
		err = c.BindJSON(&updateBook)

		if err != nil {
			c.String(http.StatusBadRequest, "Infection Informations")
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

		if updateBook.Name != nil {
			book.Name = *updateBook.Name
		}
		if updateBook.Year != nil {
			book.Year = *updateBook.Year
		}
		err = crud.UpdateBook(&book)
		if err != nil {
			c.String(http.StatusInternalServerError, "Book could not created")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"book":    book,
			"message": "Book was created",
		})
	}
}

func DeleteBook() gin.HandlerFunc {
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

		err = crud.DeleteBook(book.BookId)

		if err != nil {
			c.String(http.StatusInternalServerError, "Delete the book was failed")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Book was deleted",
		})
	}
}

func CreateAccountTitle() gin.HandlerFunc {
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

		type CreateAccountTitle struct {
			Name   string `json:"name" binding:"required"`
			Amount int64  `json:"amount"`
			Type   uint   `json:"type"`
		}
		var createAccountTitle CreateAccountTitle
		err = c.BindJSON(&createAccountTitle)

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "Infection Informations")
			c.Abort()
			return
		}

		var accountTitle = model.AccountTitle{
			BookId: book.BookId,
			Name:   createAccountTitle.Name,
			Amount: createAccountTitle.Amount,
			Type:   createAccountTitle.Type,
		}

		err = crud.CreateAccountTitle(&accountTitle)
		if err != nil {
			c.String(http.StatusInternalServerError, "Create Account Title was failed")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"account_title": accountTitle,
			"message":       "Account Title was created",
		})
	}
}

func GetAccountTitle() gin.HandlerFunc {
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

		tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "tid could not convert to integer")
			c.Abort()
			return
		}

		accountTitle, err := crud.GetAccountTitle(&book, tid)
		if err != nil {
			c.String(http.StatusNotFound, "No Account Title")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"account_title": accountTitle,
			"message":       "Account Title was found",
		})
	}
}

func UpdateAccountTitle() gin.HandlerFunc {
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

		type UpdateAccountTitle struct {
			Name   *string `json:"name"`
			Amount *int64  `json:"amount"`
			Type   *uint   `json:"type"`
		}
		var updateAccountTitle UpdateAccountTitle
		err = c.BindJSON(&updateAccountTitle)

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "Infection Informations")
			c.Abort()
			return
		}

		tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "tid could not convert to integer")
			c.Abort()
			return
		}

		accountTitle, err := crud.GetAccountTitle(&book, tid)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid title id")
			c.Abort()
			return
		}

		if updateAccountTitle.Name != nil {
			accountTitle.Name = *updateAccountTitle.Name
		}
		if updateAccountTitle.Amount != nil {
			accountTitle.Amount = *updateAccountTitle.Amount
		}
		if updateAccountTitle.Type != nil {
			accountTitle.Type = *updateAccountTitle.Type
		}

		err = crud.UpdateAccountTitle(&accountTitle)
		if err != nil {
			c.String(http.StatusNotFound, "The account title could not deleted")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"account_title": accountTitle,
			"message":       "Account Title was updated",
		})
	}
}

func DeleteAccountTitle() gin.HandlerFunc {
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

		tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "tid could not convert to integer")
			c.Abort()
			return
		}

		err = crud.DeleteAccountTitle(&book, tid)
		if err != nil {
			c.String(http.StatusNotFound, "The account title could not deleted")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Account Title was deleted",
		})
	}
}
