package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/hndwrapper"
	"route256/libs/httpserver"
	"route256/libs/stopsignal"
)

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	lomsClient := loms.New(cfg.Services.Loms.Url)
	productService := productservice.New(cfg.Services.ProductService.Url, cfg.Services.ProductService.Token)

	model := domain.New(lomsClient, productService)

	// routes
	http.Handle("/addToCart", hndwrapper.New(addtocart.New(model).Handle))
	http.Handle("/deleteFromCart", hndwrapper.New(deletefromcart.New(model).Handle))
	http.Handle("/listCart", hndwrapper.New(listcart.New(model).Handle))
	http.Handle("/purchase", hndwrapper.New(purchase.New(model).Handle))

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
