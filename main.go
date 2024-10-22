package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	LoadAverage         int
	MemoryCapacity      int
	MemoryUsage         int
	DiskCapacity        int
	DiskUsage           int
	NetworkCapacity     int
	NetworkUsage        int
	memoryUsagePercent  int
	diskUsagePercent    int
	networkUsagePercent int
}

func makeServerStats(val []string, ss ServerStats) ServerStats {
	var err error
	ss.LoadAverage, err = strconv.Atoi(val[0])
	if err != nil {
		fmt.Println(err)
	}
	ss.MemoryCapacity, err = strconv.Atoi(val[1])
	if err != nil {
		fmt.Println(err)
	}
	ss.MemoryUsage, err = strconv.Atoi(val[2])
	if err != nil {
		fmt.Println(err)
	}
	ss.DiskCapacity, err = strconv.Atoi(val[3])
	if err != nil {
		fmt.Println(err)
	}
	ss.DiskUsage, err = strconv.Atoi(val[4])
	if err != nil {
		fmt.Println(err)
	}
	ss.NetworkCapacity, err = strconv.Atoi(val[5])
	if err != nil {
		fmt.Println(err)
	}
	ss.NetworkUsage, err = strconv.Atoi(val[6])
	if err != nil {
		fmt.Println(err)
	}
	ss.memoryUsagePercent = int(float64(ss.MemoryUsage) / float64(ss.MemoryCapacity) * 100)
	ss.diskUsagePercent = int(float64(ss.DiskUsage) / float64(ss.DiskCapacity) * 100)
	ss.networkUsagePercent = int(float64(ss.NetworkUsage) / float64(ss.NetworkCapacity) * 100)
	return ss
}

func (ss ServerStats) checkMemoryUsagePercent() (err string, ok bool) {
	ok = true
	if ss.memoryUsagePercent > MemoryThreshold {
		err = fmt.Sprintf("Memory usage too high: %d%%", ss.memoryUsagePercent)
		ok = false
	}
	return
}

func (ss ServerStats) checkAvailableSpace() (err string, ok bool) {
	ok = true
	if ss.diskUsagePercent > DiskThreshold {
		availableSpace := (ss.DiskCapacity - ss.DiskUsage) / 1024 / 1024
		err = fmt.Sprintf("Free disk space is too low: %d Mb left", availableSpace)
		ok = false
	}
	return
}

func (ss ServerStats) checkavAilableBandwidth() (err string, ok bool) {
	ok = true
	if ss.networkUsagePercent > NetworkThreshold {
		availableBandwidth := (ss.NetworkCapacity - ss.NetworkUsage) / 1000 / 1000
		err = fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", availableBandwidth)
		ok = false
	}
	return
}

type ServerStatsList []ServerStats

func main() {
	c := http.Client{}
	ssl := ServerStatsList{}
	errCount := 0
	for i := 0; i < 100; i++ {
		resp, err := c.Get(ServerURL)
		if err != nil {
			errCount++
			fmt.Println(err)
		}
		defer resp.Body.Close()

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
		// fmt.Println(ss)
		ssl = append(ssl, ss)
	}
	// ss := makeServerStats(strings.Split("11,4915402826,1712029496,423323774247,409739069884,2482309012,365544533", ","), ServerStats{})
	// ss := makeServerStats(strings.Split("3,4915402826,2200880953,423323774247,113519465486,2482309012,403665858", ","), ServerStats{})
	// ss := makeServerStats(strings.Split("83,4915402826,4915402826,423323774247,397994209170,2482309012,554186051", ","), ServerStats{})
	// fmt.Println(ss)
	// ssl = append(ssl, ss)
	// fmt.Println(ssl)

	for _, v := range ssl {
		err, ok := v.checkMemoryUsagePercent()
		if !ok {
			fmt.Println(err)
		}
		err, ok = v.checkAvailableSpace()
		if !ok {
			fmt.Println(err)
		}
		err, ok = v.checkavAilableBandwidth()
		if !ok {
			fmt.Println(err)
		}
	}

}
