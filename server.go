// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"strings"

	"github.com/euskadi31/docker-manager/docker"
	"github.com/gorilla/mux"
	"github.com/rs/xlog"
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

	//router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("/opt/docker-manager/ui/"))))

	router.PathPrefix("/ui/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := strings.Replace(r.URL.Path, "/ui/", "/", 1)

		extensions := []string{".js", ".css", ".map", ".ico"}
		for _, ext := range extensions {
			if strings.HasSuffix(r.URL.Path, ext) {
				http.ServeFile(w, r, "/opt/docker-manager/ui/"+filename)

				return
			}
		}

		http.ServeFile(w, r, "/opt/docker-manager/ui/index.html")
	})

	xlog.Infof("Server running on %s", addr)

	return http.ListenAndServe(addr, router)
}
