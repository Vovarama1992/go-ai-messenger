package kafkaadapter

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

type ForwardMessage struct {
	ChatID      int64  `json:"chatId"`
	SenderID    int64  `json:"senderId"`
	Text        string `json:"text"`
	AIGenerated bool   `json:"aiGenerated"`
}

type ForwardHandler func(msg ForwardMessage)

type ForwardConsumerInterface interface {
	Start()
	Stop()
}

type ForwardConsumer struct {
	reader  *kafka.Reader
	handler ForwardHandler
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewForwardConsumer(topic string, handler ForwardHandler) *ForwardConsumer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	ctx, cancel := context.WithCancel(context.Background())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "", // читаем все сообщения
		Topic:   topic,
	})

	return &ForwardConsumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (fc *ForwardConsumer) Start() {
	fc.wg.Add(1)
	go func() {
		defer fc.wg.Done()
		defer fc.reader.Close()

		for {
			m, err := fc.reader.ReadMessage(fc.ctx)
			if err != nil {
				if fc.ctx.Err() != nil {
					// Context cancelled — нормальное завершение
					return
				}
				log.Printf("❌ Kafka read error: %v", err)
				continue
			}

			var msg ForwardMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("❌ Unmarshal error: %v", err)
				continue
			}

			fc.handler(msg)
		}
	}()
}

func (fc *ForwardConsumer) Stop() {
	fc.cancel()
	fc.wg.Wait()
}
