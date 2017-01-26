package controller

import (
	"github.com/euskadi31/docker-manager/server"
	"net/http"
	"strings"
)

type UiController struct {
}

func NewUiController() (*UiController, error) {
	return &UiController{}, nil
}

func (c UiController) Mount(r *server.Router) {
	r.AddPrefixRouteFunc("/ui/", c.StaticFile)
}

func (c UiController) StaticFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.Replace(r.URL.Path, "/ui/", "/", 1)

	extensions := []string{".js", ".css", ".map", ".ico"}
	for _, ext := range extensions {
		if strings.HasSuffix(r.URL.Path, ext) {
			http.ServeFile(w, r, "/opt/docker-manager/ui/"+filename)

			return
		}
	}

	http.ServeFile(w, r, "/opt/docker-manager/ui/index.html")
}
