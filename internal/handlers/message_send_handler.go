package handlers

import (
	"Lark-Bot/internal/utils"
	"context"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

func SendMessage(message string, chat_type string, id string) {
	logrus.Infoln(message)
	logrus.Infoln(id)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(chat_type).
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
