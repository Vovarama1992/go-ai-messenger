package stream

import "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"

var PendingAdviceChan = make(chan model.GptAdvice, 100)
