package kafka_consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConsumerConfig struct {
	Context              context.Context
	GroupId              string
	Brokers              []string
	Topic                string
	Handler              func(ctx context.Context, topic string, msg []byte) bool
	HandlerRetryInterval time.Duration
}

type KafkaConsumer struct {
	p                    sarama.ConsumerGroup
	Context              context.Context
	ContextCancel        context.CancelFunc
	Topic                string
	Handler              func(ctx context.Context, topic string, msg []byte) bool
	HandlerRetryInterval time.Duration
}

func NewKafkaConsumer(cfg KafkaConsumerConfig) (*KafkaConsumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// create consumer
	p, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupId, config)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(cfg.Context)

	if cfg.HandlerRetryInterval == 0 {
		cfg.HandlerRetryInterval = 2 * time.Second
	}

	return &KafkaConsumer{
		p:                    p,
		Context:              ctx,
		ContextCancel:        cancel,
		Topic:                cfg.Topic,
		Handler:              cfg.Handler,
		HandlerRetryInterval: cfg.HandlerRetryInterval,
	}, nil
}

func (o *KafkaConsumer) Stop() error {
	o.ContextCancel()
	return o.p.Close()
}

// private methods for ConsumerGroupHandler interface:

func (o *KafkaConsumer) Setup(ses sarama.ConsumerGroupSession) error {
	return nil
}

func (o *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (o *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var ok bool
	var msg *sarama.ConsumerMessage

	for {
		select {
		case <-o.Context.Done():
			return nil
		case msg, ok = <-claim.Messages():
			if !ok {
				return nil
			}
			if o.consumeClaimMessage(msg, session) != nil {
				return nil
			}
		}
	}
}

func (o *KafkaConsumer) consumeClaimMessage(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) error {
	var timer *time.Timer
	ctx := session.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if o.Handler(ctx, msg.Topic, msg.Value) {
				session.MarkMessage(msg, "")
				return nil
			} else {
				timer = time.NewTimer(o.HandlerRetryInterval)
				select {
				case <-ctx.Done():
					if !timer.Stop() {
						<-timer.C
					}
				case <-timer.C:
				}
			}
		}
	}
}
