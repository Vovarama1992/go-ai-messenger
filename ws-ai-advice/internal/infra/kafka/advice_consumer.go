package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
	"github.com/segmentio/kafka-go"
)

type adviceConsumer struct {
	reader *kafka.Reader
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewAdviceConsumer(reader *kafka.Reader) ports.AdviceConsumer {
	return &adviceConsumer{reader: reader}
}

func (c *adviceConsumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("🛑 Kafka consumer stopped")
					return
				}
				log.Printf("❌ Kafka read: %v", err)
				continue
			}

			var advice model.GptAdvice
			if err := json.Unmarshal(msg.Value, &advice); err != nil {
				log.Printf("❌ Unmarshal advice: %v", err)
				continue
			}

			stream.PendingAdviceChan <- advice
		}
	}()
	return nil
}

func (c *adviceConsumer) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.reader.Close()
}
