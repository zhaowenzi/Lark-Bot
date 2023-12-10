package constants

import "github.com/sashabaranov/go-openai"

const (
	UrlVerification      = `url_verification`
	MessageReceive       = `im.message.receive_v1`
	ChatTypeGroup        = `group`
	ChatTypeP2p          = `p2p`
	CommandChatgpt       = `/chatgpt`
	CommandChatgptConfig = `/chatgptConfig`
	SenderTypeUser       = `user`
)

var ChatGPTModels = [...]string{openai.GPT3Dot5Turbo, openai.GPT4, openai.GPT432K}

var CurrentChatGPTModel = ChatGPTModels[0]
