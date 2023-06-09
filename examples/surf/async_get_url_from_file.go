package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	start := time.Now()
	urlsChan := make(chan *surf.AsyncURL)

	file, err := os.Open("domains.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(urlsChan)

		for scanner.Scan() {
			domain := scanner.Text()
			domain = "http://" + strings.TrimSpace(domain)
			urlsChan <- surf.NewAsyncURL(domain)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opt := surf.NewOptions().
		DisableKeepAlive().
		MaxRedirects(3).
		Timeout(time.Second * 20).
		HTTP2(false).
		GetRemoteAddress().
		DNS("127.0.0.1:53")

	jobs, errors := surf.NewClient().
		SetOptions(opt).
		Async.WithContext(ctx).
		Get(urlsChan).
		Pool(100).
		Do()

	const limitBytes = 250000

	var counter int32
	for jobs != nil && errors != nil {
		atomic.AddInt32(&counter, 1)

		if counter%1000 == 0 {
			fmt.Printf(
				"number goroutines: %d started at: %s now: %s, urls counter: %d\n\n",
				runtime.NumGoroutine(),
				start.Format("2006-01-02 15:04:05"),
				time.Now().Format("2006-01-02 15:04:05"),
				counter,
			)
		}

		select {
		case job, ok := <-jobs:
			if !ok {
				jobs = nil
				continue
			}

			job.Body.Limit(limitBytes).Bytes()
			fmt.Println(job.RemoteAddress())
			// job.Body.Bytes()
		case _, ok := <-errors:
			if !ok {
				errors = nil
				continue
			}
		}
	}
}
