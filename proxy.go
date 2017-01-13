// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/docker/go-connections/sockets"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy struct
type Proxy struct {
	proxy *httputil.ReverseProxy
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// NewProxy server
func NewProxy(host string) (*Proxy, error) {
	sock, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	target, err := url.Parse("http://docker")
	if err != nil {
		return nil, err
	}

	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(
			target.Path,
			strings.Replace(req.URL.Path, "/api/", "/", 1),
		)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	transport := new(http.Transport)
	sockets.ConfigureTransport(transport, sock.Scheme, sock.Path)

	proxy := &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}

	return &Proxy{
		proxy: proxy,
	}, nil
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

// see http://www.darul.io/post/2015-07-22_go-lang-simple-reverse-proxy
// url target: http://docker
// Auth: oauth2/jwt
