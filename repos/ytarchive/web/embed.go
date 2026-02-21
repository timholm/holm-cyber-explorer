// Package web provides embedded frontend assets
package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFS embed.FS

// GetFS returns the embedded filesystem rooted at dist/
func GetFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
