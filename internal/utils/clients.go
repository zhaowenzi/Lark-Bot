package utils

import (
	"github.com/joho/godotenv"
	"github.com/larksuite/oapi-sdk-go/v3"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var openaiClient *openai.Client

var larkClient *lark.Client

var cache = cache2.New(10*time.Hour, 12*time.Hour)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
		return
	}
	openaiClient = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	larkClient = lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), lark.WithEnableTokenCache(true))
}

func GetLarkClient() *lark.Client {
	return larkClient
}
