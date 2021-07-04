package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/Jordens1/go-microservice/log"
	"github.com/Jordens1/go-microservice/portal"
	"github.com/Jordens1/go-microservice/registry"
	"github.com/Jordens1/go-microservice/service"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "6060"
	url := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registion{
		ServiceName: registry.PortalService,
		Url:         url,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.PortalService,
		},
		ServiceUpdateUrl: url + "/services",
		HeartbeatURL:     url + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(),
		host, port, portal.RegisterHandlers, r)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		log.SetClientLog(logProvider, r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("shutdown protal.")

}
