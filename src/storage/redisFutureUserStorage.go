package storage

import (
	"auth/src/dto"
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
	return stor.rdb.Set(ctx, user.Key+stor.prefix, user.Gmail, stor.expiration).Err()
}

func (stor gmailWithKeyPairRedisStorage) GetByUniqueKey(key string) (dto.GmailWithKeyPairDTO, error) {
	gmail, err := stor.rdb.Get(ctx, key+stor.prefix).Result()

	if err != nil {
		if err == redis.Nil {
			return dto.GmailWithKeyPairDTO{}, fmt.Errorf("key not found: %v", key)
		}

		return dto.GmailWithKeyPairDTO{}, err
	}

	return dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   key,
	}, nil
}

func (stor gmailWithKeyPairRedisStorage) Delete(user dto.GmailWithKeyPairDTO) error {
	return stor.rdb.Del(ctx, user.Key+stor.prefix).Err()
}
