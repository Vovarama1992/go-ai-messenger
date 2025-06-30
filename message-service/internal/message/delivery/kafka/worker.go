package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg dto.IncomingMessage) error
}

func StartMessageWorkers(ctx context.Context, wg *sync.WaitGroup, handler MessageHandler, msgChan <-chan dto.IncomingMessage, n int) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					log.Printf("worker %d: shutting down", id)
					return
				case msg := <-msgChan:
					if err := handler.Handle(ctx, msg); err != nil {
						log.Printf("worker %d: failed to handle msg: %v", id, err)
					}
				}
			}
		}(i)
	}
}
