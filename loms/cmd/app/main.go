package main

import (
	"context"
	"net/http"
	"os"
	dbPg "route256/libs/db/pg"
	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/kafka_producer"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/stopsignal"
	"route256/libs/tracer"
	"route256/loms/internal/domain"
	"route256/loms/internal/handler"
	repoPg "route256/loms/internal/repo/pg"
	"route256/loms/pkg/proto/loms_v1"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := ConfigLoad()

	logger.Init(cfg.LogLevel, cfg.Debug)

	err := tracer.InitGlobal(cfg.JaegerHostPort, "loms")
	if err != nil {
		logger.Fatalw(nil, err, "tracer.InitGlobal")
	}

	metrics.Init("loms")

	db, err := dbPg.New(cfg.DbDsn)
	if err != nil {
		logger.Fatalw(nil, err, "dbPg.New")
	}

	err = db.Migrate("migrations")
	if err != nil {
		logger.Fatalw(nil, err, "db.Migrate")
	}

	repo := repoPg.New(db)

	var eventProducerOrderStatusChange domain.IEventProducer = nil

	if len(cfg.OrderStatusChangeNotifyBrokers) > 0 && cfg.OrderStatusChangeNotifierTopic != "" {
		eventProducerOrderStatusChange, err = kafka_producer.NewKafkaProducer(kafka_producer.KafkaProducerConfig{
			Brokers:        cfg.OrderStatusChangeNotifyBrokers,
			Topic:          cfg.OrderStatusChangeNotifierTopic,
			Compress:       false,
			AssuranceLevel: kafka_producer.AssuranceLevelExactlyOnce,
		})
		if err != nil {
			logger.Fatalw(nil, err, "kafka_producer.NewKafkaProducer")
		}
	}

	dm := domain.New(db, repo, eventProducerOrderStatusChange)

	// grpc
	grpcHandler := handler.New(dm)
	grpcSrv := grpcserver.New(
		grpc.ChainUnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer.GetTracer()),
			logger.MiddlewareGRPC,
			tracer.MiddlewareGRPC,
			metrics.MiddlewareGRPC,
		),
	)
	reflection.Register(grpcSrv.Server)
	loms_v1.RegisterLomsServer(grpcSrv.Server, grpcHandler)

	// http
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = loms_v1.RegisterLomsHandlerFromEndpoint(context.Background(), mux, "localhost:"+cfg.GrpcPort, opts)
	if err != nil {
		logger.Fatalw(nil, err, "loms_v1.RegisterLomsHandlerFromEndpoint")
	}
	// add health check handler
	if err = mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
	}); err != nil {
		logger.Fatalw(nil, err, "mux.HandlePath")
	}
	// add metrics handler
	if err = mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		logger.Fatalw(nil, err, "something wrong with metrics handler")
	}

	// start server

	logger.Infow(nil, "Start...")

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

	logger.Infow(nil, "Shutdown...")

	if !grpcSrv.Shutdown() {
		exitCode = 1
	}

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	logger.Infow(nil, "Exit...")

	os.Exit(exitCode)
}
