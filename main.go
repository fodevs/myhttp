package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const defaultParallel = 10

func main() {
	maxP := flag.Int("parallel", defaultParallel, "max number of parallel requests")
	flag.Parse()
	if *maxP <= 0 {
		*maxP = defaultParallel
	}

	ch := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for _, arg := range flag.Args() {
			u, err := url.Parse(arg)
			if err != nil {
				fmt.Printf("\nParse link `%s` error: %v", arg, err)
				continue
			}
			ch <- "http://" + u.Host + u.Path
		}
		close(ch)
	}()

	var wg sync.WaitGroup
	wg.Add(*maxP)
	for i := 0; i < *maxP; i++ {
		go func() {
			defer wg.Done()
			hasher := md5.New()

			for link := range ch {
				h, err := computeHash(ctx, hasher, link)
				if err != nil {
					fmt.Println(link, err)
					continue
				}
				fmt.Println(link, h)
				hasher.Reset()
			}
		}()
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		fmt.Printf("got signal %s", (<-sigc).String())
		cancel()
	}()
	wg.Wait()
}

func computeHash(ctx context.Context, hasher hash.Hash, link string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(hasher, resp.Body)
	if err != nil {
		return "", fmt.Errorf("do copy: %w", err)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
