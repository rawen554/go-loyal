package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

const tokenExp = time.Hour * 3
const CookieName = "jwt-token"
const UserIDKey = "userID"

var ErrTokenNotValid = errors.New("token is not valid")
var ErrNoUserInToken = errors.New("no user data in token")

func BuildJWTString(userID uint64, seed string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(seed))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string, seed string) (uint64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(seed), nil
		})
	if err != nil {
		if !token.Valid {
			return 0, ErrTokenNotValid
		} else {
			return 0, errors.New("parsing error")
		}
	}

	if claims.UserID == 0 {
		return 0, ErrNoUserInToken
	}

	return claims.UserID, nil
}

func AuthMiddleware(seed string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(CookieName)
		if err != nil {
			log.Printf("Error reading cookie[%v]: %v", CookieName, err)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}

		userID, err := GetUserID(cookie, seed)
		if err != nil {
			if errors.Is(err, ErrNoUserInToken) || errors.Is(err, ErrTokenNotValid) {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Abort()
				return
			} else {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				c.Abort()
				return
			}
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
