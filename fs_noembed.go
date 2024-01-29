//go:build !embed

package main

import "net/http"

var staticFS http.FileSystem
