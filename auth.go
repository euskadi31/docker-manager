// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/RangelReale/osin"
)

type OAuthStorage struct {
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *OAuthStorage) Clone() osin.Storage {
	return nil
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *OAuthStorage) Close() {

}

// GetClient loads the client by id (client_id)
func (s *OAuthStorage) GetClient(id string) (osin.Client, error) {
	return nil, nil
}

// SaveAuthorize saves authorize data.
func (s *OAuthStorage) SaveAuthorize(*osin.AuthorizeData) error {
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *OAuthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	return nil, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *OAuthStorage) RemoveAuthorize(code string) error {
	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *OAuthStorage) SaveAccess(*osin.AccessData) error {
	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *OAuthStorage) LoadAccess(token string) (*osin.AccessData, error) {
	return nil, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *OAuthStorage) RemoveAccess(token string) error {
	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *OAuthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	return nil, nil
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *OAuthStorage) RemoveRefresh(token string) error {
	return nil
}
