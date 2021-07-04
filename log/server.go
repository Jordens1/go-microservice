package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type filelog string

func (fl filelog) Write(date []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	return f.Write(date)

}

func Run(destination string) {
	log = stlog.New(filelog(destination), "[go] : ", stlog.LstdFlags)

}

func RegisterHandlers() {
	http.HandleFunc("/log", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v \n", message)
}
