package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/Jordens1/go-microservice/grades"
	"github.com/Jordens1/go-microservice/log"
	"github.com/Jordens1/go-microservice/registry"
	"github.com/Jordens1/go-microservice/service"
)

func main() {

	host, port := "localhost", "8888"

	url := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registion{
		ServiceName:      registry.GradingService,
		Url:              url,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateUrl: url + "/services",
		HeartbeatURL:     url + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), host, port, grades.RegisterHandlers, r)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("logging service found at : %s \n", logProvider)
		log.SetClientLog(logProvider, r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("shutdown grading service .")

}
