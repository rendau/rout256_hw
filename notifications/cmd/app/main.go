package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"route256/libs/kafka_consumer"
	"route256/libs/stopsignal"
	"route256/notifications/internal/clients/telegram"
	"route256/notifications/internal/domain"
	"time"
)

type OrderStatusChangeEventSt struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	// telegram
	tg, err := telegram.New(cfg.TelegramToken, cfg.TelegramChatId)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	dm := domain.New(tg, cfg.OrderStatusChangeEventTemplate)

	// consumer
	consumer, err := kafka_consumer.NewKafkaConsumer(kafka_consumer.KafkaConsumerConfig{
		Context: context.Background(),
		Brokers: cfg.KafkaBrokers,
		GroupId: cfg.KafkaGroup,
		Topic:   cfg.OrderStatusChangeTopic,
		Handler: func(ctx context.Context, topic string, msg []byte) bool {
			obj := &OrderStatusChangeEventSt{}
			if json.Unmarshal(msg, obj) != nil { // ignore invalid messages
				return true
			}

			if dm.HandleOrderStatusChangeEvent(domain.OrderStatusChangeEventSt{
				OrderID: obj.OrderID,
				Status:  obj.Status,
			}) != nil {
				return false
			}

			return true
		},
		RetryInterval: 3 * time.Second,
	})
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	log.Println("Start service...")

	consumer.Start()

	// wait for stop signal
	<-stopsignal.StopSignal()

	log.Println("Shutdown service...")

	consumer.Stop()

	log.Println("Exit...")

	os.Exit(0)
}
