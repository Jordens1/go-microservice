package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registion) error {

	heartbeaURL, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeaURL.Path, func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	serviceUpdateURL, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buffer)
	if err != nil {
		return err

	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("faild to registry service , repose code is %v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("update reviced %v\n", p)
	prov.update(p)
}

func ShutdownService(url string) error {

	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("faild to registry service , registry reponse ceode is %v", res.StatusCode)

	}
	return nil
}

type provides struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *provides) update(pat patch) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				p.services[patchEntry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
			}
		}
	}

}

//可能是多个url返回，应该是是[]string
func (p provides) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no provides for this service : %v", name)
	}

	idx := int(rand.Float32() * float32(len(providers)))

	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = provides{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
