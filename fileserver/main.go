package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Println("Missing required arguments")
		log.Printf("USAGE %s hostname:port rootFolder\n", path.Base(os.Args[0]))
		return
	}

	address := args[0]
	rootFolder := args[1]

	srv := http.Server{
		Addr:    address,
		Handler: http.FileServer(http.Dir(rootFolder)),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down ...")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Server starting: %v\n", address)
	if result := srv.ListenAndServe(); result != http.ErrServerClosed {
		log.Printf("Server exited with: %v\n", result)
	}

	<-idleConnsClosed

	log.Println("Done.")
}
