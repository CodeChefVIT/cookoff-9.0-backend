package auth

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
)

var TokenAuth *jwtauth.JWTAuth

func InitJWT() {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		panic("JWT_KEY environment variable not set")
	}
	TokenAuth = jwtauth.New("HS256", []byte(jwtKey), nil)
}

func GenerateJWT(user *db.User, isRefresh bool) (string, time.Time, error) {
	var expirationTime time.Time

	if !isRefresh {
		expirationTime = time.Now().Add(time.Hour / 4)
		_, tokenString, err := TokenAuth.Encode(map[string]interface{}{
			"username": user.Name,
			"user_id":  user.ID,
			"role":     user.Role,
			"type":     "access",
			"exp":      expirationTime.Unix(),
		})
		return tokenString, expirationTime, err
	}

	expirationTime = time.Now().Add(time.Hour*1 + time.Minute*30)
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     expirationTime.Unix(),
	})
	return tokenString, expirationTime, err
}
