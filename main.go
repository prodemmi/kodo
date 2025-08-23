package main

import (
	"embed"
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/prodemmi/kodo/core"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

//go:embed web/dist/*
//go:embed web/dist/assets/*
var staticFiles embed.FS

func printHelp() {
	banner := figure.NewFigure("KODO ", "slant", true)
	color.Cyan(banner.String())
	fmt.Println()

	fmt.Println(color.YellowString("A source-code kanban and note app"))
	fmt.Println(color.GreenString("--------------------------------------------------"))
	fmt.Println(color.WhiteString("Usage:"))
	fmt.Println(color.WhiteString("  kodo [flags]"))
	fmt.Println()
	fmt.Println(color.WhiteString("Available Flags:"))
	fmt.Println(color.WhiteString("  -c, --config <path>     Path to config file (default .kodo)"))
	fmt.Println(color.WhiteString("  -i, --investor          Run in investor mode (default false)"))
	fmt.Println(color.WhiteString("  -h, --help              Show this help message"))
	fmt.Println(color.GreenString("--------------------------------------------------"))
	fmt.Println()
}

func main() {
	config := core.NewDefaultConfig()

	pflag.StringVarP(&config.Flags.Config, "config", "c", config.Flags.Config, "Path to config file")
	pflag.BoolVarP(&config.Flags.Investor, "investor", "i", config.Flags.Investor, "Run in investor mode")
	help := pflag.BoolP("help", "h", false, "Show help message")

	pflag.Parse()

	if help != nil && *help == true {
		printHelp()
		return
	}

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
