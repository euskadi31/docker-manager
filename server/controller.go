package server

type Controller interface {
	Mount(r *Router)
}
