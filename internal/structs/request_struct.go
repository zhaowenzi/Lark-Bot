package structs

type LarkSubscriptionEventEncryptRequest struct {
	Encrypt *string `json:"encrypt" binding:"required"`
}

type LarkSubscriptionEventDecryptedEventHeader struct {
	EventId    *string `json:"event_id"`
	EventType  *string `json:"event_type"`
	CreateTime *string `json:"create_time"`
	Token      *string `json:"token"`
	AppId      *string `json:"app_id"`
	TenantKey  *string `json:"tenant_key"`
}

type LarkSubscriptionEventDecryptedUserId struct {
	UnionId *string `json:"union_id"`
	UserId  *string `json:"user_id"`
	OpenId  *string `json:"open_id"`
}

type LarkSubscriptionEventDecryptedEventSender struct {
	SenderId   LarkSubscriptionEventDecryptedUserId `json:"sender_id"`
	SenderType *string                              `json:"sender_type"`
	TenantKey  *string                              `json:"tenant_key"`
}

type LarkSubscriptionEventDecryptedMentionEvent struct {
	Key       *string                              `json:"key"`
	Id        LarkSubscriptionEventDecryptedUserId `json:"id"`
	Name      *string                              `json:"name"`
	TenantKey *string                              `json:"tenant_key"`
}

type LarkSubscriptionEventDecryptedEventMessage struct {
	MessageId   *string                                      `json:"message_id"`
	RootId      *string                                      `json:"root_id"`
	ParentId    *string                                      `json:"parent_id"`
	CreateTime  *string                                      `json:"create_time"`
	ChatId      *string                                      `json:"chat_id"`
	ChatType    *string                                      `json:"chat_type"`
	MessageType *string                                      `json:"message_type"`
	Content     *string                                      `json:"content"`
	Mentions    []LarkSubscriptionEventDecryptedMentionEvent `json:"mentions"`
}

type LarkSubscriptionEventDecryptedEvent struct {
	Sender  LarkSubscriptionEventDecryptedEventSender  `json:"sender"`
	Message LarkSubscriptionEventDecryptedEventMessage `json:"message"`
}

type LarkSubscriptionEventDecryptedRequest struct {
	Challenge *string                                   `json:"challenge"`
	Token     *string                                   `json:"token"`
	Type      *string                                   `json:"type"`
	Schema    *string                                   `json:"schema"`
	Header    LarkSubscriptionEventDecryptedEventHeader `json:"header"`
	Event     LarkSubscriptionEventDecryptedEvent       `json:"event"`
}
