package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func main() {
	var (
		timeout         time.Duration
		dialerTimeout   time.Duration
		url             string
		maxIdleConns    int
		idleConnTimeout time.Duration
	)

	flag.DurationVar(&timeout, "timeout", 0, "Set http.Client.Timeout")
	flag.DurationVar(&dialerTimeout, "dialerTimeout", 30*time.Second, "Set net.Dialer.Timeout")
	flag.IntVar(&maxIdleConns, "maxIdleConns", 100, "Set http.Transport.MaxIdleConns")
	flag.DurationVar(&idleConnTimeout, "idleConnTimeout", 30*time.Second, "Set http.Transport.IdleConnTimeout")
	flag.StringVar(&url, "url", "http://localhost:8888/", "Set the URL to make requests to")
	flag.Parse()

	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   dialerTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:    maxIdleConns,
			IdleConnTimeout: idleConnTimeout,
		},
	}

	start := time.Now()
	body := strings.NewReader("Client OK")
	resp, err := client.Post(url, "text/plain", body)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// This is necessary or it may not be possible to reuse keep alive connections
	// in the connection pool.
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	log.Printf("duration:%v status:%d body:\"%s\"\n", time.Since(start), resp.StatusCode, buf)
}
