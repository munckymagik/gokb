package main

// Examples from https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func slowAPICall() string {
	d := rand.Intn(5)
	select {
	case <-time.After(time.Duration(d) * time.Second):
		log.Printf("Slow API call done after %d seconds.\n", d)
		return "foobar"
	}
}

func slowAPICallHandler(w http.ResponseWriter, req *http.Request) {
	result := slowAPICall()
	_, _ = io.WriteString(w, result+"\n")
}

func main() {
	srv := http.Server{
		Addr:         ":8888",
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(slowAPICallHandler), 1*time.Second, "Timeout\n"),
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server exited with: %s\n", err)
	}
}
