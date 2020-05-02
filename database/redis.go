package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

const (
	memgoKey = "memgo"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redisClient()
}

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

func clearRedis(pattern string) (int64, error) {
	keys, e := RedisClient.Keys(pattern).Result()
	if e != nil {
		return 0, e
	}
	if len(keys) == 0 {
		return 0, nil
	}
	return RedisClient.Del(keys...).Result()

}
