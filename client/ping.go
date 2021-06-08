// Copyright 2021 William Perron. All rights reserved. MIT License.
package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/wperron/zombie/config"
)

var (
	defaultDelay  = 10000 * time.Millisecond
	defaultJitter = 0.2
	DefaultPinger = &pinger{
		client: http.DefaultClient,
	}
)

type Pinger interface {
	Ping(config.Target, chan<- string)
}

type pinger struct {
	client *http.Client
}

func NewPinger(c http.Client) *pinger {
	return &pinger{
		client: &c,
	}
}

func (p *pinger) Ping(t config.Target, out chan<- string) {
	u, err := url.Parse(t.Url)
	if err != nil {
		log.Fatalf("unable to parse URL %s", t.Url)
	}

	req := http.Request{
		Method: "GET",
		URL:    u,
	}

	if t.Headers != nil && len(*t.Headers) > 0 {
		req.Header = *t.Headers
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

		res, err := p.client.Do(&req)
		if err != nil {
			out <- fmt.Sprintf("error: %s", err)
		} else {
			prefix := t.Name
			if prefix == "" {
				prefix = t.Url
			}
			out <- fmt.Sprintf("[%s] GET %s %s", prefix, res.Status, u.String())
		}
	}
}
