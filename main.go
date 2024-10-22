package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	RequestTimout     = 2 * time.Second
	RequestsFrequency = 300 * time.Millisecond
	ErrorThreshold    = 3
	ServerURL         = "http://srv.msk01.gigacorp.local/_stats"
	LoadThreshold     = 30
	MemoryThreshold   = 80
	DiskThreshold     = 90
	NetworkThreshold  = 90
)

type ServerStats struct {
	LoadAverage     int
	MemoryCapacity  int
	MemoryUsage     int
	DiskCapacity    int
	DiskUsage       int
	NetworkCapacity int
	NetworkUsage    int
}

func makeServerStats(val []string, ss ServerStats) ServerStats {
	return ss
}

type ServerStatsList []ServerStats

func main() {
	c := http.Client{}
	ssl := ServerStatsList{}
	errCount := 0
	for i := 0; i < 3; i++ {
		resp, err := c.Get(ServerURL)
		if err != nil {
			errCount++
			fmt.Println(err)
		}
		if resp.StatusCode != http.StatusOK {
			errCount++
			fmt.Printf("failed to send request %s\n", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errCount++
			fmt.Printf("failed to parse response %s\n", err)
			continue
		}
		ss := makeServerStats(strings.Split(string(body), ","), ServerStats{})
		fmt.Println(ss)
		ssl = append(ssl, ss)
	}
	fmt.Println(ssl)
}
