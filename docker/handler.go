package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
	"net/http"
)

type key int

const (
	dockerKey key = iota
)

var (
	errContextEmpry = errors.New("Context is nil")
	errNotFound     = errors.New("Docker is not found in the Context")
)

// NewHandler docker
func NewHandler(c *client.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r != nil {
				r = r.WithContext(context.WithValue(r.Context(), dockerKey, c))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// FromContext gets the docker client out of the context.
// If not docker client is stored in the context, a nil is returned.
func FromContext(ctx context.Context) (*client.Client, error) {
	if ctx == nil {
		return nil, errContextEmpry
	}

	dc, ok := ctx.Value(dockerKey).(*client.Client)
	if !ok {
		return nil, errNotFound
	}

	return dc, nil
}
