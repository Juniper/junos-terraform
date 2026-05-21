//go:build tools

package tools

import (
	// release
	_ "github.com/goreleaser/goreleaser/v2"

	// publish
	_ "github.com/chrismarget/lambda-tf-registry/cmd/register"
)
