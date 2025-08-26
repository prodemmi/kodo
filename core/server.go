package core

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/prodemmi/kodo/core/cli"
	"github.com/prodemmi/kodo/core/entities"
	"github.com/prodemmi/kodo/core/handlers"
	"github.com/prodemmi/kodo/core/services"
	"go.uber.org/zap"
)

type Server struct {
	config         *entities.Config
	scannerService *services.ScannerService
	staticFiles    embed.FS
	logger         *zap.Logger

	noteHandler     *handlers.NoteHandler
	historyHandler  *handlers.HistoryHandler
	chatHandler     *handlers.ChatHandler
	settingsHandler *handlers.SettingHandler
	itemHandler     *handlers.ItemHandler
}

func NewServer(
	config *entities.Config,
	logger *zap.Logger,
	noteHandler *handlers.NoteHandler,
	historyHandler *handlers.HistoryHandler,
	chatHandler *handlers.ChatHandler,
	settingsHandler *handlers.SettingHandler,
	itemHandler *handlers.ItemHandler,
	staticFiles embed.FS,
	scannerService *services.ScannerService,
) *Server {
	return &Server{
		config:          config,
		staticFiles:     staticFiles,
		scannerService:  scannerService,
		logger:          logger,
		noteHandler:     noteHandler,
		historyHandler:  historyHandler,
		chatHandler:     chatHandler,
		settingsHandler: settingsHandler,
		itemHandler:     itemHandler,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.Handle("/", s.withCORS(http.HandlerFunc(s.serveStatic)))

	s.registerItemRoutes(mux)
	s.registerNoteRoutes(mux)
	s.registerHistoryRoutes(mux)
	s.registerSettingsRoutes(mux)
	s.registerMiscRoutes(mux)

	port := 8080
	url := fmt.Sprintf("http://%s:%d", cli.GetLocalIP(), port)
	cli.ShowServerInfo(url, s.config)

	s.logger.Debug("found items in the project",
		zap.Int("length", s.scannerService.GetItemsLength()))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		s.logger.Fatal("server failed", zap.Error(err))
	}
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) serveStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	filePath := "web/dist/" + path
	data, err := s.staticFiles.ReadFile(filePath)
	if err != nil {
		s.logger.Debug("file not found", zap.String("filename", filePath), zap.Error(err))
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", s.detectContentType(path))
	_, _ = w.Write(data)
}

func (s *Server) detectContentType(path string) string {
	switch ext := filepath.Ext(path); ext {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".svg":
		return "image/svg+xml; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}

func (s *Server) handleInvestor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{
		"investor": s.config.Flags.Investor,
	})
}

func (s *Server) registerItemRoutes(mux *http.ServeMux) {
	mux.Handle("/api/items", s.withCORS(http.HandlerFunc(s.itemHandler.HandleItems)))
	mux.Handle("/api/items/update", s.withCORS(http.HandlerFunc(s.itemHandler.HandleUpdateTodo)))
	mux.Handle("/api/items/open-file", s.withCORS(http.HandlerFunc(s.itemHandler.HandleOpenFile)))
	mux.Handle("/api/items/get-context", s.withCORS(http.HandlerFunc(s.itemHandler.HandleGetContext)))
}

func (s *Server) registerNoteRoutes(mux *http.ServeMux) {
	mux.Handle("/api/notes", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNotes)))
	mux.Handle("/api/notes/search", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteSearch)))
	mux.Handle("/api/notes/update", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteUpdate)))
	mux.Handle("/api/notes/delete", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteDelete)))
	mux.Handle("/api/notes/move", s.withCORS(http.HandlerFunc(s.noteHandler.HandleMoveNotes)))
	mux.Handle("/api/notes/export", s.withCORS(http.HandlerFunc(s.noteHandler.HandleExportNotes)))
	mux.Handle("/api/notes/tags", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteTags)))
	mux.Handle("/api/notes/sync", s.withCORS(http.HandlerFunc(s.noteHandler.HandleSyncNotes)))

	mux.Handle("/api/notes/history", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteHistory)))
	mux.Handle("/api/notes/history/history", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteHistoryStats)))
	mux.Handle("/api/notes/history/cleanup", s.withCORS(http.HandlerFunc(s.noteHandler.HandleCleanupHistory)))
	mux.Handle("/api/notes/with-history", s.withCORS(http.HandlerFunc(s.noteHandler.HandleNoteWithHistory)))

	mux.Handle("/api/notes/activity/author", s.withCORS(http.HandlerFunc(s.noteHandler.HandleAuthorActivity)))
	mux.Handle("/api/notes/activity/branch", s.withCORS(http.HandlerFunc(s.noteHandler.HandleBranchActivity)))

	mux.Handle("/api/notes/folders", s.withCORS(http.HandlerFunc(s.noteHandler.HandleFolders)))
	mux.Handle("/api/notes/folders/update", s.withCORS(http.HandlerFunc(s.noteHandler.HandleFolderUpdate)))
	mux.Handle("/api/notes/folders/delete", s.withCORS(http.HandlerFunc(s.noteHandler.HandleFolderDelete)))
	mux.Handle("/api/notes/folders/tree", s.withCORS(http.HandlerFunc(s.noteHandler.HandleFolderTree)))
	mux.Handle("/api/notes/categories", s.withCORS(http.HandlerFunc(s.noteHandler.HandleCategories)))
}

func (s *Server) registerHistoryRoutes(mux *http.ServeMux) {
	mux.Handle("/api/history", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStats)))
	mux.Handle("/api/history/history", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsHistory)))
	mux.Handle("/api/history/compare", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsCompare)))
	mux.Handle("/api/history/cleanup", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsCleanup)))
	mux.Handle("/api/history/items", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsItems)))
	mux.Handle("/api/history/items/by-file", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsItemsByFile)))
	mux.Handle("/api/history/trends", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsTrends)))
	mux.Handle("/api/history/changes", s.withCORS(http.HandlerFunc(s.historyHandler.HandleStatsChanges)))
}

func (s *Server) registerSettingsRoutes(mux *http.ServeMux) {
	mux.Handle("/api/settings", s.withCORS(http.HandlerFunc(s.settingsHandler.HandleSettings)))
	mux.Handle("/api/settings/update", s.withCORS(http.HandlerFunc(s.settingsHandler.HandleSettingsUpdate)))
}

// func (s *Server) registerChatRoutes(mux *http.ServeMux) {
// 	mux.Handle("/api/chat/project-files", s.withCORS(http.HandlerFunc(s.chatHandler.HandleProjectFiles)))
// }

func (s *Server) registerMiscRoutes(mux *http.ServeMux) {
	mux.Handle("/api/investor", s.withCORS(http.HandlerFunc(s.handleInvestor)))
}
