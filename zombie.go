// Copyright 2021 William Perron. All rights reserved. MIT License.

// Command zombie is a natural load generator to simulate real-life traffic
// on a system.
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultDelay  = 10000 * time.Millisecond
	defaultJitter = 0.2
)

type Target struct {
	Url    string
	Header http.Header
	Delay  time.Duration
	Jitter float64
}

func main() {
	targets := []Target{
		// {
		// 	Url: "https://example.org",
		// 	Header: map[string][]string{
		// 		"Authorization": {"Bearer eyJ0eXAiOiEXAMPLE"},
		// 	},
		// 	Delay:  1000 * time.Millisecond,
		// 	Jitter: 0.2,
		// },
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	out := make(chan string)

	for _, t := range targets {
		go ping(t, out)
	}

	for message := range out {
		log.Println(message)
	}
}

func ping(t Target, out chan string) {
	u, err := url.Parse(t.Url)
	if err != nil {
		log.Fatalf("unable to parse URL %s", t.Url)
	}

	req := http.Request{
		Method: "GET",
		URL:    u,
	}

	if len(t.Header) > 0 {
		req.Header = t.Header
	}

	for {
		delay := float64(t.Delay)
		if delay == 0.0 {
			delay = float64(defaultDelay)
		}

		jitter := float64(t.Jitter)
		if jitter == 0.0 {
			jitter = float64(defaultJitter)
		}

		time.Sleep(time.Duration(delay * (1 + (jitter * (rand.Float64()*2 - 1)))))

		res, err := http.DefaultClient.Do(&req)
		if err != nil {
			out <- fmt.Sprintf("error: %s", err)
		} else {
			out <- fmt.Sprintf("GET %s %s", res.Status, u.String())
		}
	}
}
