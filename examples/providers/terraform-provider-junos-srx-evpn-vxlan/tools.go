//go:build tools

package main

import (
	// release
	_ "github.com/goreleaser/goreleaser"

	// publish
	_ "github.com/chrismarget/lambda-tf-registry/cmd/register"
)
