package main

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
}

func (h *Handler) handle(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "hello")
}

func main() {
	fmt.Println("hello")
	handler := Handler{}
	http.HandleFunc("/", handler.handle)

	log.Fatal(http.ListenAndServe(":8070", nil))
}
