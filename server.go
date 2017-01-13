// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/euskadi31/docker-manager/docker"
	"github.com/gorilla/mux"
	"github.com/rs/xlog"
	"net/http"
	"net/http/httputil"
)

// Server struct
type Server struct {
	proxy *httputil.ReverseProxy
}

// NewServer create a Server
func NewServer() (*Server, error) {
	proxy, err := docker.New(Config.DockerHost)
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

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler).Methods("GET", "HEAD")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
	router.PathPrefix("/api/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.proxy.ServeHTTP(w, r)
	})

	xlog.Infof("Server running on %s", addr)

	return http.ListenAndServe(addr, router)
}
