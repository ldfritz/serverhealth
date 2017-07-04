package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	startTime = time.Date(2017, time.July, 3, 22, 15, 30, 0, time.UTC)
)

func uptime() {
	fmt.Printf("\nuptime: %s\n", time.Since(startTime))
}

func main() {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	go func() {
		for {
			select {
			case <-interrupts:
				uptime()
				os.Exit(0)
			}
		}
	}()
	for {
		_, err := http.Get("http://192.168.0.2")
		if err != nil {
			uptime()
			log.Fatal("Connection error")
		}
		log.Print(".")
		time.Sleep(1 * time.Minute)
	}
}
