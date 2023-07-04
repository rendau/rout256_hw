package kafka_producer

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type assuranceLevel string

const (
	AssuranceLevelAtLeastOnce assuranceLevel = "at_least_once"
	AssuranceLevelAtMostOnce  assuranceLevel = "at_most_once"
	AssuranceLevelExactlyOnce assuranceLevel = "exactly_once"
)

type KafkaProducerConfig struct {
	Brokers        []string
	Topic          string
	Compress       bool
	AssuranceLevel assuranceLevel
}

type KafkaProducer struct {
	p     sarama.SyncProducer
	Topic string
}

func NewKafkaProducer(cfg KafkaProducerConfig) (*KafkaProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Partitioner = sarama.NewHashPartitioner

	// assurance level
	if cfg.AssuranceLevel == AssuranceLevelAtLeastOnce {
		config.Producer.RequiredAcks = sarama.WaitForLocal
	} else if cfg.AssuranceLevel == AssuranceLevelAtMostOnce {
		config.Producer.RequiredAcks = sarama.NoResponse
	} else { // default is AssuranceLevelExactlyOnce
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Idempotent = true
		config.Net.MaxOpenRequests = 1
	}

	// compression
	if cfg.Compress {
		config.Producer.CompressionLevel = sarama.CompressionLevelDefault
		config.Producer.Compression = sarama.CompressionGZIP
	} else {
		config.Producer.Compression = sarama.CompressionNone
	}

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &KafkaProducer{
		p:     producer,
		Topic: cfg.Topic,
	}, nil
}

func (o *KafkaProducer) SendMessage(key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: o.Topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := o.p.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (o *KafkaProducer) Close() error {
	return o.p.Close()
}
