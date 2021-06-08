// Copyright 2021 William Perron. All rights reserved. MIT License.
package config

import (
	"net/http"
	"time"

	_ "gopkg.in/yaml.v3"
)

// Config of the zombie process
type Config struct {
	// API configuration
	Api *Api `yaml:"api,omitempty"`

	// List of Targets
	Targets []Target `yaml:"targets"`
}

// API serving status and metrics info about the zombie process
type Api struct {
	Enabled bool   `yaml:"enabled"`
	Addr    string `yaml:"addr,omitempty"`
}

// Target to crawl
type Target struct {
	// URL to be requested
	Url string `yaml:"url"`

	// Name to print out in the log, defaults to URL if left empty
	Name string `yaml:"name,omitempty"`

	// Headers to add to the request
	Headers *http.Header `yaml:"headers,omitempty"`

	// Delay to wait between each request. This parameter is affected byt the
	// Jitter parameter
	Delay time.Duration `yaml:"delay"`

	// Jitter applied to the Delay between each request. Jitter is a modifier
	// applied in each direction so that a value of `0.2` means ±20%
	Jitter float64 `yaml:"jitter"`
}
