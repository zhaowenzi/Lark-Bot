package utils

import (
	"Lark-Bot/internal/constants"
	"Lark-Bot/internal/structs"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)

func GetUUID() string {
	return uuid.New().String()
}

func Decrypt(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}

func CallOpenAI(request []structs.RedisMessage) string {
	var chatCompletionMessages []openai.ChatCompletionMessage
	for _, message := range request {
		if *message.SenderType == constants.SenderTypeUser {
			chatCompletionMessages = append(chatCompletionMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: gjson.Get(*message.Content, "text").String(),
			})
		} else {
			chatCompletionMessages = append(chatCompletionMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: gjson.Get(*message.Content, "text").String(),
			})
		}
	}

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: chatCompletionMessages,
		})

	if err != nil {
		logrus.Warnf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func CheckCache(key string) bool {
	_, found := cache.Get(key)
	if found {
		return true
	} else {
		cache.SetDefault(key, 1)
		return false
	}
}

func GetAllParentMessages(key string) (messages []structs.RedisMessage) {
	var slice *redis.StringSliceCmd
	slice = GetLarkRedisClient().LRange(ctx, key, 0, -1)

	for index := range slice.Val() {
		var eachDecoded structs.RedisMessage
		json.Unmarshal([]byte(slice.Val()[index]), &eachDecoded)
		messages = append(messages, eachDecoded)
	}
	return messages
}

func StoreMessage(key string, redisMessage structs.RedisMessage) {
	encoded, _ := json.Marshal(redisMessage)
	GetLarkRedisClient().RPush(ctx, key, encoded)
}
