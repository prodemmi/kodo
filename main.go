package main

import (
	"embed"
	"os"

	"github.com/prodemmi/kodo/core"
	"github.com/prodemmi/kodo/core/cli"
	"github.com/prodemmi/kodo/core/entities"
	"github.com/prodemmi/kodo/core/handlers"
	"github.com/prodemmi/kodo/core/services"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

//go:embed web/dist/*
//go:embed web/dist/assets/*
var staticFiles embed.FS

func main() {
	config := entities.NewDefaultConfig()

	pflag.StringVarP(&config.Flags.Config, "config", "c", config.Flags.Config, "Path to config file")
	pflag.BoolVarP(&config.Flags.Investor, "investor", "i", config.Flags.Investor, "Run in investor mode")
	pflag.BoolVarP(&config.Flags.Silent, "silent", "s", config.Flags.Silent, "Silent the logger")
	showHelp := pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()

	if *showHelp {
		cli.PrintHelp()
		return
	}

	var logger *zap.Logger
	if config.Flags.Silent {
		logger = services.NewSilenceLogger()
	} else {
		logger = services.NewLogger()
	}
	defer logger.Sync()

	// Initialize services
	settingsService := services.NewSettingsService(config, logger)

	// Prepare settings
	if err := settingsService.Initialize(); err != nil {
		logger.Fatal("failed to initialize settings", zap.Error(err))
		os.Exit(1)
	}

	noteService := services.NewNoteService(config, logger)
	historyService := services.NewHistoryService(config, logger)
	scannerService := services.NewScannerService(config, settingsService, historyService, logger)
	remoteService := services.NewRemoteManager(logger, settingsService, noteService)

	// Initialize handlers
	noteHandler := handlers.NewNoteHandler(logger, noteService, remoteService)
	historyHandler := handlers.NewHistoryHandler(logger, scannerService, historyService, settingsService)
	chatHandler := handlers.NewChatHandler(logger)
	settingsHandler := handlers.NewSettingHandler(logger, settingsService, scannerService)
	itemHandler := handlers.NewItemHandler(logger, scannerService, historyService, settingsService)

	// Prepare history service
	if err := historyService.Initialize(); err != nil {
		logger.Fatal("failed to initialize history service", zap.Error(err))
		os.Exit(1)
	}

	// Initialize and start the server
	server := core.NewServer(
		config,
		logger,
		noteHandler,
		historyHandler,
		chatHandler,
		settingsHandler,
		itemHandler,
		staticFiles,
		scannerService,
	)

	server.Start()
}
