package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type CounterStore struct {
	sync.Mutex
	counters map[string]int
}

func (cs *CounterStore) get(w http.ResponseWriter, req *http.Request) {
	cs.Lock()
	defer cs.Unlock()

	log.Printf("get %v", req)
	name := req.URL.Query().Get("name")
	if val, ok := cs.counters[name]; ok {
		fmt.Fprintf(w, "%s: %d\n", name, val)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}

func (cs *CounterStore) inc(response http.ResponseWriter, request *http.Request) {
	cs.Lock()
	defer cs.Unlock()

	name := request.URL.Query().Get("name")
	cs.counters[name]++
	response.WriteHeader(200)
}

func main() {
	store := CounterStore{counters: map[string]int{}}
	http.HandleFunc("/inc", store.inc)
	http.HandleFunc("/get", store.get)

	port := 8070
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
