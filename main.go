// Copyright 2021 William Perron. All rights reserved. MIT License.

// Command zombie is a natural load generator to simulate real-life traffic
// on a system.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wperron/zombie/api"
	"github.com/wperron/zombie/client"
	"github.com/wperron/zombie/config"
)

var (
	configPath = flag.String("config", "", "The location of the config file.")
	// TODO(wperron) add verbose and quiet options
)

func init() {
	if client.DefaultPinger == nil {
		log.Fatal("default pinger is nil")
	}
}

func main() {
	// Set up channel on which to send termination signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	// Parse command line args
	flag.Parse()

	// Load the configuration file
	conf, err := config.LoadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	printSummary(*conf)

	// Start the API if enabled
	if conf.Api != nil && conf.Api.Enabled {
		go api.Serve(conf.Api.Addr)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	out := make(chan string)

	for _, t := range conf.Targets {
		go client.DefaultPinger.Ping(t, out)
	}

	go func() {
		for message := range out {
			log.Println(message)
		}
	}()

	// Block until a signal is received.
	s := <-sigs
	log.Println("Got signal:", s)
}

func printSummary(c config.Config) {
	fmt.Println("Zombie started")
	if c.Api != nil && c.Api.Enabled {
		fmt.Printf("API enabled on %s\n", c.Api.Addr)
	}
	for _, t := range c.Targets {
		if t.Name != "" {
			fmt.Printf("target name: %s at %s, base delay: %d ms, jitter: %f\n", t.Name, t.Url, t.Delay.Milliseconds(), t.Jitter)
		} else {
			fmt.Printf("target name: %s, base delay: %dms, jitter: %f\n", t.Url, t.Delay, t.Jitter)
		}
	}
	fmt.Println("")
}
