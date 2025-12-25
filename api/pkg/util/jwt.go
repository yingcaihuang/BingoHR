package util

import (
	"hr-api/pkg/app"
	"hr-api/pkg/e"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

type Claims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(id int, username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour)

	claims := Claims{
		id,
		username,
		EncodeMD5(password),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "bingo-hr",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// GetcurrentUI
func GetCurrentUid(c *gin.Context) int {
	appG := app.Gin{C: c}

	if val, exists := c.Get("uid"); exists {
		if uid, ok := val.(int); ok {
			return uid
		}
	}

	appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
	return 0
}
