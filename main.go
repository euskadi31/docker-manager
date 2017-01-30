// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/docker/docker/client"
	"github.com/euskadi31/docker-manager/controller"
	"github.com/euskadi31/docker-manager/database"
	"github.com/euskadi31/docker-manager/docker"
	"github.com/euskadi31/docker-manager/server"
	"github.com/rs/xlog"
	"net/http"
)

func main() {
	addr := fmt.Sprintf(":%d", Config.Port)

	db, err := storm.Open("/var/lib/docker-manager/docker-manager.db")
	if err != nil {
		xlog.Fatal(err)
	}

	dc, err := client.NewEnvClient()
	if err != nil {
		xlog.Fatal(err)
	}

	router := server.NewRouter()
	router.Use(database.NewHandler(db))
	router.Use(docker.NewHandler(dc))

	dockerController, err := controller.NewDockerController(Config.DockerHost)
	if err != nil {
		xlog.Fatal(err)
	}

	uiController, err := controller.NewUiController()
	if err != nil {
		xlog.Fatal(err)
	}

	wsController, err := controller.NewWsController()
	if err != nil {
		xlog.Fatal(err)
	}

	apiController, err := controller.NewApiController()
	if err != nil {
		xlog.Fatal(err)
	}

	router.AddController(apiController)
	router.AddController(dockerController)
	router.AddController(uiController)
	router.AddController(wsController)

	xlog.Infof("Server running on %s", addr)

	xlog.Error(http.ListenAndServe(addr, router))
}
