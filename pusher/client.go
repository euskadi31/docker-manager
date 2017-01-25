// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pusher

import (
	"fmt"
	"github.com/euskadi31/docker-manager/pusher/message"
	"github.com/gorilla/websocket"
	"github.com/rs/xlog"
	"github.com/satori/go.uuid"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

// Client struct
type Client struct {
	ID     string
	ws     *websocket.Conn
	pusher *Pusher
}

// NewClient constructor
func NewClient(ws *websocket.Conn, pusher *Pusher) *Client {
	return &Client{
		ID:     uuid.NewV4().String(),
		ws:     ws,
		pusher: pusher,
	}
}

// Write event to client
func (c *Client) Write(event *message.Event) {
	b, err := message.Encode(event)
	if err != nil {
		xlog.Error(err)

		return
	}

	if err := c.ws.WriteMessage(websocket.TextMessage, b); err != nil {
		xlog.Error(err)

		c.Close()
	}
}

// Close client
func (c *Client) Close() {

}

func (c *Client) process(event *message.Event) {
	switch event.Type {
	case "message":
		//c.pusher.Write(event)

	case "subscribe":
		c.pusher.JoinChannel(event.Channel, c)

	case "unsubscribe":
		c.pusher.LeaveChannel(event.Channel, c)

	case "ping":
		c.Write(&message.Event{
			Type: "pong",
		})

	default:
		err := fmt.Sprintf("Unsupported %s event type", event.Type)

		xlog.Error(err)

		c.Write(&message.Event{
			Type: "error",
			Data: err,
		})
	}
}

// Listen func
func (c *Client) Listen() {
	c.ws.SetReadLimit(maxMessageSize)

	for {
		_, msg, err := c.ws.ReadMessage()
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

		event, err := message.Decode(msg)
		if err != nil {
			xlog.Error(err)
		} else {
			c.process(event)
		}

		time.Sleep(time.Millisecond * 100)
	}
}
