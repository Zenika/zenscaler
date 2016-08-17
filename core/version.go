package core

import "runtime"

var (
	// Version number
	Version = "unspecified-dev"
	// GoVersion at build time
	GoVersion = runtime.Version()
)
