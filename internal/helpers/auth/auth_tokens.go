package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type AuthTokenManager struct {
	client        *redis.Client
	AccessSecret  string
	RefreshSecret string
}

var Tokens *AuthTokenManager

func Init(client *redis.Client) {
	Tokens = &AuthTokenManager{
		client:        client,
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
	}
}

func (tm *AuthTokenManager) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(tm.AccessSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (tm *AuthTokenManager) GetUserID(ctx context.Context, accessToken string) (string, error) {
	userID, err := tm.client.Get(ctx, fmt.Sprintf("access_token:%s", accessToken)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("access token not found")
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

func (tm *AuthTokenManager) VerifyRefreshToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.RefreshSecret), nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("expiration claim is not a valid float64")
	}

	if time.Now().Unix() > int64(exp) {
		return false, fmt.Errorf("token has expired")
	}

	_, err = tm.client.Get(context.Background(), fmt.Sprintf("refresh_token:%s", tokenString)).
		Result()
	if err == redis.Nil {
		return false, fmt.Errorf("refresh token expired or not found")
	} else if err != nil {
		return false, fmt.Errorf("error retrieving refresh token: %v", err)
	}
	return true, nil
}
