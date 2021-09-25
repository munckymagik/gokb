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
	listener, err := net.Listen("tcp", "localhost:8088")
	if err != nil {
		log.Fatalf("creating listener: %v", err)
	}

	fmt.Println("I'm listening ...")
	connID := 0
	for {
		conn, err := listener.Accept()
		connID += 1
		if err != nil {
			log.Fatalf("accepting: %v", err)
		}

		go handleConn(connID, conn)
	}
}

func handleConn(id int, conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("%d: closing conn: %v\n", id, err)
		}
	}()

	inputReader := bufio.NewReaderSize(os.Stdin, 80)
	connReader := bufio.NewReaderSize(conn, 80)

	for {
		fmt.Printf("%d: Awaiting request ...\n", id)

		// Set an idle timeout
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err := connReader.Peek(1)
		if err != nil {
			fmt.Printf("idling: %v\n", err)
			break
		}

		// Set a read timeout
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		var request string
		request, err = connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("%d while awaiting request: %v\n", id, err)
			break
		}

		fmt.Printf("%d REQUEST: %s", id, request)
		fmt.Printf("%d Write a response and press Return: ", id)

		response, err := inputReader.ReadBytes('\n')
		if err == io.EOF {
			fmt.Printf("%d Stdin got EOF\n", id)
			break
		}
		if err != nil {
			fmt.Printf("%d Stdin error: %v\n", id, err)
			break
		}

		// Set a write timeout
		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

		_, err = conn.Write(response)
		if err != nil {
			fmt.Printf("writing response: %v\n", err)
			break
		}
	}

	fmt.Printf("%d: Bye!\n", id)
}
