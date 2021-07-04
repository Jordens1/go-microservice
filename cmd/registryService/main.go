package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Jordens1/go-microservice/registry"
)

func main() {
	registry.SetupRegistryService()

	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = "localhost:" + registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("started. press any key to stop ")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("shutdown registry service")

}
