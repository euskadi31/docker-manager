// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caarlos0/env"
	"github.com/rs/xlog"
	"log"
)

// Configuration struct
type Configuration struct {
	Debug  bool   `env:"DEBUG" envDefault:"false"`
	Port   int    `env:"PORT" envDefault:"8080"`
	Target string `env:"TARGET" envDefault:"http://docker"`
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
