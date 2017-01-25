// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pusher

import (
	"sync"
)

// ChannelManager struct
type ChannelManager struct {
	channels    map[string]*Channel
	channelsMtx *sync.Mutex
}

// NewChannelManager constructor
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels:    make(map[string]*Channel),
		channelsMtx: &sync.Mutex{},
	}
}

// Add channel to manager
func (c *ChannelManager) Add(channel *Channel) {
	c.channelsMtx.Lock()
	defer c.channelsMtx.Unlock()

	c.channels[channel.ID] = channel
}

// Get channel by id
func (c *ChannelManager) Get(ID string) (*Channel, bool) {
	c.channelsMtx.Lock()
	defer c.channelsMtx.Unlock()

	channel, ok := c.channels[ID]

	return channel, ok
}

// Remove channel to manager
func (c *ChannelManager) Remove(channel *Channel) {
	c.channelsMtx.Lock()
	defer c.channelsMtx.Unlock()

	delete(c.channels, channel.ID)
}
