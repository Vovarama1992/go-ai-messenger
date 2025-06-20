package stream

import (
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
)

var BindingInitChan = make(chan dto.AiBindingInitPayload, 100)
var BindingResultChan = make(chan dto.ThreadResult, 100)
var FeedChan = make(chan dto.AiFeedPayload, 100)
var AutoReplyChan = make(chan dto.AiAutoReplyResult, 100)
var AdviceRequestChan = make(chan dto.AdviceRequestPayload, 100)
var AdviceResponseChan = make(chan dto.GptAdvice, 100)
