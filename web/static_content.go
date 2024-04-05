package web

import "embed"

//go:embed dist/*
var StaticContent embed.FS
