package handlers

import (
	"Lark-Bot/internal/constants"
	"Lark-Bot/internal/structs"
	"Lark-Bot/internal/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strings"
)

func HandleReceivedMessage(larkSubscriptionEventDecryptedRequest structs.LarkSubscriptionEventDecryptedRequest) {
	var larkMessageText structs.LarkMessageText
	json.Unmarshal([]byte(larkSubscriptionEventDecryptedRequest.Event.Message.Content), &larkMessageText)

	switch {
	case strings.HasPrefix(larkMessageText.Text, constants.CommandChatgpt):
		resp := utils.CallOpenAI(larkMessageText.Text)
		logrus.Infof("Response from OpenAI: %s\n", resp)
		returnObject := structs.LarkMessageText{Text: resp}
		returnText, _ := json.Marshal(returnObject)
		switch larkSubscriptionEventDecryptedRequest.Event.Message.ChatType {
		case constants.ChatTypeGroup:
			SendMessage(string(returnText), `chat_id`, larkSubscriptionEventDecryptedRequest.Event.Message.ChatId)
		case constants.ChatTypeP2p:
			SendMessage(string(returnText), `open_id`, larkSubscriptionEventDecryptedRequest.Event.Sender.SenderId.OpenId)
		}
	}
}
