package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("closing conn: %v\n", err)
		}
	}()

	input := inputService(os.Stdin, "stdin")
	server := inputService(conn, "server")

loop:
	for {
		select {
		case buf, ok := <-input:
			if !ok {
				break loop
			}
			_, err := fmt.Fprintf(conn, "Client: %s", buf)
			if err != nil {
				log.Fatalf("failed writing to conn: %v", err)
			}
		case buf, ok := <-server:
			if !ok {
				break loop
			}
			fmt.Printf("Server: %s", buf)
		}
	}

	fmt.Println("Bye!")
}

func inputService(rd io.Reader, name string) chan []byte {
	stream := bufio.NewReader(rd)
	output := make(chan []byte)

	go func() {
		defer close(output)
		fmt.Printf("Listening to %s\n", name)

		for {
			bytes, err := stream.ReadBytes('\n')
			if err == io.EOF {
				fmt.Printf("EOF in %s reader\n", name)
				return
			}
			if err != nil {
				fmt.Printf("ERROR in %s reader: %v\n", name, err)
				return
			}

			output <- bytes
		}
	}()

	return output
}
