// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pusher

import (
	"sync"
)

// ClientManager struct
type ClientManager struct {
	clients    map[string]*Client
	clientsMtx *sync.Mutex
}

// NewClientManager constructor
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[string]*Client),
		clientsMtx: &sync.Mutex{},
	}
}

// Add client to manager
func (c *ClientManager) Add(client *Client) {
	c.clientsMtx.Lock()
	defer c.clientsMtx.Unlock()

	c.clients[client.ID] = client
}

// Remove client to manager
func (c *ClientManager) Remove(client *Client) {
	c.clientsMtx.Lock()
	defer c.clientsMtx.Unlock()

	delete(c.clients, client.ID)
}

// Clients list
func (c *ClientManager) Clients() map[string]*Client {
	return c.clients
}
