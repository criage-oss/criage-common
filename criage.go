// Package criage-common provides shared types, configurations, and utilities
// for the Criage package manager ecosystem.
//
// This module contains common components used across:
// - criage-client: The main package manager CLI
// - criage-server: HTTP repository server
// - criage-mcp: MCP server for AI integration
//
// Key components:
// - types: Common data structures and types
// - config: Configuration management
// - archive: Archive creation and extraction utilities
package main

import (
	"fmt"
	"runtime"
)

const (
	// Version of the common module
	Version = "1.0.0"

	// ModuleName for identification
	ModuleName = "criage-common"
)

// GetVersion returns the current version of the module
func GetVersion() string {
	return Version
}

// GetInfo returns basic information about the module
func GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        ModuleName,
		"version":     Version,
		"go_version":  runtime.Version(),
		"platform":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"description": "Common types and utilities for Criage package manager",
	}
}

// This file serves as the main entry point for the module
// though it's primarily used as a library by other components
func main() {
	fmt.Printf("Criage Common Module v%s\n", Version)
	fmt.Println("This is a library module for the Criage package manager ecosystem.")
	fmt.Println("See documentation for usage examples.")
}
