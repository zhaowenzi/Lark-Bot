package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

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
		var larkSubscriptionEventRequest LarkSubscriptionEventEncryptRequest
		c.BindJSON(&larkSubscriptionEventRequest)
		decodeMessage, _ := Decrypt(larkSubscriptionEventRequest.Encrypt, encryptKey)
		var larkSubscriptionEventDecryptedRequest LarkSubscriptionEventDecryptedRequest
		json.Unmarshal([]byte(decodeMessage), &larkSubscriptionEventDecryptedRequest)
		logrus.WithFields(logrus.Fields{
			"request": decodeMessage,
		}).Info("Received request from Lark")

		if larkSubscriptionEventDecryptedRequest.Type == UrlVerification {
			c.JSON(http.StatusOK, LarkUrlVerificationResponse{Challenge: larkSubscriptionEventDecryptedRequest.Challenge})
			return
		}

		c.JSON(http.StatusOK, nil)
	})
	r.Run(":8888")
}
