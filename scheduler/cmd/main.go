package main

import (
	"log"
	"time"

	"github.com/munckymagik/gokb/scheduler"
)

func main() {
	scheduler := scheduler.New(time.Second*2, func(stopping bool) {
		log.Printf("tick: %v", stopping)
	})

	scheduler.Start()
	time.Sleep(time.Second * 10)
	scheduler.Stop()
}
