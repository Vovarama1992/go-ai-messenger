package kafkaadapter

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

type InviteMessage struct {
	ChatID int64 `json:"chatId"`
	UserID int64 `json:"userId"`
}

type InviteHandler func(msg InviteMessage)

type InviteConsumerInterface interface {
	Start()
	Stop()
}

type InviteConsumer struct {
	reader  *kafka.Reader
	handler InviteHandler
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewInviteConsumer(topic string, handler InviteHandler) *InviteConsumer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	ctx, cancel := context.WithCancel(context.Background())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "", // читаем всё, без групп
		Topic:   topic,
	})

	return &InviteConsumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ic *InviteConsumer) Start() {
	ic.wg.Add(1)
	go func() {
		defer ic.wg.Done()
		defer ic.reader.Close()

		for {
			m, err := ic.reader.ReadMessage(ic.ctx)
			if err != nil {
				if ic.ctx.Err() != nil {
					return
				}
				log.Printf("❌ Kafka read error (invite): %v", err)
				continue
			}

			var msg InviteMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("❌ Unmarshal error (invite): %v", err)
				continue
			}

			ic.handler(msg)
		}
	}()
}

func (ic *InviteConsumer) Stop() {
	ic.cancel()
	ic.wg.Wait()
}
