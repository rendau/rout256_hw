package main

import (
	"context"
	"log"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/stopsignal"
	"route256/loms/internal/domain"
	"route256/loms/internal/handler"
	"route256/loms/pkg/proto/loms_v1"

	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	model := domain.New()

	// grpc
	grpcHandler := handler.New(model)
	grpcSrv := grpcserver.New()
	reflection.Register(grpcSrv.Server)
	loms_v1.RegisterLomsServer(grpcSrv.Server, grpcHandler)

	// http
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = loms_v1.RegisterLomsHandlerFromEndpoint(context.Background(), mux, "localhost:"+cfg.GrpcPort, opts)
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

	log.Println("Shutdown service...")

	if !grpcSrv.Shutdown() {
		exitCode = 1
	}

	log.Println("Exit...")

	os.Exit(exitCode)
}
