// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	ascii = `
___ _ _  _ _   _ ___  ____ _    ____ _  _ ____ ____ ____ 
 |  | |\ |  \_/  |__] |__| |    |__| |\ | |    |___ |__/ 
 |  | | \|   |   |__] |  | |___ |  | | \| |___ |___ |  \                                        
`
)

// Config configuration details of balancer
type Config struct {
	SSLCertificateKey string      `yaml:"ssl_certificate_key"`
	Location          []*Location `yaml:"location"`
	Schema            string      `yaml:"schema"`
	Port              int         `yaml:"port"`
	SSLCertificate    string      `yaml:"ssl_certificate"`
}

// Location routing details of balancer
type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string) (*Config, error) {
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Print print config details
func (c *Config) Print() {
	fmt.Printf("%s\nSchema: %s\nPort: %d\nLocation:\n", ascii, c.Schema, c.Port)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxyPass: %s\n\tMode: %s\n",
			l.Pattern, l.ProxyPass, l.BalanceMode)
	}
}

// Validation verify the configuration details of the balancer
func (c *Config) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return errors.New(fmt.Sprintf("the schema \"%s\" not supported", c.Schema))
	}
	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}
	if c.Schema == "https" && (len(c.SSLCertificate) == 0 || len(c.SSLCertificateKey) == 0) {
		return errors.New("the https proxy requires ssl_certificate_key and ssl_certificate")
	}
	return nil
}
