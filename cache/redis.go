package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"project_restapi/models"
	"time"
)

type Cache interface {
	GetClient() *redis.Client
	SetCache(key string, value interface{}, expired time.Duration) error
	GetCacheUser(key string) []*models.User
	GetCacheBook(key string) []*models.Book
	GetCacheReview(key string) []*models.BookReview
	DestroyCache(key ...string)
}

type cache struct {
	host     string
	password string
	db       int
}

func NewCache(host string, password string, db int) Cache {
	return &cache{
		host:     host,
		password: password,
		db:       db,
	}
}

func (c *cache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: c.password,
		DB:       c.db,
	})
}

func (c *cache) SetCache(key string, value interface{}, expired time.Duration) error {
	r := c.GetClient()
	defer r.Close()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.Set(key, json, expired).Err()

	return err
}

func (c *cache) GetCacheUser(key string) []*models.User {
	r := c.GetClient()
	defer r.Close()

	val, err := r.Get(key).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	var users []*models.User

	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		log.Println(err)
		return nil
	}

	return users
}

func (c *cache) GetCacheBook(key string) []*models.Book {
	r := c.GetClient()
	defer r.Close()

	val, err := r.Get(key).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	var book []*models.Book

	err = json.Unmarshal([]byte(val), &book)
	if err != nil {
		log.Println(err)
		return nil
	}

	return book
}

func (c *cache) GetCacheReview(key string) []*models.BookReview {
	r := c.GetClient()
	defer r.Close()

	val, err := r.Get(key).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	var review []*models.BookReview

	err = json.Unmarshal([]byte(val), &review)
	if err != nil {
		log.Println(err)
		return nil
	}

	return review
}

func (c *cache) DestroyCache(key ...string) {
	r := c.GetClient()
	defer r.Close()

	for _, i := range key {
		r.Del(i)
	}
}