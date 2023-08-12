package storage

import (
	"auth/src/application/dto"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type gmailWithKeyPairRedisStorage struct {
	rdb        *redis.Client
	prefix     string
	expiration time.Duration
}

func NewRedisTemporaryStorage(address string, expiration time.Duration, prefix string) *gmailWithKeyPairRedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return &gmailWithKeyPairRedisStorage{
		rdb:        rdb,
		prefix:     prefix,
		expiration: expiration,
	}
}

func (stor gmailWithKeyPairRedisStorage) Create(user dto.GmailWithKeyPairDTO) error {
	err := stor.rdb.Set(ctx, user.Key+stor.prefix, user.Gmail, stor.expiration).Err()
	if err != nil {
		return fmt.Errorf("can not create pair for %s in redis: %v", user.Gmail, err)
	}

	return nil
}

func (stor gmailWithKeyPairRedisStorage) GetByUniqueKey(key string) (dto.GmailWithKeyPairDTO, error) {
	gmail, err := stor.rdb.Get(ctx, key+stor.prefix).Result()
	if err != nil {
		if err == redis.Nil {
			return dto.GmailWithKeyPairDTO{}, fmt.Errorf("key %s not found in redis", key)
		}

		return dto.GmailWithKeyPairDTO{}, fmt.Errorf("can not get by key %s in redis: %v", key, err)
	}

	return dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   key,
	}, nil
}

func (stor gmailWithKeyPairRedisStorage) Delete(user dto.GmailWithKeyPairDTO) error {
	gmail, err := stor.rdb.Get(ctx, user.Key+stor.prefix).Result()
	if err != nil {
		return fmt.Errorf("can not get %s key from redis: %v", gmail, err)
	}

	if gmail != user.Gmail {
		return fmt.Errorf("key not found in redis: %v", user.Key)
	}

	err = stor.rdb.Del(ctx, user.Key+stor.prefix).Err()
	if err != nil {
		return fmt.Errorf("can not delete pair for %s from redis: %v", gmail, err)
	}

	return nil
}
