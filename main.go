package main

import (
	"embed"

	"github.com/prodemmi/KodoBoard/core"
	"go.uber.org/zap"
)

//go:embed frontend/dist/*
//go:embed frontend/dist/assets/*
var staticFiles embed.FS

// TODO: this is a test todo
// DONE 2025-08-18 01:00 by prodemmi
func main() {
	var logger *zap.Logger
	if true {
		logger = core.NewSilenceLogger()
	} else {
		logger = core.NewLogger()
	}

	scanner := core.NewScanner(logger)

	server := core.NewServer(logger, staticFiles, scanner)

	server.StartServer()
}
