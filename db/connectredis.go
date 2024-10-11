package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var ConnectRedisDB *redis.Client; 

func ConnectRedis() *redis.Client{

    redisPassword, _ := os.LookupEnv("REDIS_PASSWORD")
    host, _ := os.LookupEnv("REDIS_HOST")
    port, _ := os.LookupEnv("REDIS_PORT")


    client := redis.NewClient(&redis.Options{
        Addr:	  host+":"+port,
        Password: redisPassword,
        DB:		  0,
    })
    ConnectRedisDB=client
	return client
}