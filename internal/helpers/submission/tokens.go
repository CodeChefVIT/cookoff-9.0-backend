package submission

import (
	"context"
	"fmt"
	"os"
	"time"

	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/redis/go-redis/v9"
)

type TokenManager struct {
	client *redis.Client
}

var Tokens *TokenManager

func Init() {
	host := os.Getenv("DRAGONFLY_HOST")
	port := os.Getenv("DRAGONFLY_PORT")
	pwd := os.Getenv("DRAGONFLY_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		DB:           0,
		Password:     pwd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     50,
		PoolTimeout:  10 * time.Second,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Errof("Redis Init Failed: " + err.Error())
	}
	logger.Infof("Connect to dragonflydb")
	Tokens = &TokenManager{client: client}
}

func (tm *TokenManager) AddToken(ctx context.Context, token string, subID string) error {
	err := tm.client.Set(ctx, fmt.Sprintf("token:%s", token), subID, 0).Err()
	if err != nil {
		return err
	}

	err = tm.client.SAdd(ctx, fmt.Sprintf("sub:%s:tokens", subID), token).Err()
	if err != nil {
		return err
	}
	return nil
}

func (tm *TokenManager) GetSubID(ctx context.Context, token string) (string, error) {
	subID, err := tm.client.Get(ctx, fmt.Sprintf("token:%s", token)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found")
	} else if err != nil {
		return "", err
	}
	return subID, nil
}

func (tm *TokenManager) DeleteToken(ctx context.Context, token string) error {
	subID, err := tm.GetSubID(ctx, token)
	if err != nil {
		return err
	}

	err = tm.client.Del(ctx, fmt.Sprintf("token:%s", token)).Err()
	if err != nil {
		return err
	}

	err = tm.client.SRem(ctx, fmt.Sprintf("sub:%s:tokens", subID), token).Err()
	if err != nil {
		return err
	}

	setSize, err := tm.client.SCard(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Result()
	if err != nil {
		return err
	}

	if setSize == 0 {
		err = tm.client.Del(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
