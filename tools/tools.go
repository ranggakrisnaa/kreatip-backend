//go:build tools

// Package tools pins dev-only CLI tools so their versions are tracked in go.mod/go.sum.
// Run: go generate ./tools/... to install all tools.
package tools

import (
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
)

//go:generate go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
