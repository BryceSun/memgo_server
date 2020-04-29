package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
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
