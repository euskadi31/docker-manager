// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/rs/xlog"
	"net/http"
)

// Server struct
type Server struct {
	proxy *Proxy
}

// NewServer create a Server
func NewServer() (*Server, error) {
	proxy, err := NewProxy(Config.Target)
	if err != nil {
		return nil, err
	}

	return &Server{
		proxy: proxy,
	}, nil
}

// Listen Server
func (s *Server) Listen() error {
	addr := fmt.Sprintf(":%d", Config.Port)

	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/api/", s.proxy.handle)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	xlog.Infof("Server running on %s", addr)

	return http.ListenAndServe(addr, nil)
}
