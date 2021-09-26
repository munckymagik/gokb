package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type handler struct {
	readLatency  time.Duration
	writeLatency time.Duration
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t0 := time.Now()
	log.Printf("proto:%s method:%s path:%s\n", req.Proto, req.Method, req.URL.Path)

	select {
	case <-req.Context().Done():
		log.Printf("∟ Context was cancelled after %v while delaying read\n", time.Since(t0))
		return
	case <-time.After(h.readLatency):
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("ERROR while writing response: %v\n", err)
			return
		}
		log.Printf("∟ body:\"%s\"\n", body)
	}

	select {
	case <-req.Context().Done():
		log.Printf("∟ Context was cancelled after %v waiting to write\n", time.Since(t0))
		return
	case <-time.After(h.writeLatency):
		if _, err := io.WriteString(w, "Server OK"); err != nil {
			log.Printf("∟ ERROR while writing response: %v\n", err)
			return
		}
		log.Printf("∟ duration:%v\n", time.Since(t0))
	}
}

func main() {
	var (
		readHeaderTimeout time.Duration
		readTimeout       time.Duration
		handlerTimeout    time.Duration
		writeTimeout      time.Duration
		readLatency       time.Duration
		writeLatency      time.Duration
		idleTimeout       time.Duration
		addrHost          string
		addrPort          uint
	)
	flag.DurationVar(&readHeaderTimeout, "readHeaderTimeout", 0, "Set ReadHeaderTimeout")
	flag.DurationVar(&readTimeout, "readTimeout", 0, "Set ReadTimeout")
	flag.DurationVar(&handlerTimeout, "handlerTimeout", 30*time.Second, "Set HandlerTimeout")
	flag.DurationVar(&writeTimeout, "writeTimeout", 0, "Set WriteTimeout")
	flag.DurationVar(&readLatency, "readLatency", 0, "Set read latency of the handler")
	flag.DurationVar(&writeLatency, "writeLatency", 0, "Set write latency of the handler")
	flag.DurationVar(&idleTimeout, "idleTimeout", 0, "Set IdleTimeout")
	flag.StringVar(&addrHost, "host", "", "Set the IP address to bind to")
	flag.UintVar(&addrPort, "port", 8888, "Set the TCP port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", addrHost, addrPort)

	handler := handler{readLatency: readLatency, writeLatency: writeLatency}

	srv := http.Server{
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		Handler:           http.TimeoutHandler(&handler, handlerTimeout, "Timeout\n"),
	}

	log.Printf("Listening at %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Server exited with: %s\n", err)
	}
}
