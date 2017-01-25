// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/json"
)

// Decode json string to event
func Decode(message []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(message, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
