package endpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Prokuma/ProkumaLabAccount-Backend/crud"
	model "github.com/Prokuma/ProkumaLabAccount-Backend/models"
	util "github.com/Prokuma/ProkumaLabAccount-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var NoAuthorizationError = errors.New("No Authorization")

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		type CreateUser struct {
			Email    string `json:"email" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var createUser CreateUser
		err := c.BindJSON(&createUser)
		salt := os.Getenv("HASH_SALT")

		if err != nil {
			c.String(http.StatusBadRequest, "Request is failed "+err.Error())
			return
		}

		r := sha256.Sum256([]byte(createUser.Password + salt))
		hash := hex.EncodeToString(r[:])

		err = crud.CreateUser(&model.User{
			Email:    createUser.Email,
			Name:     createUser.Name,
			Password: hash,
		})

		if err != nil {
			c.String(http.StatusBadRequest, "Request is failed "+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Created User: %s", createUser.Name),
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		type LoginWithEmailAndPassword struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var loginInfo LoginWithEmailAndPassword
		err := c.BindJSON(&loginInfo)
		salt := os.Getenv("HASH_SALT")

		if err != nil {
			c.String(http.StatusBadRequest, "Request is failed "+err.Error())
			return
		}

		r := sha256.Sum256([]byte(loginInfo.Password + salt))
		hash := hex.EncodeToString(r[:])

		user, err := crud.GetUserFromEmail(loginInfo.Email)

		if err != nil {
			c.String(http.StatusBadRequest, "Request is failed "+err.Error())
			return
		}

		if user.Password != hash {
			c.String(http.StatusForbidden, "Password is not currect")
			return
		}

		token, tokenLifeTime, err := util.GenerateToken(user.UserId)

		if err != nil {
			c.String(http.StatusInternalServerError, "Making Token was failed")
			return
		}

		c.SetCookie("token", token, tokenLifeTime*3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Succesed Login: %s", user.Name),
		})
	}
}

func getUserIdFromJWT(c *gin.Context) (model.User, error) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, err
	}

	token, err := util.ParseToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid Token: "+err.Error())
		return model.User{}, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)
	fmt.Println(userId)

	user, err := crud.GetUser(userId)
	if err != nil {
		c.String(http.StatusUnauthorized, "User not found")
		return model.User{}, err
	}
	return user, nil
}

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
			"book_id": book.BookId,
			"message": "Created Book Successed",
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
			"message": "The book was deleted",
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

		err = crud.CreateAccountTitle(&model.AccountTitle{
			BookId: book.BookId,
			Name:   createAccountTitle.Name,
			Amount: createAccountTitle.Amount,
			Type:   createAccountTitle.Type,
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Create Account Title was failed")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Account Title was created",
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
