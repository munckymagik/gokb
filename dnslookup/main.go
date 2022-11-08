package main

import (
	"fmt"
	"net"
	"os"
	"path"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide at least hostname to lookup")
		fmt.Printf("USAGE %s hostname[, hostname]...\n", path.Base(os.Args[0]))
		return
	}

	for _, host := range args {
		fmt.Printf("Looking up IP for %s\n", host)

		addrs, err := net.LookupHost(host)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if len(addrs) == 0 {
			fmt.Println("No results found")
		}

		for _, addr := range addrs {
			fmt.Println(addr)
		}

		fmt.Println()
	}
}
