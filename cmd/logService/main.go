package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/Jordens1/go-microservice/log"
	"github.com/Jordens1/go-microservice/registry"
	"github.com/Jordens1/go-microservice/service"
)

func main() {
	log.Run("./distributed.log")

	host, port := "localhost", "8080"

	url := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registion{
		ServiceName:      registry.LogService,
		Url:              url,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateUrl: url + "/services",
		HeartbeatURL:     url + "/heartbeat",
	}
	ctx, err := service.Start(context.Background(), host, port, log.RegisterHandlers, r)
	if err != nil {

		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutdown log service .")

}
