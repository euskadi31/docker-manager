// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"strings"

	"github.com/pkg/errors"

	// "github.com/RangelReale/osin"
	"io/ioutil"

	"strconv"

	"github.com/asdine/storm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/euskadi31/docker-manager/docker"
	"github.com/euskadi31/docker-manager/entity"
	"github.com/euskadi31/docker-manager/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/justinas/alice"
	"github.com/rs/xlog"
)

type DockerLog struct {
	Type      string            `json:"Type"`
	Labels    map[string]string `json:"Labels"`
	Timestamp string            `json:"Timestamp"`
	IP        string            `json:"IP"`
	Message   string            `json:"Message"`
}

func NewDockerLog(b []byte) *DockerLog {
	// It is encoded on the first 8 bytes like this:
	//
	// header := [8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}
	//
	// `STREAM_TYPE` can be:
	//
	// -   0: stdin (will be written on stdout)
	// -   1: stdout
	// -   2: stderr
	//
	// `SIZE1, SIZE2, SIZE3, SIZE4` are the 4 bytes of
	// the uint32 size encoded as big endian.

	h := make([]byte, 8)
	buf := bytes.NewBuffer(b)
	buf.Read(h)

	var t string

	switch h[0] {
	case 0:
		t = "stdin"
	case 1:
		t = "stdout"
	case 2:
		t = "stderr"
	}

	//xlog.Debugf("Docker Log: %s", string(buf.Bytes()))

	logmsg := bytes.SplitN(buf.Bytes(), []byte(" - - "), 2)

	part := bytes.SplitN(logmsg[0], []byte(" "), 3)

	labels := make(map[string]string)

	l := string(part[1])

	if l != "" {
		items := strings.Split(l, ",")

		for _, val := range items {
			pair := strings.SplitN(val, "=", 2)

			labels[pair[0]] = pair[1]
		}
	}

	// 2017-01-20T05:50:42.047552194Z  10.0.1.3 - - [20/Jan/2017:05:50:41 +0000] "GET /logo.png HTTP/1.1" 200 13133 "http://localhost:8012/" "Mozilla
	// 2017-01-20T05:50:42.047490838Z com.docker.swarm.node.id=8dsrkozezn9j2lbvmciwtn2wf,com.docker.swarm.service.id=3hwqpkilg4saotx96jeiwgggg,com.docker.swarm.task.id=z037zkyr99u649qimsumnho0w 10.0.1.3 - - [20/Jan/2017:05:50:41 +0000] "GET / HTTP/1.1" 200 485 "-" "Mozilla

	return &DockerLog{
		Type:      t,
		Labels:    labels,
		Timestamp: string(part[0]),
		IP:        string(part[2]),
		Message:   string(logmsg[1]),
	}
}

var (
	PONG = []byte(`{"type":"pong"}`)
	PING = []byte(`{"type":"ping"}`)
)

const maxMessageSize = 10240

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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

	dctx := context.Background()

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

	//
	router.Handle("/ws/container/{id}/log", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		go func() {
			ws.SetReadLimit(maxMessageSize)

			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
					) {
						xlog.Error(err)
					}
					break
				}

				if bytes.Compare(msg, PING) == 0 {
					ws.WriteMessage(websocket.TextMessage, PONG)
				}

				time.Sleep(time.Millisecond * 500)
			}
		}()

		ctx := r.Context()

		dc := DockerFromContext(ctx)

		responseBody, err := dc.ContainerLogs(dctx, vars["id"], types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     true,
		})
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		defer responseBody.Close()

		scanner := bufio.NewScanner(responseBody)

		for {
			if scanner.Scan() {
				b, err := json.Marshal(NewDockerLog(scanner.Bytes()))
				if err != nil {
					continue
				}

				ws.WriteMessage(websocket.TextMessage, b)
			} else {
				time.Sleep(time.Millisecond * 1000)
			}
		}
	}))

	router.Handle("/ws/service/{name}/log", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		go func() {
			ws.SetReadLimit(maxMessageSize)

			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
					) {
						xlog.Error(err)
					}
					break
				}

				if bytes.Compare(msg, PING) == 0 {
					ws.WriteMessage(websocket.TextMessage, PONG)
				}

				time.Sleep(time.Millisecond * 500)
			}
		}()

		ctx := r.Context()

		dc := DockerFromContext(ctx)

		responseBody, err := dc.ServiceLogs(dctx, vars["name"], types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Tail:       "5",
			Follow:     true,
			// Details:    true,
		})
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		defer responseBody.Close()

		scanner := bufio.NewScanner(responseBody)

		for {
			if scanner.Scan() {
				b, err := json.Marshal(NewDockerLog(scanner.Bytes()))
				if err != nil {
					continue
				}

				ws.WriteMessage(websocket.TextMessage, b)
			} else {
				time.Sleep(time.Millisecond * 1000)
			}
		}
	}))

	router.Handle("/ws/events", middleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		go func() {
			ws.SetReadLimit(maxMessageSize)

			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
					) {
						xlog.Error(err)
					}
					break
				}

				if bytes.Compare(msg, PING) == 0 {
					ws.WriteMessage(websocket.TextMessage, PONG)
				}

				time.Sleep(time.Millisecond * 500)
			}
		}()

		dc := DockerFromContext(r.Context())

		eventq, errq := dc.Events(context.Background(), types.EventsOptions{})

		for {
			select {
			case event := <-eventq:
				b, err := json.Marshal(event)
				if err != nil {
					continue
				}
				ws.WriteMessage(websocket.TextMessage, b)
			case <-errq:
				return
			}
		}
	}))

	// Docker API proxy
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
