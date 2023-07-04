package kafka_consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConsumerConfig struct {
	Context       context.Context
	GroupId       string
	Brokers       []string
	Topic         string
	Handler       func(ctx context.Context, topic string, msg []byte) bool
	RetryInterval time.Duration
	SkipUnread    bool
}

type KafkaConsumer struct {
	cg            sarama.ConsumerGroup
	Context       context.Context
	ContextCancel context.CancelFunc
	Topic         string
	Handler       func(ctx context.Context, topic string, msg []byte) bool
	RetryInterval time.Duration
	wg            *sync.WaitGroup
}

func NewKafkaConsumer(cfg KafkaConsumerConfig) (*KafkaConsumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	if cfg.SkipUnread {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// create kafka consumer group
	cg, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupId, config)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	if cfg.RetryInterval == 0 {
		cfg.RetryInterval = 2 * time.Second
	}

	ctx, cancel := context.WithCancel(cfg.Context)
	wg := &sync.WaitGroup{}

	return &KafkaConsumer{
		cg:            cg,
		Context:       ctx,
		ContextCancel: cancel,
		Topic:         cfg.Topic,
		Handler:       cfg.Handler,
		RetryInterval: cfg.RetryInterval,
		wg:            wg,
	}, nil
}

func (o *KafkaConsumer) Start() {
	o.wg.Add(1)
	go o.consumeRoutine()
}

func (o *KafkaConsumer) Stop() {
	//_ = o.cg.Close() // sarama has a bug https://github.com/Shopify/sarama/issues/1351
	o.ContextCancel()
	o.wg.Wait()
}

func (o *KafkaConsumer) consumeRoutine() {
	defer o.wg.Done()

	var err error

	for {
		err = o.cg.Consume(o.Context, []string{o.Topic}, o)
		if err != nil {
			fmt.Println("Error occurred on consume:", err)
		}

		if !sleepWithContext(o.Context, o.RetryInterval) {
			return
		}
	}
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
			if o.handleMessage(msg, session) != nil {
				return nil
			}
		}
	}
}

func (o *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) error {
	ctx := o.Context
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if o.Handler(ctx, msg.Topic, msg.Value) {
				session.MarkMessage(msg, "")
				return nil
			} else {
				if !sleepWithContext(ctx, o.RetryInterval) {
					return nil
				}
			}
		}
	}
}

func sleepWithContext(ctx context.Context, dur time.Duration) bool {
	timer := time.NewTimer(dur)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return false
	case <-timer.C:
		return true
	}
}
