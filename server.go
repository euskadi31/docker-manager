// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"strings"

	"github.com/pkg/errors"

	// "github.com/RangelReale/osin"
	"io/ioutil"

	"strconv"

	"github.com/asdine/storm"
	"github.com/docker/docker/client"
	"github.com/euskadi31/docker-manager/docker"
	"github.com/euskadi31/docker-manager/entity"
	"github.com/euskadi31/docker-manager/server"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/xlog"
)

// Server struct
type Server struct {
	proxy *httputil.ReverseProxy
	// oauth2 *osin.Server
	db *storm.DB
	dc *client.Client
}

// NewServer create a Server
func NewServer() (*Server, error) {
	// oauth2 := osin.NewServer(osin.NewServerConfig(), &OAuthStorage{})

	proxy, err := docker.NewProxy(Config.DockerHost)
	if err != nil {
		return nil, err
	}

	db, err := storm.Open("/var/lib/docker-manager/docker-manager.db")
	if err != nil {
		return nil, err
	}

	dc, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Server{
		proxy: proxy,
		db:    db,
		dc:    dc,
		// oauth2: oauth2,
	}, nil
}

// Listen Server
func (s *Server) Listen() error {
	defer s.db.Close()
	addr := fmt.Sprintf(":%d", Config.Port)

	middleware := alice.New(
		NewStormHandler(s.db),
		NewDockerHandler(s.dc),
	)

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler).Methods("GET", "HEAD")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	router.Handle("/api/registries", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		db := StormFromContext(r.Context())

		var registries []entity.Registry
		if err := db.All(&registries); err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)
		}

		server.JSON(w, http.StatusOK, registries)
	})).Methods("GET")

	router.Handle("/api/registries", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var registry entity.Registry

		if err := json.NewDecoder(r.Body).Decode(&registry); err != nil {
			server.FailureFromError(w, http.StatusBadRequest, err)

			return
		}
		defer r.Body.Close()

		/*dc := DockerFromContext(ctx)

		auth, err := dc.RegistryLogin(ctx, types.AuthConfig{
			Username:      registry.Username,
			Password:      registry.Password,
			ServerAddress: "https://" + registry.Server + "/v2/",
		})
		if err != nil {
			server.FailureFromError(w, http.StatusBadRequest, err)

			return
		}

		xlog.Debugf("Auth Registry: %#v", auth)
		*/
		db := StormFromContext(ctx)

		if err := db.Save(&registry); err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		server.JSON(w, http.StatusCreated, registry)
	})).Methods("POST")

	router.Handle("/api/registries/{id:[0-9]+}", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		xlog.Infof("ID:", vars["id"])
	})).Methods("PUT")

	router.Handle("/api/registries/{id:[0-9]+}", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		xlog.Infof("ID:", vars["id"])
	})).Methods("DELETE")

	router.Handle("/api/registries/{id:[0-9]+}/repositories", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			server.NotFoundFailure(w, r)

			return
		}

		db := StormFromContext(r.Context())

		var registry entity.Registry

		if err := db.One("ID", ID, &registry); err != nil {
			server.FailureFromError(w, http.StatusNotFound, errors.Wrapf(err, "Cannot find registry by ID: %d", ID))

			return
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/_catalog", registry.Server), nil)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		// Add header with json of username and password
		req.SetBasicAuth(registry.Username, registry.Password)

		httpClient := &http.Client{}

		resp, err := httpClient.Do(req)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		xlog.Debugf("Response: %s", string(b))

		//json.NewDecoder(req.Body).Decode(&)

		xlog.Infof("ID:", vars["id"])
	})).Methods("GET")

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

	// Access token endpoint
	/*router.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		resp := s.oauth2.NewResponse()
		defer resp.Close()

		if ar := s.oauth2.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			s.oauth2.FinishAccessRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	}).Methods("POST")
	*/

	xlog.Infof("Server running on %s", addr)

	return http.ListenAndServe(addr, router)
}
