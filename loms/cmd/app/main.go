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

	dbPg "route256/libs/db/pg"
	"route256/libs/grpcserver"
	"route256/libs/httpserver"
	"route256/libs/stopsignal"
	"route256/loms/internal/domain"
	"route256/loms/internal/handler"
	repoPg "route256/loms/internal/repo/pg"
	"route256/loms/pkg/proto/loms_v1"
)

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	db, err := dbPg.New(cfg.DbDsn)
	if err != nil {
		log.Fatalln("pg.New: ", err)
	}

	err = db.Migrate("migrations")
	if err != nil {
		log.Fatalln("db.Migrate: ", err)
	}

	repo := repoPg.New(db)

	dm := domain.New(db, repo)

	// grpc
	grpcHandler := handler.New(dm)
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

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	log.Println("Exit...")

	os.Exit(exitCode)
}
