package handlers

import (
	"Lark-Bot/internal/structs"
	"Lark-Bot/internal/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func HandleReceivedMessage(message string, openId string) {
	resp := utils.CallOpenAI(message)
	logrus.Infof("Response from OpenAI: %s\n", resp)
	returnObject := structs.LarkMessageText{Text: resp}
	returnText, _ := json.Marshal(returnObject)
	SendMessage(string(returnText), openId)
}
