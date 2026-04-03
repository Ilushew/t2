package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CodeService struct {
	redis *redis.Client
}

func NewCodeService(redisClient *redis.Client) *CodeService {
	return &CodeService{
		redis: redisClient,
	}
}

func (s *CodeService) SetCode(ctx context.Context, userID string, code string) error {
	key := fmt.Sprintf("verify_code:%s", userID)
	return s.redis.Set(ctx, key, code, 15*time.Minute).Err()
}

func (s *CodeService) VerifyCode(ctx context.Context, userID string, code string) error{
	key := fmt.Sprintf("verify_code:%s", userID)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil{
		return fmt.Errorf("код не найден или истёк")
	}
	if err != nil{
		return err
	}
	if storedCode != code{
		return fmt.Errorf("неверный код")
	}
	s.redis.Del(ctx, key)
	return nil
}