package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)

}

type studentHandler struct{}

func (sh studentHandler) getAll(rw http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJson(students)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("faild to serialize students : %q", err)

		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(data)
}

func (sh studentHandler) toJson(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("faild to serialize students : %q", err)
	}
	return buffer.Bytes(), nil
}

func (sh studentHandler) getOne(rw http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetById(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("faild to serialize students : %q", err)

		return
	}

	data, err := sh.toJson(student)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("faild to serialize students : %q", err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(data)
}

func (sh studentHandler) addGrade(rw http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetById(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("faild to serialize students : %q", err)

		return
	}

	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	student.Grades = append(student.Grades, g)
	rw.WriteHeader(http.StatusCreated)
	data, err := sh.toJson(g)
	if err != nil {
		log.Println(err)
		return

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(data)
}

func (sh studentHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		sh.getAll(rw, r)

	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(rw, r, id)
	case 4:

		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(rw, r, id)
	default:
		rw.WriteHeader(http.StatusNotFound)

	}

}
