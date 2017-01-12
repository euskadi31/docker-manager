// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy struct
type Proxy struct {
	proxy *httputil.ReverseProxy
}

// NewProxy server
func NewProxy(target string) (*Proxy, error) {
	uri, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	/*proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
	})
	proxy.Director = func(r *http.Request) {
		log.Printf("%#v", r.URL)
		r.Host = uri.Host
		r.URL.Host = uri.Host
		r.URL.Scheme = uri.Scheme
	}*/

	return &Proxy{
		proxy: httputil.NewSingleHostReverseProxy(uri),
	}, nil
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.Replace(r.RequestURI, "/api/", "/", 1)
	p.proxy.ServeHTTP(w, r)
}

// see http://www.darul.io/post/2015-07-22_go-lang-simple-reverse-proxy
// url target: http://docker
// Auth: oauth2/jwt
