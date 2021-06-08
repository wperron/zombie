// Copyright 2021 William Perron. All rights reserved. MIT License.
package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func Load(s string) (*Config, error) {
	cfg := &Config{}

	if err := yaml.Unmarshal([]byte(s), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFile(fp string) (*Config, error) {
	bs, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %s", fp, err)
	}

	cfg, err := Load(string(bs))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file: %s", err)
	}
	return cfg, nil
}
