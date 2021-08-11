// Copyright 2021 William Perron. All rights reserved. MIT License.
package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type Result struct {
	Name       string
	Method     string
	Status     int
	StatusText string
	URL        string
	Latency    int
	TraceID    string
}

func NewPinger(c http.Client) *pinger {
	return &pinger{
		client: &c,
	}
}

func NewInstrumentedClient() *http.Client {
	client := &http.Client{}
	*client = *http.DefaultClient
	client.Timeout = 1 * time.Second

	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "client_in_flight_requests",
		Help: "A gauge of in-flight requests for the wrapped client.",
	})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "client_api_requests_total",
			Help: "A counter for requests from the wrapped client.",
		},
		[]string{"code", "method"},
	)

	// dnsLatencyVec uses custom buckets based on expected dns durations.
	// It has an instance label "event", which is set in the
	// DNSStart and DNSDonehook functions defined in the
	// InstrumentTrace struct below.
	dnsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_duration_seconds",
			Help:    "Trace dns latency histogram.",
			Buckets: []float64{.005, .01, .025, .05},
		},
		[]string{"event"},
	)

	// tlsLatencyVec uses custom buckets based on expected tls durations.
	// It has an instance label "event", which is set in the
	// TLSHandshakeStart and TLSHandshakeDone hook functions defined in the
	// InstrumentTrace struct below.
	tlsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tls_duration_seconds",
			Help:    "Trace tls latency histogram.",
			Buckets: []float64{.05, .1, .25, .5},
		},
		[]string{"event"},
	)

	// histVec has no labels, making it a zero-dimensional ObserverVec.
	histVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(counter, tlsLatencyVec, dnsLatencyVec, histVec, inFlightGauge)

	// Define functions for the available httptrace.ClientTrace hook
	// functions that we want to instrument.
	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_start").Observe(t)
		},
		DNSDone: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_done").Observe(t)
		},
		TLSHandshakeStart: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_start").Observe(t)
		},
		TLSHandshakeDone: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_done").Observe(t)
		},
	}

	// Wrap the default RoundTripper with middleware.
	roundTripper := promhttp.InstrumentRoundTripperInFlight(inFlightGauge,
		promhttp.InstrumentRoundTripperCounter(counter,
			promhttp.InstrumentRoundTripperTrace(trace,
				promhttp.InstrumentRoundTripperDuration(histVec, http.DefaultTransport),
			),
		),
	)

	// Set the RoundTripper on our client.
	client.Transport = roundTripper
	return client
}

func NewInstrumentedPinger() *pinger {
	client := NewInstrumentedClient()
	return &pinger{
		client: client,
	}
}

func (p *pinger) Ping(t config.Target, out chan<- Result, e chan<- error) {
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
		delay := float64(t.Duration())
		if delay == 0.0 {
			delay = float64(defaultDelay)
		}

		jitter := float64(t.Jitter)
		if jitter == 0.0 {
			jitter = float64(defaultJitter)
		}

		time.Sleep(time.Duration(delay * (1 + (jitter * (rand.Float64()*2 - 1)))))

		start := time.Now()
		res, err := p.client.Do(&req)
		lat := time.Since(start)
		if err != nil {
			e <- fmt.Errorf("error: %s", err)
		} else {
			out <- Result{
				Name:       t.Name,
				Method:     "GET",
				Status:     res.StatusCode,
				StatusText: res.Status,
				URL:        t.Url,
				Latency:    int(lat.Milliseconds()),
				TraceID:    res.Header.Get(t.TraceHeader),
			}
		}
	}
}
