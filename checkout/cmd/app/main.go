package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/checkout/internal/client/loms"
	"route256/checkout/internal/client/productservice"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handler"
	"route256/checkout/pkg/proto/checkout_v1"
	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/stopsignal"
)

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	lomsClient, err := loms.New(cfg.Services.Loms.Url)
	if err != nil {
		log.Fatalln("loms.New: ", err)
	}

	productService, err := productservice.New(cfg.Services.ProductService.Url, cfg.Services.ProductService.Token)
	if err != nil {
		log.Fatalln("productservice.New: ", err)
	}

	model := domain.New(lomsClient, productService)

	// grpc
	grpcHandler := handler.New(model)
	grpcSrv := grpcserver.New()
	reflection.Register(grpcSrv.Server)
	checkout_v1.RegisterCheckoutServer(grpcSrv.Server, grpcHandler)

	// http
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = checkout_v1.RegisterCheckoutHandlerFromEndpoint(context.Background(), mux, "localhost:"+cfg.GrpcPort, opts)
	if err != nil {
		log.Fatalln("RegisterCheckoutHandlerFromEndpoint: ", err)
	}

	// start server

	err = grpcSrv.Start(cfg.GrpcPort)
	if err != nil {
		log.Fatalln("grpcSrv.Start: ", err)
	}

	httpSrv := httpserver.Start(cfg.HttpPort, mux)
	if err != nil {
		log.Fatalln("httpserver.Start: ", err)
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

	log.Println("Shutdown service...")

	if !grpcSrv.Shutdown() {
		exitCode = 1
	}

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	// exit

	log.Println("Exit...")

	os.Exit(exitCode)
}
