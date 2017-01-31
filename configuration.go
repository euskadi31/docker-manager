// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/rs/xlog"
)

// Configuration struct
type Configuration struct {
	Debug      bool   `env:"DEBUG" envDefault:"false"`
	Port       int    `env:"PORT" envDefault:"8080"`
	DockerHost string `env:"DOCKER_HOST" envDefault:"unix:///var/run/docker.sock"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
}

var (
	// Config object
	Config *Configuration
)

func init() {
	Config = &Configuration{}

	if err := env.Parse(Config); err != nil {
		log.Fatal(err)
	}

	xlog.Infof("Config: %#v", Config)
}
