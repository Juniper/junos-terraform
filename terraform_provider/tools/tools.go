//go:build tools

package tools

import (
	// release
	_ "github.com/goreleaser/goreleaser"

	// publish
	_ "github.com/chrismarget/lambda-tf-registry/cmd/register"
)
