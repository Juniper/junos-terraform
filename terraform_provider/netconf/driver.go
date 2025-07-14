// Copyright (c) 2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

// Driver interface for building drivers that are self-contained from a user's perspective.
type Driver interface {
	Lock(ds string) (*RPCReply, error)
	Unlock(ds string) (*RPCReply, error)

	Close() error
	Dial() error
	DialTimeout() error
	SendRaw(rawxml string) (*RPCReply, error)
	GetConfig() (*RPCReply, error)
}

// New is an interface that checks compliancy
func NewDriver(d Driver) Driver {
	return d
}
