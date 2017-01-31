package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/euskadi31/docker-manager/database"
	"github.com/euskadi31/docker-manager/entity"
	"github.com/euskadi31/docker-manager/server"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/xlog"
)

type ApiController struct {
}

func NewApiController() (*ApiController, error) {
	return &ApiController{}, nil
}

func (c ApiController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/registries", c.GetRegistriesHandler).Methods("GET")
	r.AddRouteFunc("/api/registries", c.PostRegistryHandler).Methods("POST")
	r.AddRouteFunc("/api/registries/{id:[0-9]+}", c.PutRegistryHandler).Methods("PUT")
	r.AddRouteFunc("/api/registries/{id:[0-9]+}", c.DeleteRegistryHandler).Methods("DELETE")
	r.AddRouteFunc("/api/registries/{id:[0-9]+}/repositories", c.GetRegistryRepositoriesHandler).Methods("GET")
}

// GetRegistriesHandler /api/registries
func (c ApiController) GetRegistriesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	var registries []entity.Registry
	if err := db.All(&registries); err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, registries)
}

// PostRegistryHandler /api/registries
func (c ApiController) PostRegistryHandler(w http.ResponseWriter, r *http.Request) {
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
	db, err := database.FromContext(ctx)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	if err := db.Save(&registry); err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, registry)
}

// PutRegistryHandler /api/registries/{id:[0-9]+}
func (c ApiController) PutRegistryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	xlog.Infof("ID:", vars["id"])
}

// DeleteRegistryHandler /api/registries/{id:[0-9]+}
func (c ApiController) DeleteRegistryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	xlog.Infof("ID:", vars["id"])
}

// GetRegistryRepositoriesHandler /api/registries/{id:[0-9]+}/repositories
func (c ApiController) GetRegistryRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		server.NotFoundFailure(w, r)

		return
	}

	db, err := database.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

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
}
