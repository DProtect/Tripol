//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func receive(client *http.Client, url string, size int, path string) int {
	reply, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer reply.Body.Close()
	if reply.StatusCode != 200 {
		panic("unexpected response status")
	}
	if path == "" {
		data, err := io.ReadAll(reply.Body)
		if err != nil || len(data) != size {
			panic("unexpected response")
		}
		return len(data)
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	written, err := io.Copy(file, reply.Body)
	syncErr := file.Sync()
	if closeErr := file.Close(); err != nil || syncErr != nil || closeErr != nil || written != int64(size) {
		panic("unexpected download")
	}
	return int(written)
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	line := func() string {
		if !input.Scan() {
			panic("missing input")
		}
		return input.Text()
	}
	integer := func() int {
		value, err := strconv.Atoi(line())
		if err != nil {
			panic(err)
		}
		return value
	}
	url, count, mode, width, size, directory := line(), integer(), line(), integer(), integer(), line()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 4
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	fmt.Println("READY")
	line()
	total := 0
	start := time.Now()
	if mode == "concurrent" {
		done := make(chan int, width)
		for start := 0; start < count; start += width {
			for i := 0; i < width; i++ {
				go func() { done <- receive(client, url, size, "") }()
			}
			for i := 0; i < width; i++ {
				total += <-done
			}
		}
	} else {
		for i := 0; i < count; i++ {
			path := ""
			if mode == "stream" {
				path = filepath.Join(directory, fmt.Sprintf("%d.bin", i))
			}
			total += receive(client, url, size, path)
		}
	}
	fmt.Println(time.Since(start).Nanoseconds(), total)
	input.Scan()
}
