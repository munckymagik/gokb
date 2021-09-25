package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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

	inputReader := bufio.NewReaderSize(os.Stdin, 80)
	connReader := bufio.NewReaderSize(conn, 80)

	for {
		fmt.Print("Write a request and press Return: ")

		request, err := inputReader.ReadBytes('\n')
		if err == io.EOF {
			fmt.Println("Stdin got EOF")
			break
		}
		if err != nil {
			fmt.Printf("Stdin error: %v\n", err)
			break
		}

		fmt.Fprintf(conn, "REQUEST: %s", request)

		fmt.Println("Awaiting response ...")
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		response, err := connReader.ReadString('\n')
		if err == nil {
			fmt.Printf("RESPONSE: %s", response)
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("ERROR timeout awaiting response: %v\n", err)
		} else if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
			fmt.Printf("ERROR temporary error awaiting response: %v\n", err)
		} else if err == io.EOF {
			fmt.Println("conn got EOF")
			break
		} else {
			fmt.Printf("conn error: %v\n", err)
			break
		}

	}

	fmt.Println("Bye!")
}
