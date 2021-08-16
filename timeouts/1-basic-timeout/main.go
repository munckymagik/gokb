package main

// Examples from https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request ...")

	// even this is no received by the client
	w.WriteHeader(http.StatusOK)

	time.Sleep(2 * time.Second)

	bytesWritten, err := io.WriteString(w, "I am slow!\n")
	if err != nil {
		log.Printf("ERROR %v\n", err)
	}
	log.Printf("Handling request done. Bytes written = %d\n", bytesWritten)
}

func main() {
	srv := http.Server{
		Addr:         ":8888",
		WriteTimeout: 1 * time.Second,
		Handler:      http.HandlerFunc(slowHandler),
	}

	log.Printf("Serving on %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server exited with: %s\n", err)
	}
}
