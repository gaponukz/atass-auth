package storage

import (
	"auth/src/entities"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type FutureUserRedisStorage struct {
	rdb        *redis.Client
	expiration time.Duration
}

func NewFutureUserRedisStorage(expiration time.Duration) IFutureUserStorage {
	return FutureUserRedisStorage{
		rdb: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		expiration: expiration,
	}
}

func (stor FutureUserRedisStorage) Create(user entities.FutureUser) error {
	stringJson, err := futureUserToJson(user)

	if err != nil {
		return err
	}

	return stor.rdb.Set(ctx, user.UniqueKey, stringJson, stor.expiration).Err()
}

func (stor FutureUserRedisStorage) GetByUniqueKey(key string) (entities.FutureUser, error) {
	stringJson, err := stor.rdb.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return entities.FutureUser{}, fmt.Errorf("key not found: %v", key)
		}

		return entities.FutureUser{}, err
	}

	return futureUserFromJson(stringJson)
}

func (stor FutureUserRedisStorage) Delete(user entities.FutureUser) error {
	return stor.rdb.Del(ctx, user.UniqueKey).Err()
}

func futureUserToJson(user entities.FutureUser) (string, error) {
	data, err := json.MarshalIndent(user, "", " ")

	return string(data), err
}

func futureUserFromJson(stringJson string) (entities.FutureUser, error) {
	var user entities.FutureUser

	err := json.Unmarshal([]byte(stringJson), &user)

	return user, err
}
