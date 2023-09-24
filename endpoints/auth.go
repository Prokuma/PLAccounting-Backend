package endpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Prokuma/PLAccounting-Backend/crud"
	model "github.com/Prokuma/PLAccounting-Backend/models"
	util "github.com/Prokuma/PLAccounting-Backend/utils"
	"github.com/gin-gonic/gin"
)

var NoAuthorizationError = errors.New("No Authorization")

// Ping godoc
// @Summary Ping
// @Tags Ping
// @Description Ping
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"pong"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser godoc
// @Summary Create User
// @Tags User
// @Description Create User
// @Accept  json
// @Produce  json
// @Param user body CreateUserRequest true "Create User"
// @Success 200 {string} string	"Created User"
// @Failure 400 {string} string	"Request is failed"
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var createUser CreateUserRequest
	err := c.BindJSON(&createUser)
	salt := os.Getenv("HASH_SALT")

	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed ")
		c.Abort()
		return
	}

	r := sha256.Sum256([]byte(createUser.Password + salt))
	hash := hex.EncodeToString(r[:])

	userInfo := model.User{
		Email:    createUser.Email,
		Name:     createUser.Name,
		Password: hash,
	}

	triedUser, err := crud.GetUserFromEmail(createUser.Email)
	if err == nil {
		if triedUser.Email == createUser.Email {
			c.String(http.StatusBadRequest, "Request was failed ")
			c.Abort()
			return
		}
	}

	token, err := util.GenerateMailConfirmationToken(&userInfo)
	if err != nil {
		c.String(http.StatusInternalServerError, "Making Token was failed")
		c.Abort()
		return
	}

	err = util.SendRealCreateUserMail(createUser.Email, token)
	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed ")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Created User Requested: %s", createUser.Name),
	})
}

// CreateUserAtDatabase godoc
// @Summary Create User At Database
// Note: This endpoint is not for API.
func CreateUserAtDatabase(c *gin.Context) {
	tokenQuery := c.Query("token")
	email, err := util.Redis.Get(util.Context, tokenQuery).Result()
	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed ")
		c.Abort()
		return
	}

	tokenInfo, err := util.Redis.Get(util.Context, email+".info").Result()
	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed ")
		c.Abort()
		return
	}

	var user model.User
	err = json.Unmarshal([]byte(tokenInfo), &user)
	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed ")
		c.Abort()
		return
	}

	err = crud.CreateUser(&user)
	if err != nil {
		c.String(http.StatusInternalServerError, "Create User was failed ")
		c.Abort()
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Created User"))
}

type LoginWithEmailAndPassword struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login
// @Tags User
// @Description Login
// @Accept  json
// @Produce  json
// @Param user body LoginWithEmailAndPassword true "Login"
// @Success 200 {string} string	"Login"
// @Failure 400 {string} string	"Request is failed"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginInfo LoginWithEmailAndPassword
	err := c.BindJSON(&loginInfo)
	salt := os.Getenv("HASH_SALT")

	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed "+err.Error())
		c.Abort()
		return
	}

	r := sha256.Sum256([]byte(loginInfo.Password + salt))
	hash := hex.EncodeToString(r[:])

	user, err := crud.GetUserFromEmail(loginInfo.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "Request was failed "+err.Error())
		c.Abort()
		return
	}

	if user.Password != hash {
		c.String(http.StatusForbidden, "Password was not currect")
		c.Abort()
		return
	}

	token, tokenLifeTime, err := util.GenerateToken(user.UserId)

	if err != nil {
		c.String(http.StatusInternalServerError, "Making Token was failed")
		c.Abort()
		return
	}

	c.SetCookie("token", token, tokenLifeTime*3600, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"email":   user.Email,
		"name":    user.Name,
		"expire":  tokenLifeTime * 3600,
		"message": fmt.Sprintf("Succesed Login: %s", user.Name),
	})
}

// GetUser godoc
// @Summary Get User
// @Tags User
// @Description Get User
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"Get User"
// @Failure 400 {string} string	"Request is failed"
// @Router /user [get]
func GetUser(c *gin.Context) {
	user, err := getUserIdFromJWT(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"id":     user.UserId,
		"email":  user.Email,
		"name":   user.Name,
	})
}

// Logout godoc
// @Summary Logout
// @Tags User
// @Description Logout
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"Logout"
// @Failure 400 {string} string	"Request is failed"
// @Router /logout [get]
func Logout(c *gin.Context) {
	user, err := getUserIdFromJWT(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Not Logged In")
		c.Abort()
		return
	}

	util.Redis.Del(util.Context, user.UserId)
	c.SetCookie("token", "", 0, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Succesed Logout",
	})
}

func getUserIdFromJWT(c *gin.Context) (model.User, error) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, err
	}

	token, err := util.ParseToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, err
	}

	claims := token.Claims.(*util.TokenClaims)
	userId := claims.UserId
	exp := claims.Exp

	ret, err := util.Redis.Get(util.Context, tokenString).Result()
	if err != nil {
		fmt.Println("Error: ", err)
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, err
	}

	retInt, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		fmt.Println("Error: ", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return model.User{}, err
	}
	if retInt < exp {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, NoAuthorizationError
	}

	user, err := crud.GetUser(userId)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return model.User{}, err
	}

	return user, nil
}
