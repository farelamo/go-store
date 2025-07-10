package repository

import (
	"context"
	"store/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	SaveRedis(authRedis model.AuthRedis) error
	GetByKey(key string) (string, error)
	Delete(key string) error
}

type authRepository struct {
	Rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) AuthRepository {
	return &authRepository{
		Rdb: rdb,
	}
}

func (a *authRepository) SaveRedis(auth model.AuthRedis) error {
	err := a.Rdb.Set(context.Background(), auth.Key, auth.Value, 30*time.Minute).Err()
	if err != nil {
	}
	return nil
}

func (a *authRepository) GetByKey(key string) (string, error) {
	value, err := a.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (a *authRepository) Delete(key string) error {
	err := a.Rdb.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}
