package handlers

import (
	"Lark-Bot/internal/constants"
	"Lark-Bot/internal/structs"
	"Lark-Bot/internal/utils"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
	"strings"
)

var ctx = context.Background()

// HandleReceivedMessage The steps after receiving message:
//  1. Store to Redis
//  2. Get all parent messages
//  3. Parse the first message
//  4. Check the prefix of the first message
func HandleReceivedMessage(larkSubscriptionEventDecryptedRequest structs.LarkSubscriptionEventDecryptedRequest) {
	// Store To Redis
	redisMessage := structs.RedisMessage{
		MessageId:   larkSubscriptionEventDecryptedRequest.Event.Message.MessageId,
		RootId:      larkSubscriptionEventDecryptedRequest.Event.Message.RootId,
		ParentId:    larkSubscriptionEventDecryptedRequest.Event.Message.ParentId,
		CreateTime:  larkSubscriptionEventDecryptedRequest.Event.Message.CreateTime,
		ChatId:      larkSubscriptionEventDecryptedRequest.Event.Message.ChatId,
		ChatType:    larkSubscriptionEventDecryptedRequest.Event.Message.ChatType,
		MessageType: larkSubscriptionEventDecryptedRequest.Event.Message.MessageType,
		Content:     larkSubscriptionEventDecryptedRequest.Event.Message.Content,
		SenderId:    larkSubscriptionEventDecryptedRequest.Event.Sender.SenderId.OpenId,
		SenderType:  larkSubscriptionEventDecryptedRequest.Event.Sender.SenderType,
	}

	// Get all parent message
	var messages []structs.RedisMessage
	if larkSubscriptionEventDecryptedRequest.Event.Message.RootId == nil {
		utils.StoreMessage(*larkSubscriptionEventDecryptedRequest.Event.Message.MessageId, redisMessage)
		messages = utils.GetAllParentMessages(*larkSubscriptionEventDecryptedRequest.Event.Message.MessageId)
	} else {
		utils.StoreMessage(*larkSubscriptionEventDecryptedRequest.Event.Message.RootId, redisMessage)
		messages = utils.GetAllParentMessages(*larkSubscriptionEventDecryptedRequest.Event.Message.RootId)
	}

	switch {
	case strings.HasPrefix(gjson.Get(*messages[0].Content, "text").String(), constants.CommandChatgpt):
		resp := utils.CallOpenAI(messages)
		logrus.Infof("Response from OpenAI: %s\n", resp)
		returnObject := structs.LarkMessageText{Text: resp}
		returnText, _ := json.Marshal(returnObject)
		ReplyMessage(*larkSubscriptionEventDecryptedRequest.Event.Message.MessageId, string(returnText))
	case strings.HasPrefix(gjson.Get(*messages[0].Content, "text").String(), constants.CommandChatgptConfig):
		config := strings.Split(gjson.Get(*messages[0].Content, "text").String(), " ")[1]
		resp := make(map[string]string)
		if slices.Contains(constants.ChatGPTModels[:], config) {
			constants.CurrentChatGPTModel = config
			resp["ConfigResult"] = "Fail"
			resp["CurrentModel"] = constants.CurrentChatGPTModel
			resp["SupportedModels"] = strings.Join(constants.ChatGPTModels[:], ",")
		} else {
			resp["ConfigResult"] = "Success"
			resp["CurrentModel"] = constants.CurrentChatGPTModel
			resp["SupportedModels"] = strings.Join(constants.ChatGPTModels[:], ",")
		}
		respText, _ := json.Marshal(resp)
		returnObject := structs.LarkMessageText{Text: string(respText)}
		returnText, _ := json.Marshal(returnObject)
		ReplyMessage(*larkSubscriptionEventDecryptedRequest.Event.Message.MessageId, string(returnText))
	}
}
