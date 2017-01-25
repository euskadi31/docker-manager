// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pusher

import (
	"github.com/euskadi31/docker-manager/pusher/message"
)

// Channel struct
type Channel struct {
	ID       string
	clients  *ClientManager
	eventsCh chan *message.Event
}

// NewChannel constructor
func NewChannel(ID string) *Channel {
	return &Channel{
		ID:       ID,
		clients:  NewClientManager(),
		eventsCh: make(chan *message.Event),
	}
}

// Join client to channel
func (c *Channel) Join(client *Client) {
	c.clients.Add(client)

	client.Write(&message.Event{
		Type:    "subscribed",
		Channel: c.ID,
	})
}

// Leave client to channel
func (c *Channel) Leave(client *Client) {
	c.clients.Remove(client)

	client.Write(&message.Event{
		Type:    "unsubscribed",
		Channel: c.ID,
	})
}

// Write to channel
func (c *Channel) Write(event *message.Event) {
	c.eventsCh <- event
}

// Listen channel
func (c *Channel) Listen() {
	for {
		select {
		case e := <-c.eventsCh:
			for _, client := range c.clients.Clients() {
				client.Write(e)
			}
		}
	}
}
