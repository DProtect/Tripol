//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("READY")
	input.Scan()
	url := input.Text()
	input.Scan()
	count, err := strconv.Atoi(input.Text())
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	total := 0
	start := time.Now()
	for i := 0; i < count; i++ {
		reply, err := client.Get(url)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(reply.Body)
		reply.Body.Close()
		if err != nil || reply.StatusCode != 200 || len(data) != 512 {
			panic("unexpected response")
		}
		total += len(data)
	}
	fmt.Println(time.Since(start).Nanoseconds(), total)
	input.Scan()
}
