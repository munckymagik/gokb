package main

import (
	"bufio"
	"fmt"
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
		if err != nil {
			fmt.Printf("Stdin error: %v\n", err)
			break
		}

		_, err = conn.Write(request)
		if err != nil {
			fmt.Printf("writing request: %v\n", err)
			break
		}

		fmt.Println("Awaiting response ...")
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		response, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("awaiting response: %v\n", err)
			break
		}

		fmt.Printf("RESPONSE: %s", response)
	}

	fmt.Println("Bye!")
}
