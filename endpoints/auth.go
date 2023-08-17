package endpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Prokuma/PLAccounting-Backend/crud"
	model "github.com/Prokuma/PLAccounting-Backend/models"
	util "github.com/Prokuma/PLAccounting-Backend/utils"
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
