package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var Context context.Context
var Redis *redis.Client

type TokenClaims struct {
	UserId string `json:"user_id"`
	Exp    int64  `json:"exp"`
	jwt.RegisteredClaims
}

func InitRedis() {
	Context = context.Background()
	// Redis Initilization
	hostAddr := os.Getenv("REDIS_HOST")
	hostPort := os.Getenv("REDIS_PORT")
	hostPassword := os.Getenv("REDIS_PASSWORD")
	hostDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     hostAddr + ":" + hostPort,
		Password: hostPassword,
		DB:       hostDB,
		PoolSize: 1000,
	})
}

func GenerateToken(userId string) (string, int, error) {
	secretKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	tokenLifeTime, err := strconv.Atoi(os.Getenv("JWT_TOKEN_LIFETIME"))

	secretKeyFile, err := os.ReadFile(secretKeyPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", 0, err
	}

	secretKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretKeyFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", 0, err
	}

	exp := time.Now().Add(time.Hour * time.Duration(tokenLifeTime)).Unix()
	claims := TokenClaims{
		UserId: userId,
		Exp:    exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", 0, err
	}

	Redis.Set(Context, userId, strconv.FormatInt(exp, 10), 0)

	return tokenString, tokenLifeTime, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	verifyKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	verifyKeyFile, err := os.ReadFile(verifyKeyPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyKeyFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if _, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, errors.New("Invalid Token")
	}
}
