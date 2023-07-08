package main

import (
	"context"
	"net/http"
	"os"
	"route256/checkout/internal/client/loms"
	"route256/checkout/internal/client/productservice"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handler"
	repoPg "route256/checkout/internal/repo/pg"
	"route256/checkout/pkg/proto/checkout_v1"
	dbPg "route256/libs/db/pg"
	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/stopsignal"
	"route256/libs/tracer"
	"time"

	"github.com/opentracing-contrib/go-grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := ConfigLoad()

	logger.Init(cfg.LogLevel, cfg.Debug)

	err := tracer.InitGlobal(cfg.JaegerHostPort, "checkout")
	if err != nil {
		logger.Fatalw(nil, err, "tracer.InitGlobal")
	}

	metrics.Init("checkout")

	db, err := dbPg.New(cfg.DbDsn)
	if err != nil {
		logger.Fatalw(nil, err, "dbPg.New")
	}

	err = db.Migrate("migrations")
	if err != nil {
		logger.Fatalw(nil, err, "db.Migrate")
	}

	repo := repoPg.New(db)

	lomsClient, err := loms.New(cfg.Services.Loms.Url)
	if err != nil {
		logger.Fatalw(nil, err, "loms.New")
	}

	productService, err := productservice.New(cfg.Services.ProductService.Url, cfg.Services.ProductService.Token)
	if err != nil {
		logger.Fatalw(nil, err, "productservice.New")
	}

	dm := domain.New(repo, lomsClient, productService)

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
	checkout_v1.RegisterCheckoutServer(grpcSrv.Server, grpcHandler)

	// http
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = checkout_v1.RegisterCheckoutHandlerFromEndpoint(context.Background(), mux, "localhost:"+cfg.GrpcPort, opts)
	if err != nil {
		logger.Fatalw(nil, err, "checkout_v1.RegisterCheckoutHandlerFromEndpoint")
	}
	// add health check handler
	if err = mux.HandlePath(http.MethodGet, "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
		logger.Fatalw(nil, err, "httpserver.Start")
	}

	exitCode := 0

	select {
	case <-grpcSrv.Wait():
		exitCode = 1
	case <-httpSrv.Wait():
		exitCode = 1
	case <-stopsignal.StopSignal():
	}

	// shutdown

	logger.Infow(nil, "Shutdown...")

	if !grpcSrv.Shutdown() {
		exitCode = 1
	}

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	// exit

	logger.Infow(nil, "Exit...")

	os.Exit(exitCode)
}
