package server

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	//"github.com/rs/xlog"
	"net/http"
)

type Router struct {
	*mux.Router
	middleware alice.Chain
}

// NewRouter constructor
func NewRouter() *Router {
	return &Router{
		Router: mux.NewRouter(),
	}
}

func (r *Router) Use(middleware ...alice.Constructor) {
	r.middleware = r.middleware.Append(middleware...)
}

func (r *Router) AddController(controller Controller) {
	controller.Mount(r)
}

func (r *Router) AddRoute(path string, handler http.Handler) *mux.Route {
	return r.Handle(path, handler)
}

func (r *Router) AddRouteFunc(path string, handler http.HandlerFunc) *mux.Route {
	return r.Handle(path, handler)
}

func (r *Router) AddPrefixRoute(prefix string, handler http.Handler) *mux.Route {
	return r.PathPrefix(prefix).Handler(handler)
}

func (r *Router) AddPrefixRouteFunc(prefix string, handler http.HandlerFunc) *mux.Route {
	return r.PathPrefix(prefix).HandlerFunc(handler)
}
