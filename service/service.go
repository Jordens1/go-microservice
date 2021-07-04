package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Jordens1/go-microservice/registry"
)

func Start(ctx context.Context, host, port string, RegisterHandlersFunc func(), reg registry.Registion) (context.Context, error) {
	RegisterHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)

	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil

}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancle := context.WithCancel(ctx)
	var svc http.Server
	svc.Addr = host + ":" + port

	go func() {
		log.Println(svc.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancle()
	}()

	go func() {
		fmt.Printf("%v started. press any key to stop ", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		svc.Shutdown(ctx)
		cancle()
	}()

	return ctx
}
