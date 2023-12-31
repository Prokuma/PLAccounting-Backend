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
		c.String(http.StatusUnauthorized, "User was not found")
		c.Abort()
		return
	}

	book, err := crud.GetBook(c.Param("bid"))
	if err != nil {
		c.String(http.StatusNotFound, "Book was not found")
		c.Abort()
		return
	}

	pages := crud.GetBookPages(book.BookId, 20)

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
		"pages":   pages,
		"message": "Get Book Successed",
	})
}

// CreateBookFromOldBook godoc
// @Summary Create Book From Old Book
// @Tags Book
// @Description Create Book From Old Book
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Create Book From Old Book"
// @Failure 400 {string} string	"Request is failed"
// @Router /migrate/{bid} [post]
func CreateBookFromOldBook(c *gin.Context) {
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

	oldBook, err := crud.GetBook(c.Param("bid"))
	if err != nil {
		c.String(http.StatusNotFound, "Book was not found")
		c.Abort()
		return
	}

	err = crud.CreateBookAndAccountTitleFromBook(createBook.Year, createBook.Name, &user, &oldBook)
	if err != nil {
		c.String(http.StatusInternalServerError, "Book could not created")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book was created",
	})
}

// GetAllBooks godoc
// @Summary Get All Books
// @Tags Book
// @Description Get All Books
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"Get All Books"
// @Failure 400 {string} string	"Request is failed"
// @Router /book [get]
func GetAllBooks(c *gin.Context) {
	user, err := getUserIdFromJWT(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "User was not found")
		c.Abort()
		return
	}

	books, err := crud.GetAllBooks(&user)
	if err != nil {
		c.String(http.StatusNotFound, "No Book")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books":   books,
		"message": "Books was found",
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

type CreateBookAuthorizationRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	Authority string `json:"authority" binding:"required"`
}

// CreateBookAuthorization godoc
// @Summary Create Book Authorization
// @Tags Book Authorization
// @Description Create Book Authorization
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Create Book Authorization"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/bookAuthorization [post]
func CreateBookAuthorization(c *gin.Context) {
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
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}
	if strings.Index(bookAuthorization.Authority, "admin") == -1 {
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}

	var createBookAuthorization CreateBookAuthorizationRequest
	err = c.BindJSON(&createBookAuthorization)

	if err != nil {
		c.String(http.StatusBadRequest, "Infection Informations")
		c.Abort()
		return
	}

	var inputBookAuthorization = model.BookAuthorization{
		BookId:    book.BookId,
		UserId:    createBookAuthorization.UserId,
		Authority: createBookAuthorization.Authority,
	}
	err = crud.CreateBookAuthorization(&inputBookAuthorization)
	if err != nil {
		c.String(http.StatusInternalServerError, "Book Authorization could not created")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book_authorization": inputBookAuthorization,
		"message":            "Book Authorization was created",
	})
}

type UpdateBookAuthorizationRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	Authority string `json:"authority" binding:"required"`
}

// UpdateBookAuthorization godoc
// @Summary Update Book Authorization
// @Tags Book Authorization
// @Description Update Book Authorization
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Param uid path string true "User ID"
// @Success 200 {string} string	"Update Book Authorization"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/bookAuthorization/{uid} [patch]
func UpdateBookAuthorization(c *gin.Context) {
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
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}
	if strings.Index(bookAuthorization.Authority, "admin") == -1 {
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}

	var updateBookAuthorization UpdateBookAuthorizationRequest
	err = c.BindJSON(&updateBookAuthorization)
	if err != nil {
		c.String(http.StatusBadRequest, "Infection Informations")
		c.Abort()
		return
	}

	var inputBookAuthorization = model.BookAuthorization{
		BookId:    book.BookId,
		UserId:    updateBookAuthorization.UserId,
		Authority: updateBookAuthorization.Authority,
	}

	err = crud.UpdateBookAuthorization(&inputBookAuthorization)
	if err != nil {
		c.String(http.StatusInternalServerError, "Book Authorization could not updated")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book_authorization": inputBookAuthorization,
		"message":            "Book Authorization was updated",
	})
}

func GetBookAuthorizations(c *gin.Context) {
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
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}
	if strings.Index(bookAuthorization.Authority, "admin") == -1 {
		c.String(http.StatusUnauthorized, "You are not authorized")
		c.Abort()
		return
	}

	bookAuthorizations, err := crud.GetBookAuthorizations(&book)
	if err != nil {
		c.String(http.StatusInternalServerError, "Book Authorization could not updated")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book_authorizations": bookAuthorizations,
		"message":             "Book Authorization was updated",
	})
}

type CreateAccountTitleRequest struct {
	Name       string `json:"name" binding:"required"`
	Amount     int64  `json:"amount"`
	AmountBase int64  `json:"amount_base"`
	Type       uint   `json:"type"`
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
		BookId:     book.BookId,
		Name:       createAccountTitle.Name,
		Amount:     createAccountTitle.Amount,
		AmountBase: createAccountTitle.AmountBase,
		Type:       createAccountTitle.Type,
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

// GetAllAccountTitles godoc
// @Summary Get All Account Titles
// @Tags Account Title
// @Description Get All Account Titles
// @Accept  json
// @Produce  json
// @Param bid path string true "Book ID"
// @Success 200 {string} string	"Get All Account Titles"
// @Failure 400 {string} string	"Request is failed"
// @Router /book/{bid}/accountTitle [get]
func GetAllAccountTitles(c *gin.Context) {
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

	accountTitles, err := crud.GetAllAccountTitles(&book)
	if err != nil {
		c.String(http.StatusNotFound, "No Account Title")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_titles": accountTitles,
		"message":        "Account Titles was found",
	})
}

type UpdateAccountTitleRequest struct {
	Name       *string `json:"name"`
	Amount     *int64  `json:"amount"`
	AmountBase *int64  `json:"amount_base"`
	Type       *uint   `json:"type"`
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
	if updateAccountTitle.AmountBase != nil {
		accountTitle.AmountBase = *updateAccountTitle.AmountBase
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
