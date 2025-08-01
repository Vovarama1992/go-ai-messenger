package kafkaadapter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-utils/kafkautil"

	kafka "github.com/segmentio/kafka-go"
)

var ErrProducerClosed = errors.New("producer closed")

type KafkaProducer struct {
	writer             *kafka.Writer
	aiBindingInitTopic string
	adviceTopic        string
	mu                 sync.Mutex
	closed             bool

	retryer *kafkautil.Retry
	breaker *kafkautil.Breaker
}

func (p *KafkaProducer) WithRetryBreaker(retry *kafkautil.Retry, breaker *kafkautil.Breaker) {
	p.retryer = retry
	p.breaker = breaker
}

func NewKafkaProducer(bindingInitTopic, adviceTopic string, writer *kafka.Writer) *KafkaProducer {
	return &KafkaProducer{
		writer:             writer,
		aiBindingInitTopic: bindingInitTopic,
		adviceTopic:        adviceTopic,
	}
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

	writeFn := func(ctx context.Context, msg kafka.Message) error {
		return p.writer.WriteMessages(ctx, msg)
	}

	var writeErr error
	if p.retryer != nil && p.breaker != nil {
		writeErr = p.retryer.Do(func() error {
			return p.breaker.Do(func() error {
				return writeFn(ctx, msg)
			})
		})
	} else {
		writeErr = writeFn(ctx, msg)
	}

	if writeErr != nil {
		log.Printf("❌ Failed to write message to Kafka: %v", writeErr)
		return writeErr
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

func (p *KafkaProducer) SendAdviceRequest(ctx context.Context, payload model.AdviceRequestPayload) error {
	return p.Produce(ctx, p.adviceTopic, payload)
}

func (p *KafkaProducer) PublishAdviceRequest(threadID string) error {
	payload := model.AdviceRequestPayload{
		ThreadID: threadID,
	}
	return p.SendAdviceRequest(context.Background(), payload)
}

func (p *KafkaProducer) SendInvite(ctx context.Context, payload interface{}, topic string) error {
	return p.Produce(ctx, topic, payload)
}
