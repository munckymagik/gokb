package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/munckymagik/gokb/mapsync"
)

type HandlerState struct {
	repo mapsync.Repo
}

func (handlerState *HandlerState) get(w http.ResponseWriter, req *http.Request) {
	log.Printf("get %v", req)
	name := req.URL.Query().Get("name")
	if val := handlerState.repo.Get(name); val != nil {
		fmt.Fprintf(w, "%s: %d\n", name, val)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}

func main() {
	store := HandlerState{repo: mapsync.NewNaiveRepo()}
	http.HandleFunc("/get", store.get)

	port := 8070
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
