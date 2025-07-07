# Go-NETCONF

[![GoDoc](https://godoc.org/github.com/Juniper/go-netconf/netconf?status.svg)](https://godoc.org/github.com/Juniper/go-netconf/netconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/Juniper/go-netconf)](https://goreportcard.com/report/github.com/Juniper/go-netconf)
[![Build Status](https://travis-ci.org/Juniper/go-netconf.svg?branch=master)](https://travis-ci.org/Juniper/go-netconf)

This library is a simple NETCONF client based on [RFC6241](http://tools.ietf.org/html/rfc6241) and [RFC6242](http://tools.ietf.org/html/rfc6242) (although not fully compliant yet).

> **Note:** This library is a fork of the Juniper go-netconf library. This fork is open to pull requests and suggestions. Libraries like this should be easy to use and this library has a nicer abstraction layer using a 'Driver' interface. This is work-in-progress.

## Features
* Support for SSH transport using go.crypto/ssh. (Other transports are planned).
* Built in RPC support (in progress).
* Support for custom RPCs.
* Independent of XML library.  Free to choose encoding/xml or another third party library to parse the results.

## Install
* Requires Go 1.4 or later
* Tests pass for 1.8 and later using Travis-CI (this is WIP)
* Please use `git clone` because this is a library and not an application

```bash
go get github.com/Juniper/go-netconf
```

## Example
* See examples in `examples/` directory.

This is the interface that `go-netconf` uses.

```go
// Driver interface for building drivers that are self-contained from a user's perspective.
type Driver interface {
	Lock(ds string) (*rpc.RPCReply, error)
	Unlock(ds string) (*rpc.RPCReply, error)

	Close() error
	Dial() error
	DialTimeout() error
	SendRaw(rawxml string) (*rpc.RPCReply, error)
	GetConfig() (*rpc.RPCReply, error)
}
```

Please note, that two transport types exist; one for direct use on Junos and the other socket based over SSH. This driver mechanism promotes the easy interchange of lower level transports without messing about.

## Documentation
You can view full API documentation at GoDoc: http://godoc.org/github.com/arsonistgopher/go-netconf
*The docs need updating. For now please review the source code*

## License
(BSD 2)

Copyright © 2013-2018 Juniper Networks, Inc. All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

(1) Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

(2) Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those of the authors and should not be interpreted as representing official policies, either expressed or implied, of Juniper Networks.

Authors and Contributors
------------------------
* [Brandon Bennett](https://github.com/nemith), Facebook
* [Charl Matthee](https://github.com/charl)
* [Jade Auer](https://github.com/jda)
* [Wayne Tucker](https://github.com/wtucker)
* [Christian Giese](https://github.com/GIC-de), Juniper Networks
* [David Gee](https://github.com/Juniper), Juniper Networks, IPEngineer.net
