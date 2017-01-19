package main

import (
	"context"
	"net/http"

	"github.com/asdine/storm"
	"github.com/docker/docker/client"
)

type contextKey int

const (
	stormKey contextKey = iota
	dockerKey
)

// NewDockerHandler middleware
func NewDockerHandler(c *client.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r != nil {
				r = r.WithContext(context.WithValue(r.Context(), dockerKey, c))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// NewStormHandler middleware
func NewStormHandler(db *storm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r != nil {
				r = r.WithContext(context.WithValue(r.Context(), stormKey, db))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// StormFromContext gets the strom out of the context.
// If not storm is stored in the context, a nil is returned.
func StormFromContext(ctx context.Context) *storm.DB {
	if ctx == nil {
		return nil
	}

	db, ok := ctx.Value(stormKey).(*storm.DB)
	if !ok {
		return nil
	}

	return db
}

// DockerFromContext gets the docker client out of the context.
// If not docker client is stored in the context, a nil is returned.
func DockerFromContext(ctx context.Context) *client.Client {
	if ctx == nil {
		return nil
	}

	dc, ok := ctx.Value(dockerKey).(*client.Client)
	if !ok {
		return nil
	}

	return dc
}
