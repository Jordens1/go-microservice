package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = "3000"
const ServicesURL = "http://localhost:" + ServerPort + "/services"

type registry struct {
	registerRegistions []Registion
	mutex              *sync.RWMutex
}

func (r *registry) add(reg Registion) error {
	r.mutex.Lock()
	r.registerRegistions = append(r.registerRegistions, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(
		patch{
			Added: []patchEntry{
				patchEntry{
					Name: reg.ServiceName,
					URL:  reg.Url,
				},
			},
		})
	return err
}

func (r *registry) notify(pat patch) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, reg := range r.registerRegistions {
		go func(reg Registion) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sedUpdate := false
				for _, added := range pat.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sedUpdate = true
					}
				}
				for _, removed := range pat.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sedUpdate = true
					}
				}
				if sedUpdate {
					err := r.sendPath(p, reg.ServiceUpdateUrl)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}

	return nil

}

func (r *registry) sendRequiredServices(reg Registion) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch

	for _, serviceReg := range r.registerRegistions {
		for _, reqService := range reg.RequiredServices {
			if reqService == serviceReg.ServiceName {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.Url,
				})
			}

		}
	}

	err := r.sendPath(p, reg.ServiceUpdateUrl)
	if err != nil {
		return err
	}

	return nil
}

func (r *registry) sendPath(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))

	if err != nil {
		return err
	}
	return err

}

var reg = registry{
	registerRegistions: make([]Registion, 0),
	mutex:              new(sync.RWMutex),
}

func (r *registry) remove(url string) error {
	for i := range reg.registerRegistions {
		if reg.registerRegistions[i].Url == url {
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: r.registerRegistions[i].ServiceName,
						URL:  r.registerRegistions[i].Url,
					},
				},
			})

			r.mutex.Lock()
			reg.registerRegistions = append(reg.registerRegistions[:i], reg.registerRegistions[i+1:]...)
			r.mutex.Unlock()

			return nil
		}
	}
	return fmt.Errorf("service at url :%s not found ", url)
}

type RegistryService struct{}

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registerRegistions {
			wg.Add(1)
			go func(reg Registion) {
				defer wg.Done()
				success := true
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("heartbeat check passed for %v", reg.ServiceName)
						if !success {
							success = false
							r.remove(reg.Url)
						}
						time.Sleep(1 * time.Second)

					}

				}
			}(reg)
		}
		wg.Wait()
		time.Sleep(freq)
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(time.Second * 3)
	})
}

func (s RegistryService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("request reviced")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registion
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Add service: %v with URl %s \n ", r.ServiceName, r.Url)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)

		log.Printf("removing service %s ", url)
		err = reg.remove(url)

		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
