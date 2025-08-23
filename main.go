package main

import (
	"embed"

	"github.com/prodemmi/kodo/core"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

//go:embed web/dist/*
//go:embed web/dist/assets/*
var staticFiles embed.FS

func main() {
	config := core.NewDefaultConfig()

	pflag.StringVarP(&config.Flags.Config, "config", "c", config.Flags.Config, "Path to config file")
	pflag.BoolVarP(&config.Flags.Investor, "investor", "i", config.Flags.Investor, "Run in investor mode")

	pflag.Parse()

	var logger *zap.Logger
	if true {
		logger = core.NewSilenceLogger()
	} else {
		logger = core.NewLogger()
	}

	scanner := core.NewScanner(config, logger)

	server := core.NewServer(config, logger, staticFiles, scanner)

	server.StartServer()
}
