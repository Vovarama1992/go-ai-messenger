package kafkaadapter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"

	kafka "github.com/segmentio/kafka-go"
)

var ErrProducerClosed = errors.New("producer closed")

type KafkaProducer struct {
	writer             *kafka.Writer
	aiBindingInitTopic string
	mu                 sync.Mutex
	closed             bool
}

func NewKafkaProducer(aiBindingInitTopic string) *KafkaProducer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer:             writer,
		aiBindingInitTopic: aiBindingInitTopic,
	}
}

// Start — если нужно будет запускать, сейчас просто заглушка
func (p *KafkaProducer) Start() error {
	return nil
}

func (p *KafkaProducer) Produce(ctx context.Context, topic string, payload interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return ErrProducerClosed
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Value: bytes,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("❌ Failed to write message to Kafka: %v", err)
		return err
	}

	log.Printf("✅ Sent message to Kafka topic: %s", topic)
	return nil
}

func (p *KafkaProducer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	err := p.writer.Close()
	if err != nil {
		return err
	}

	p.closed = true
	return nil
}

func (p *KafkaProducer) SendAiBindingInit(ctx context.Context, payload model.AiBindingInitPayload) error {
	return p.Produce(ctx, p.aiBindingInitTopic, payload)
}
