package main

import (
	"fmt"
	"log"
	"net/http"
)

type CommandType int

const (
	GetCommand = iota
	SetCommand
	IncCommand
)

type Command struct {
	ty        CommandType
	name      string
	val       int
	replyChan chan int
}

func startCounterManager() chan<- Command {
	counters := make(map[string]int)
	cmds := make(chan Command)

	go func() {
		for cmd := range cmds {
			switch cmd.ty {
			case GetCommand:
				cmd.replyChan <- counters[cmd.name]
			case SetCommand:
				counters[cmd.name] = cmd.val
				cmd.replyChan <- cmd.val
			case IncCommand:
				counters[cmd.name]++
				cmd.replyChan <- counters[cmd.name]
			default:
				log.Fatal("unknown command type", cmd.ty)
			}
		}
	}()

	return cmds
}

type Server struct {
	cmds chan<- Command
}

func (server *Server) get(response http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	replyChan := make(chan int, 1)
	server.cmds <- Command{ty: GetCommand, name: name, replyChan: replyChan}

	reply := <-replyChan
	fmt.Fprintf(response, "%s: %d\n", name, reply)
}

func (server *Server) inc(response http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	replyChan := make(chan int, 1)
	server.cmds <- Command{ty: IncCommand, name: name, replyChan: replyChan}

	reply := <-replyChan
	fmt.Fprintf(response, "%s: %d\n", name, reply)
}

func limitNumClients(fn http.HandlerFunc, maxClients int) http.HandlerFunc {
	sema := make(chan struct{}, maxClients)

	return func(response http.ResponseWriter, request *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		fn(response, request)
	}
}

func main() {
	cmds := startCounterManager()
	store := Server{cmds: cmds}

	http.HandleFunc("/inc", store.inc)
	http.HandleFunc("/get", store.get)

	port := 8070
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
