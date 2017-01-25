// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pusher

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/xlog"
)

// ErrorHandler func
type ErrorHandler func(w http.ResponseWriter, code int, err error)

// Pusher struct
type Pusher struct {
	upgrader     websocket.Upgrader
	ErrorHandler ErrorHandler
	clients      *ClientManager
	channels     *ChannelManager
}

// NewPusher server
func NewPusher(config Configuration) *Pusher {
	return &Pusher{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ErrorHandler: func(w http.ResponseWriter, code int, err error) {
			xlog.Error(err)

			http.Error(w, err.Error(), code)
		},
		clients:  NewClientManager(),
		channels: NewChannelManager(),
	}
}

// JoinChannel added client to channel
func (p *Pusher) JoinChannel(ID string, client *Client) {
	// get or create channel
	channel, ok := p.channels.Get(ID)
	if ok == false {
		channel = NewChannel(ID)

		p.channels.Add(channel)
	}

	channel.Join(client)
}

// LeaveChannel removed client to channel
func (p *Pusher) LeaveChannel(ID string, client *Client) {
	if channel, ok := p.channels.Get(ID); ok {
		channel.Leave(client)
	}
}

// ServeHTTP handler
func (p *Pusher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.ErrorHandler(w, http.StatusInternalServerError, err)

		return
	}

	client := NewClient(ws, p)
	p.clients.Add(client)

	client.Listen()

}
