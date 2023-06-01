package main

import (
	"log"
	"net/http"
	"os"
	"route256/libs/hndwrapper"
	"route256/libs/httpserver"
	"route256/libs/stopsignal"
	"route256/loms/internal/domain"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
	"time"
)

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	model := domain.New()

	// routes
	http.Handle("/createOrder", hndwrapper.New(createorder.New(model).Handle))
	http.Handle("/listOrder", hndwrapper.New(listorder.New(model).Handle))
	http.Handle("/orderPayed", hndwrapper.New(orderpayed.New(model).Handle))
	http.Handle("/cancelOrder", hndwrapper.New(cancelorder.New(model).Handle))
	http.Handle("/stocks", hndwrapper.New(stocks.New(model).Handle))

	// start server
	httpSrv := httpserver.Start(cfg.HttpListen, http.DefaultServeMux)

	exitCode := 0

	select {
	case <-httpSrv.Wait():
		exitCode = 1
	case <-stopsignal.StopSignal():
	}

	log.Println("Shutdown service...")

	if !httpSrv.Shutdown(10 * time.Second) {
		exitCode = 1
	}

	log.Println("Exit...")

	os.Exit(exitCode)
}
