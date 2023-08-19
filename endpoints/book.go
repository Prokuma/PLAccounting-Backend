package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Prokuma/PLAccounting-Backend/crud"
	model "github.com/Prokuma/PLAccounting-Backend/models"
	"github.com/gin-gonic/gin"
)

type CreateBookRequest struct {
	Name string `json:"name" binding:"required"`
	Year uint   `json:"year" binding:"required"`
}

// CreateBook godoc
// @Summary Create Book
// @Tags Book
// @Description Create Book
// @Accept  json
// @Produce  json
// @Param book body CreateBookRequest true "Create Book"
// @Success 200 {string} string	"Created Book"
// @Failure 400 {string} string	"Request is failed"
// @Router /book [post]
func CreateBook(c *gin.Context) {
	user, err := getUserIdFromJWT(c)
	if err != nil {
		c.Abort()
		return
	}
	var createBook CreateBookRequest
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

// GetBook godoc
// @Summary Get Book
// @Tags Book
// @Description Get Book
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Get Book"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid} [get]
func GetBook(c *gin.Context) {
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

type UpdateBookRequest struct {
	Name *string `json:"name"`
	Year *uint   `json:"year"`
}

// UpdateBook godoc
// @Summary Update Book
// @Tags Book
// @Description Update Book
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param book body UpdateBookRequest true "Update Book"
// @Success 200 {string} string	"Update Book"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid} [patch]
func UpdateBook(c *gin.Context) {
	user, err := getUserIdFromJWT(c)
	if err != nil {
		c.Abort()
		return
	}
	var updateBook UpdateBookRequest
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

// DeleteBook godoc
// @Summary Delete Book
// @Tags Book
// @Description Delete Book
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Delete Book"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid} [delete]
func DeleteBook(c *gin.Context) {
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

type CreateAccountTitleRequest struct {
	Name   string `json:"name" binding:"required"`
	Amount int64  `json:"amount"`
	Type   uint   `json:"type"`
}

// CreateAccountTitle godoc
// @Summary Create Account Title
// @Tags Account Title
// @Description Create Account Title
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param accountTitle body CreateAccountTitleRequest true "Create Account Title"
// @Success 200 {string} string	"Create Account Title"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/accountTitle [post]
func CreateAccountTitle(c *gin.Context) {
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

	var createAccountTitle CreateAccountTitleRequest
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

// GetAccountTitle godoc
// @Summary Get Account Title
// @Tags Account Title
// @Description Get Account Title
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Account Title ID"
// @Success 200 {string} string	"Get Account Title"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/accountTitle/{tid} [get]
func GetAccountTitle(c *gin.Context) {
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

type UpdateAccountTitleRequest struct {
	Name   *string `json:"name"`
	Amount *int64  `json:"amount"`
	Type   *uint   `json:"type"`
}

// UpdateAccountTitle godoc
// @Summary Update Account Title
// @Tags Account Title
// @Description Update Account Title
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Account Title ID"
// @Param accountTitle body UpdateAccountTitleRequest true "Update Account Title"
// @Success 200 {string} string	"Update Account Title"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/accountTitle/{tid} [patch]
func UpdateAccountTitle(c *gin.Context) {
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

	var updateAccountTitle UpdateAccountTitleRequest
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

// DeleteAccountTitle godoc
// @Summary Delete Account Title
// @Tags Account Title
// @Description Delete Account Title
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param tid path string true "Account Title ID"
// @Success 200 {string} string	"Delete Account Title"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/accountTitle/{tid} [delete]
func DeleteAccountTitle(c *gin.Context) {
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
