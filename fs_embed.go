//go:build embed

package main

import (
	"embed"
	"net/http"
)

var (
	//go:embed static
	embedFS  embed.FS
	staticFS http.FileSystem
)

func init() {
	staticFS = http.FS(embedFS)
}
