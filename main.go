package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Server struct {
	StartTime time.Time
	Address   string
}

func LogUptime(startTime time.Time) {
	log.Printf("uptime: %s\n", time.Since(startTime) - (4 * time.Hour))
}

func (s Server) Check() {
	resp, err := http.Get(s.Address)
	if err != nil {
		LogUptime(s.StartTime)
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("%s %v", s.Address, resp.Status)
}

func main() {
	addr := os.Args[1]
	if !strings.Contains(addr, ":") {
		addr = "http://" + addr
	}

	format := "2006-01-02T15:04"
	startTime, err := time.Parse(format, os.Args[2])
	if err != nil {
		log.Printf("timestamp parsing: %v", err)
	}
	LogUptime(startTime)

	server := Server{startTime, addr}

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt)

	quick := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-quick.C:
			go server.Check()
		case <-interrupted:
			LogUptime(startTime)
			os.Exit(0)
		}
	}
}
