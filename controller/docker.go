package controller

import (
	"github.com/euskadi31/docker-manager/docker"
	"github.com/euskadi31/docker-manager/server"
	"net/http/httputil"
)

type DockerController struct {
	proxy *httputil.ReverseProxy
}

func NewDockerController(host string) (*DockerController, error) {
	proxy, err := docker.NewProxy(host)
	if err != nil {
		return nil, err
	}

	return &DockerController{
		proxy: proxy,
	}, nil
}

func (c DockerController) Mount(r *server.Router) {
	r.AddPrefixRoute("/api/", c.proxy)
}
