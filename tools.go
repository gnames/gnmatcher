//go:build tools
// +build tools

package main

import (
	// cobra is used for creating a scaffold of CLI applications
	_ "github.com/spf13/cobra"
	_ "github.com/spf13/cobra-cli"

	// benchstat runs tests multiple times and provides a summary of performance.
	_ "golang.org/x/perf/cmd/benchstat"
)
