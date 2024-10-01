package submission

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TokenManager struct {
	client *redis.Client
}

var Tokens *TokenManager

func Init(client *redis.Client) {
	Tokens = &TokenManager{client: client}
}

func (tm *TokenManager) AddToken(ctx context.Context, token string, subID string, testcaseid string) error {
	err := tm.client.Set(ctx, fmt.Sprintf("token:%s", token), fmt.Sprintf("%s:%s", subID, testcaseid), 0).Err()
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

func (tm *TokenManager) GetTokenMember(ctx context.Context, subID string) ([]string, error) {
	members, err := tm.client.SMembers(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Result()
	return members, err
}

func (tm *TokenManager) GetTokenCount(ctx context.Context, subID string) (int64, error) {
	tokenCount, err := tm.client.SCard(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Result()
	if err != nil {
		return 0, err
	}
	return tokenCount, nil
}

func (tm *TokenManager) DeleteToken(ctx context.Context, token string) error {
	subID, _, err := GetSubID(ctx, token)
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
