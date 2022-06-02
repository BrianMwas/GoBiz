package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

var RI ICache

type ICache interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
	LPush(key string, data ...interface{}) error
	LRange(key string, start int64, end int64) ([]string, error)
	Scan(key string) []string
}

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, b, expiration).Err()
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	result, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	return []byte(result), err
}

func (r *RedisCache) LPush(key string, data ...interface{}) error {
	err := r.client.LPush(context.Background(), key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) LRange(key string, start int64, end int64) ([]string, error) {
	result, err := r.client.LRange(context.Background(), key, start, end).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RedisCache) Scan(key string) []string {
	var values []string
	ctx := context.Background()
	iter := r.client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		values = append(values, iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}

	return values
}

func InitRedisCache() *RedisCache {
	var red = &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

	err := ping(red.client)
	if err != nil {
		return nil
	}

	return red
}

func ping(client *redis.Client) error {
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	logrus.Info("Connected to redis ", pong)
	return nil
}
