package main

import (
	"embed"
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/prodemmi/kodo/core"
	"github.com/prodemmi/kodo/core/entities"
	"github.com/prodemmi/kodo/core/handlers"
	"github.com/prodemmi/kodo/core/services"
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
	config := entities.NewDefaultConfig()

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
		logger = services.NewSilenceLogger()
	} else {
		logger = services.NewLogger()
	}

	settingsService := services.NewSettingsService(config, logger)
	noteService := services.NewNoteService(config, logger)
	historyService := services.NewHistoryService(config, logger)
	scannerService := services.NewScannerService(config, settingsService, historyService, logger)
	remoteService := services.NewRemoteManager(logger, settingsService, noteService)

	noteHandler := handlers.NewNoteHandler(logger, noteService, remoteService)
	historyHandler := handlers.NewHistoryHandler(logger, scannerService, historyService, settingsService)
	chatHandler := handlers.NewChatHandler(logger)
	settingsHandler := handlers.NewSettingHandler(logger)
	itemHandler := handlers.NewItemHandler(logger, scannerService, historyService, settingsService)

	historyService.Initialize()

	server := core.NewServer(config,
		logger,
		noteHandler,
		historyHandler,
		chatHandler,
		settingsHandler,
		itemHandler,
		staticFiles,
		scannerService)

	server.Start()
}
