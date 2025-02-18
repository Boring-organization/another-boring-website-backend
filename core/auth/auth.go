package auth

import (
	"TestGoLandProject/global_const"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var jwtKey = []byte(globalConst.JwtSeed)

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func GetTokenFromGinContext(context *gin.Context) (*string, error) {
	bearerToken := context.Request.Header.Get("Authorization")

	if bearerToken != "" {
		return &strings.Split(bearerToken, " ")[1], nil
	}

	return nil, fmt.Errorf("can't find bearer token")

}

func CheckTokenValidity(token string) error {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		return fmt.Errorf("can't parse token: %w", err)
	}

	return nil

}

func GetUserIdByToken(token string) (*string, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("can't parse token: %w", err)
	}

	return &claims.UserId, nil

}

func GenerateJwt(userId string, liveTime time.Duration) (string, error) {
	expirationTime := time.Now().Add(liveTime)

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)

}
