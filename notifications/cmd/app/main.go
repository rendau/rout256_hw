package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	dbPg "route256/libs/db/pg"
	"route256/libs/kafka_consumer"
	"route256/libs/logger"
	"route256/libs/stopsignal"
	"route256/notifications/internal/clients/telegram"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/domain/models"
	repoPg "route256/notifications/internal/repo/pg"
	"time"
)

type OrderStatusChangeEventSt struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	cfg := ConfigLoad()

	logger.Init(cfg.LogLevel, cfg.Debug)

	db, err := dbPg.New(cfg.DbDsn)
	if err != nil {
		logger.Fatalw(nil, err, "dbPg.New")
	}

	err = db.Migrate("migrations")
	if err != nil {
		logger.Fatalw(nil, err, "db.Migrate")
	}

	repo := repoPg.New(db)

	// telegram
	tg, err := telegram.New(cfg.TelegramToken, cfg.TelegramChatId)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	dm := domain.New(repo, tg, cfg.OrderStatusChangeEventTemplate)

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

			if dm.HandleOrderStatusEvent(&models.OrderStatusEventSt{
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
