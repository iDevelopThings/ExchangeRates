package main

import (
	"encoding/json"
	"time"
	
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})
}

func Set(key string, value interface{}) error {
	p, err := json.Marshal(value)
	
	if err != nil {
		return err
	}
	
	set := Redis.Set(key, p, time.Hour*6)
	
	if set.Err() != nil {
		return err
	}
	
	return nil
}

func Get(key string, dest interface{}) error {
	p, err := Redis.Get(key).Result()
	
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(p), dest)
}
