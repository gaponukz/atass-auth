package storage

import (
	"auth/src/entities"
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

func RedisTemporaryStorage(expiration time.Duration, prefix string) gmailWithKeyPairRedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return gmailWithKeyPairRedisStorage{
		rdb:        rdb,
		prefix:     prefix,
		expiration: expiration,
	}
}

func (stor gmailWithKeyPairRedisStorage) Create(user entities.GmailWithKeyPair) error {
	return stor.rdb.Set(ctx, user.Key+stor.prefix, user.Gmail, stor.expiration).Err()
}

func (stor gmailWithKeyPairRedisStorage) GetByUniqueKey(key string) (entities.GmailWithKeyPair, error) {
	gmail, err := stor.rdb.Get(ctx, key+stor.prefix).Result()

	if err != nil {
		if err == redis.Nil {
			return entities.GmailWithKeyPair{}, fmt.Errorf("key not found: %v", key)
		}

		return entities.GmailWithKeyPair{}, err
	}

	return entities.GmailWithKeyPair{
		Gmail: gmail,
		Key:   key,
	}, nil
}

func (stor gmailWithKeyPairRedisStorage) Delete(user entities.GmailWithKeyPair) error {
	return stor.rdb.Del(ctx, user.Key+stor.prefix).Err()
}
