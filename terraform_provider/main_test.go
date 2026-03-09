package main

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// TestRunPassesDebugToServe verifies run forwards parsed debug state to Serve.
func TestRunPassesDebugToServe(t *testing.T) {
	originalParseFlags := parseFlags
	originalServeProvider := serveProvider
	t.Cleanup(func() {
		parseFlags = originalParseFlags
		serveProvider = originalServeProvider
	})

	parseFlags = func(debug *bool) {
		*debug = true
	}

	called := false
	serveProvider = func(_ context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		called = true
		if !opts.Debug {
			t.Fatalf("expected debug=true in serve options")
		}
		if opts.Address == "" {
			t.Fatalf("expected non-empty provider address")
		}
		if providerFunc() == nil {
			t.Fatalf("expected provider constructor to return non-nil provider")
		}
		return nil
	}

	if err := run(); err != nil {
		t.Fatalf("run() returned unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected serveProvider to be called")
	}
}

// TestMainLogsFatalOnRunError verifies main logs fatally when startup fails.
func TestMainLogsFatalOnRunError(t *testing.T) {
	originalParseFlags := parseFlags
	originalServeProvider := serveProvider
	originalFatalLogger := fatalLogger
	t.Cleanup(func() {
		parseFlags = originalParseFlags
		serveProvider = originalServeProvider
		fatalLogger = originalFatalLogger
	})

	parseFlags = func(_ *bool) {}
	serveErr := errors.New("serve failed")
	serveProvider = func(_ context.Context, _ func() provider.Provider, _ providerserver.ServeOpts) error {
		return serveErr
	}

	called := false
	fatalLogger = func(v ...interface{}) {
		called = true
		if len(v) != 1 || v[0] != serveErr {
			t.Fatalf("fatalLogger called with unexpected args: %#v", v)
		}
	}

	main()

	if !called {
		t.Fatalf("expected fatalLogger to be called")
	}
}
