package utils

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/larksuite/oapi-sdk-go/v3"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()

var openaiClient *openai.Client

var larkClient *lark.Client

var cache = cache2.New(10*time.Hour, 12*time.Hour)

var larkRedisClient *redis.Client

var pasteRedisClient *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
		return
	}
	openaiClient = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	larkClient = lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), lark.WithEnableTokenCache(true))

	larkDBInt, err := strconv.Atoi(os.Getenv("LARK_DB"))
	if err != nil {
		panic(err)
	}
	pasteDBInt, err := strconv.Atoi(os.Getenv("PASTE_DB"))
	if err != nil {
		panic(err)
	}
	larkRedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       larkDBInt,
	})
	pasteRedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       pasteDBInt,
	})
}

func GetLarkClient() *lark.Client {
	return larkClient
}

func GetLarkRedisClient() *redis.Client {
	return larkRedisClient
}

func GetPasteRedisClient() *redis.Client {
	return pasteRedisClient
}
