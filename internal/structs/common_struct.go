package structs

type LarkMessageText struct {
	Text string `json:"text"`
}

type RedisMessage struct {
	MessageId   *string `json:"message_id"`
	RootId      *string `json:"root_id"`
	ParentId    *string `json:"parent_id"`
	CreateTime  *string `json:"create_time"`
	ChatId      *string `json:"chat_id"`
	ChatType    *string `json:"chat_type"`
	MessageType *string `json:"message_type"`
	Content     *string `json:"content"`
	SenderId    *string `json:"sender_id"`
	SenderType  *string `json:"sender_type"`
}
