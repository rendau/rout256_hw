package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	dbPg "route256/libs/db/pg"
	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/kafka_consumer"
	"route256/libs/logger"
	"route256/libs/stopsignal"
	"route256/notifications/internal/clients/telegram"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/domain/models"
	"route256/notifications/internal/handler"
	repoPg "route256/notifications/internal/repo/pg"
	"route256/notifications/pkg/proto/notifications_v1"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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

	// grpc
	grpcHandler := handler.New(dm)
	grpcSrv := grpcserver.New(
		grpc.ChainUnaryInterceptor(
			logger.MiddlewareGRPC,
		),
	)
	reflection.Register(grpcSrv.Server)
	notifications_v1.RegisterNotificationsServer(grpcSrv.Server, grpcHandler)

	// http
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = notifications_v1.RegisterNotificationsHandlerFromEndpoint(context.Background(), mux, "localhost:"+cfg.GrpcPort, opts)
	if err != nil {
		logger.Fatalw(nil, err, "notifications_v1.RegisterNotificationsHandlerFromEndpoint")
	}
	// add health check handler
	if err = mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
	}); err != nil {
		logger.Fatalw(nil, err, "mux.HandlePath")
	}

	log.Println("Start service...")

	consumer.Start()

	// grpc
	err = grpcSrv.Start(cfg.GrpcPort)
	if err != nil {
		logger.Fatalw(nil, err, "grpcSrv.Start")
	}

	// http
	httpSrv := httpserver.Start(cfg.HttpPort, mux)
	if err != nil {
		logger.Fatalw(nil, err, "httpSrv.Start")
	}

	exitCode := 0

	select {
	case <-grpcSrv.Wait():
		exitCode = 1
	case <-httpSrv.Wait():
		exitCode = 1
	case <-stopsignal.StopSignal():
	}

	log.Println("Shutdown service...")

	consumer.Stop()

	if !grpcSrv.Shutdown() {
		exitCode = 1
	}

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	log.Println("Exit...")

	os.Exit(exitCode)
}
