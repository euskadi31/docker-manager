// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

func main() {
	server, err := NewServer()
	if err != nil {
		log.Panic(err)
	}

	server.Listen()
}
