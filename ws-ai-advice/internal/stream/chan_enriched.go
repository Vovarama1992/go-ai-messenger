package stream

import "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"

var EnrichedAdviceChan = make(chan model.EnrichedAdvice, 100)
