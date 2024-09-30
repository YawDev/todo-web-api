package authentication

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("Secret_Key")

type Claims struct {
	Username string
	UserID   int
	jwt.RegisteredClaims
}

func GenerateAccessToken(username string, userId int) (string, error) {

	claims := &Claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(15 * time.Minute)},
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("error while creating access token")
	}
	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {

	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("error while creating refresh token")
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil && errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, errors.New("access token signature invalid")
	} else if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func ParseRefreshToken(tokenStr string) (*jwt.Claims, error) {
	var claims jwt.Claims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil && errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, errors.New("refresh token signature invalid")
	} else if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return &claims, nil
}

func payload(claims *Claims, c *gin.Context) {
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
}
