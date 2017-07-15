package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Server represents a server that can receive a GET request at Address
// and has been up since StartTime.
type Server struct {
	StartTime time.Time
	Address   string
}

// LogUptime simply logs the time since the start time.
func LogUptime(startTime time.Time) {
	log.Print(time.Since(startTime) - 4*time.Hour)
}

// Check will send a GET request to the Server and log the response
// status.  If the GET returns an error, that error and the time since
// the server's start time will be reported.
func (s Server) Check() {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(s.Address)
	if err != nil {
		LogUptime(s.StartTime)
		log.Printf("GET error: %v", err)
		return
	}
	log.Println(s.Address, resp.Status)
}

func main() {
	addr := os.Args[1]
	if !strings.Contains(addr, ":") {
		addr = "http://" + addr
	}

	format := "2006-01-02T15:04"
	startTime, err := time.Parse(format, os.Args[2])
	if err != nil {
		log.Printf("time parsing error: %v", err)
	}
	LogUptime(startTime)

	server := Server{startTime, addr}

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt)

	quick := time.NewTicker(15 * time.Second)
	defer quick.Stop()

topLoop:
	for {
		select {
		case <-quick.C:
			go server.Check()
		case <-interrupted:
			LogUptime(startTime)
			break topLoop
		}
	}
}
