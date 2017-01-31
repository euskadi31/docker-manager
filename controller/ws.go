package controller

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/euskadi31/docker-manager/docker"
	"github.com/euskadi31/docker-manager/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/xlog"
)

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

type WsController struct {
}

func NewWsController() (*WsController, error) {
	return &WsController{}, nil
}

func (c WsController) Mount(r *server.Router) {
	r.AddRouteFunc("/ws/events", c.EventHandler)
	r.AddRouteFunc("/ws/containers/{id}/logs", c.ContainerLogHandler)
	r.AddRouteFunc("/ws/services/{name}/logs", c.ServiceLogHandler)
}

// ServiceLogHandler /ws/containers/{id}/logs
func (c WsController) ContainerLogHandler(w http.ResponseWriter, r *http.Request) {
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

	dc, err := docker.FromContext(ctx)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	responseBody, err := dc.ContainerLogs(context.Background(), vars["id"], types.ContainerLogsOptions{
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
			msg, err := docker.ParseLog(scanner.Bytes())
			if err != nil {
				continue
			}

			b, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			ws.WriteMessage(websocket.TextMessage, b)
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

// ServiceLogHandler /ws/services/{name}/logs
func (c WsController) ServiceLogHandler(w http.ResponseWriter, r *http.Request) {
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

	dc, err := docker.FromContext(ctx)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	responseBody, err := dc.ServiceLogs(context.Background(), vars["name"], types.ContainerLogsOptions{
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
			msg, err := docker.ParseLog(scanner.Bytes())
			if err != nil {
				continue
			}

			b, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			ws.WriteMessage(websocket.TextMessage, b)
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

// EventHandler /ws/events
func (c WsController) EventHandler(w http.ResponseWriter, r *http.Request) {
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

	dc, err := docker.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

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
}
