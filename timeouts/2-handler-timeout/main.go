package main

// Examples from https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Delaying ...")
	time.Sleep(2 * time.Second)
	fmt.Println("Responding.")
	_, _ = io.WriteString(w, "I am slow!\n")
}

func main() {
	srv := http.Server{
		Addr:         ":8888",
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(slowHandler), 1*time.Second, "Timeout\n"),
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server exited with: %s\n", err)
	}
}
