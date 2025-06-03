package kafka

import "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"

var MessageChan = make(chan dto.IncomingMessage, 100)
