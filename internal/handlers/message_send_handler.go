package handlers

import (
	"Lark-Bot/internal/structs"
	"Lark-Bot/internal/utils"
	"context"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

func SendMessage(message string, chatType string, id string) {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(chatType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(id).
			MsgType(`text`).
			Content(message).
			Uuid(utils.GetUUID()).
			Build()).Build()

	resp, err := utils.GetLarkClient().Im.Message.Create(context.Background(), req)

	if err != nil {
		logrus.Warnln(err)
	}

	if !resp.Success() {
		logrus.Warnln(resp.Code, resp.Msg, resp.RequestId())
	}
}

func ReplyMessage(messageId string, content string) {
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			Content(content).
			MsgType(`text`).
			Build()).
		Build()

	resp, err := utils.GetLarkClient().Im.Message.Reply(context.Background(), req)

	redisMessage := structs.RedisMessage{
		MessageId:   resp.Data.MessageId,
		RootId:      resp.Data.RootId,
		ParentId:    resp.Data.ParentId,
		CreateTime:  resp.Data.CreateTime,
		ChatId:      resp.Data.ChatId,
		ChatType:    nil,
		MessageType: resp.Data.MsgType,
		Content:     resp.Data.Body.Content,
		SenderId:    resp.Data.Sender.Id,
		SenderType:  resp.Data.Sender.SenderType,
	}

	utils.StoreMessage(*resp.Data.RootId, redisMessage)

	if err != nil {
		logrus.Warnln(err)
	}

	if !resp.Success() {
		logrus.Warnln(resp.Code, resp.Msg, resp.RequestId())
	}
}
