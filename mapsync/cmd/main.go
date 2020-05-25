package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/munckymagik/gokb/mapsync"
)

type HandlerState struct {
	repo mapsync.Repo
}

func (handlerState *HandlerState) get(w http.ResponseWriter, req *http.Request) {
	// log.Printf("get %v", req)
	name := req.URL.Query().Get("name")
	if val := handlerState.repo.Get(name); val != nil {
		fmt.Fprintf(w, "%s: %d\n", name, val)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}

func usage() {
	log.Fatalln("USAGE go run ./cmd [naive|atomic|mutex]")
}

func createRepo(name string) mapsync.Repo {
	var repo mapsync.Repo

	switch name {
	case "naive":
		repo = mapsync.NewNaiveRepo()
	case "atomic":
		repo = mapsync.NewAtomicRepo()
	case "mutex":
		repo = mapsync.NewMutexRepo()
	case "rwmutex":
		repo = mapsync.NewRWMutexRepo()
	default:
		usage()
	}

	return repo
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	store := HandlerState{repo: createRepo(os.Args[1])}
	http.HandleFunc("/get", store.get)

	port := 8070
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
