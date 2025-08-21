package main

import (
	"embed"

	"github.com/prodemmi/kodo/core"
	"go.uber.org/zap"
)

//go:embed web/dist/*
//go:embed web/dist/assets/*
var staticFiles embed.FS

// TODO: this is a test todo
// This function initializes the application server. It sets up the logger
// depending on the environment, creates a new scanner for code analysis,
// and starts the HTTP server with static files and scanner integration.
// DONE 2025-08-21 19:14 by prodemmi
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
