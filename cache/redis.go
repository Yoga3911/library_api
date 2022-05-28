package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisC interface {
	GetClient() *redis.Client
	Set(key string, value interface{})
	Get(key string) string
	Del(key string)
}

type redisC struct {
	host     string
	password string
	db       int
}

func NewRedisC(host string, pass string, db int) RedisC {
	return &redisC{
		host:     host,
		password: pass,
		db:       db,
	}
}

func (r *redisC) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.host,
		Password: r.password,
		DB:       r.db,
	})
}

func (r *redisC) Set(key string, value interface{}) {
	client := r.GetClient()
	defer client.Close()

	enc, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
	}

	err = client.Set(key, enc, time.Minute * 10).Err()
	if err != nil {
		log.Println(err)
	}
}

func (r *redisC) Get(key string) string {
	client := r.GetClient()
	defer client.Close()

	val, _ := client.Get(key).Result()
	var email string

	if err := json.Unmarshal([]byte(val), &email); err != nil {
		log.Println(err)
	}
	return email
}

func (r *redisC) Del(key string) {
	client := r.GetClient()
	defer client.Close()

	client.Del(key)
}
