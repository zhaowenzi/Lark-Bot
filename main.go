package main

import (
	"Lark-Bot/internal/constants"
	"Lark-Bot/internal/handlers"
	"Lark-Bot/internal/structs"
	"Lark-Bot/internal/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
		return
	}
	encryptKey := os.Getenv("ENCRYPT_KEY")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/lark", func(c *gin.Context) {
		var larkSubscriptionEventRequest structs.LarkSubscriptionEventEncryptRequest
		c.BindJSON(&larkSubscriptionEventRequest)
		decodeMessage, _ := utils.Decrypt(*larkSubscriptionEventRequest.Encrypt, encryptKey)
		var larkSubscriptionEventDecryptedRequest structs.LarkSubscriptionEventDecryptedRequest
		json.Unmarshal([]byte(decodeMessage), &larkSubscriptionEventDecryptedRequest)
		logrus.WithFields(logrus.Fields{
			"request": decodeMessage,
		}).Info("Received request from Lark")

		if larkSubscriptionEventDecryptedRequest.Type != nil && *larkSubscriptionEventDecryptedRequest.Type == constants.UrlVerification {
			c.JSON(http.StatusOK, structs.LarkUrlVerificationResponse{Challenge: *larkSubscriptionEventDecryptedRequest.Challenge})
			return
		}

		switch *larkSubscriptionEventDecryptedRequest.Header.EventType {
		case constants.MessageReceive:
			inCache := utils.CheckCache(*larkSubscriptionEventDecryptedRequest.Event.Message.MessageId)
			if inCache {
				c.Status(http.StatusOK)
				return
			}
			handlers.HandleReceivedMessage(larkSubscriptionEventDecryptedRequest)
		}
		c.Status(http.StatusOK)
	})
	r.POST("/msg", func(c *gin.Context) {
		var textMessageRequest structs.TextMessageRequest
		c.BindJSON(&textMessageRequest)
		utils.GetLarkRedisClient().RPush(ctx, "msg", *textMessageRequest.Text)
		c.Status(http.StatusOK)
	})
	r.GET("/paste", func(c *gin.Context) {
		result, _ := utils.GetPasteRedisClient().LRange(ctx, time.Now().Format("2006-01-02"), -10, -1).Result()
		c.JSON(http.StatusOK, gin.H{
			"paste": result,
		})
	})
	r.POST("/paste", func(c *gin.Context) {
		var pastePostRequest structs.PastePostRequest
		c.BindJSON(&pastePostRequest)
		utils.GetPasteRedisClient().RPush(ctx, time.Now().Format("2006-01-02"), *pastePostRequest.PasteContent)
		c.Status(http.StatusOK)
	})
	r.Run(os.Getenv("GIN_ADDRESS"))
}
