package main

import (
	"embed"

	"github.com/prodemmi/kodo/core"
	"go.uber.org/zap"
)

//go:embed web/dist/*
//go:embed web/dist/assets/*
var staticFiles embed.FS

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
