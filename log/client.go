package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"

	"github.com/Jordens1/go-microservice/registry"
)

func SetClientLog(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[v] - %s", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})

}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("faild to send message .service response code is : %v", res.StatusCode)
	}
	return len(data), nil
}
