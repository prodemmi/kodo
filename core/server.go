package core

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Server struct {
	config      Config
	scanner     *Scanner
	staticFiles embed.FS

	logger *zap.Logger
}

func NewServer(config Config, logger *zap.Logger, staticFiles embed.FS, scanner *Scanner) *Server {
	return &Server{
		config:      config,
		staticFiles: staticFiles,
		scanner:     scanner,
		logger:      logger,
	}
}

func (s *Server) StartServer() {
	cors := func(next http.Handler) http.Handler {
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

	fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		filePath := "web/dist/" + path
		s.logger.Debug("trying to read file", zap.String("filename", filePath))
		data, err := s.staticFiles.ReadFile(filePath)
		if err != nil {
			s.logger.Debug("file not found", zap.String("filename", filePath), zap.Error(err))
			http.NotFound(w, r)
			return
		}

		switch ext := filepath.Ext(path); ext {
		case ".html", ".htm":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".woff":
			w.Header().Set("Content-Type", "font/woff")
		case ".woff2":
			w.Header().Set("Content-Type", "font/woff2")
		case ".ttf":
			w.Header().Set("Content-Type", "font/ttf")
		case ".eot":
			w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
		case ".mp4":
			w.Header().Set("Content-Type", "video/mp4")
		case ".webm":
			w.Header().Set("Content-Type", "video/webm")
		case ".ogg":
			w.Header().Set("Content-Type", "audio/ogg")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		s.logger.Debug("serving file", zap.String("filename", filePath), zap.Int("size", len(data)))
		w.Write(data)
	})

	http.Handle("/", cors(fileHandler))

	http.Handle("/api/investor", cors(http.HandlerFunc(s.handleInvestor)))

	http.Handle("/api/items", cors(http.HandlerFunc(s.handleItems)))
	http.Handle("/api/items/update", cors(http.HandlerFunc(s.handleUpdateTodo)))
	http.Handle("/api/refresh", cors(http.HandlerFunc(s.handleRefresh)))
	http.Handle("/api/open-file", cors(http.HandlerFunc(s.handleOpenFile)))
	http.Handle("/api/get-context", cors(http.HandlerFunc(s.handleGetContext)))

	http.Handle("/api/stats", cors(http.HandlerFunc(s.handleStats)))
	http.Handle("/api/stats/history", cors(http.HandlerFunc(s.handleStatsHistory)))
	http.Handle("/api/stats/compare", cors(http.HandlerFunc(s.handleStatsCompare)))
	http.Handle("/api/stats/cleanup", cors(http.HandlerFunc(s.handleStatsCleanup)))
	http.Handle("/api/stats/items", cors(http.HandlerFunc(s.handleStatsItems)))
	http.Handle("/api/stats/items/by-file", cors(http.HandlerFunc(s.handleStatsItemsByFile)))
	http.Handle("/api/stats/trends", cors(http.HandlerFunc(s.handleStatsTrends)))
	http.Handle("/api/stats/changes", cors(http.HandlerFunc(s.handleStatsChanges)))

	http.Handle("/api/notes", cors(http.HandlerFunc(s.handleNotes)))
	http.Handle("/api/notes/search", cors(http.HandlerFunc(s.handleNoteSearch)))
	http.Handle("/api/notes/update", cors(http.HandlerFunc(s.handleNoteUpdate)))
	http.Handle("/api/notes/delete", cors(http.HandlerFunc(s.handleNoteDelete)))
	http.Handle("/api/notes/move", cors(http.HandlerFunc(s.handleMoveNotes)))
	http.Handle("/api/notes/export", cors(http.HandlerFunc(s.handleExportNotes)))
	http.Handle("/api/notes/tags", cors(http.HandlerFunc(s.handleNoteTags)))
	http.Handle("/api/notes/stats", cors(http.HandlerFunc(s.handleNoteStats)))

	http.Handle("/api/notes/history", cors(http.HandlerFunc(s.handleNoteHistory)))
	http.Handle("/api/notes/history/stats", cors(http.HandlerFunc(s.handleNoteHistoryStats)))
	http.Handle("/api/notes/history/cleanup", cors(http.HandlerFunc(s.handleCleanupHistory)))
	http.Handle("/api/notes/with-history", cors(http.HandlerFunc(s.handleNoteWithHistory)))

	// Activity tracking routes
	http.Handle("/api/notes/activity/author", cors(http.HandlerFunc(s.handleAuthorActivity)))
	http.Handle("/api/notes/activity/branch", cors(http.HandlerFunc(s.handleBranchActivity)))

	http.Handle("/api/folders", cors(http.HandlerFunc(s.handleFolders)))
	http.Handle("/api/folders/update", cors(http.HandlerFunc(s.handleFolderUpdate)))
	http.Handle("/api/folders/delete", cors(http.HandlerFunc(s.handleFolderDelete)))
	http.Handle("/api/folders/tree", cors(http.HandlerFunc(s.handleFolderTree)))
	http.Handle("/api/categories", cors(http.HandlerFunc(s.handleCategories)))

	http.Handle("/api/chat/project-files", cors(http.HandlerFunc(s.handleProjectFiles)))

	port := 8080
	url := fmt.Sprintf("http://%s:%d", s.getLocalIP(), port)
	s.showServerInfo(url)

	s.logger.Debug("found items in the project", zap.Int("length", s.scanner.GetItemsLength()))

	s.scanner.LoadExistingStats()

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *Server) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func (s *Server) showServerInfo(url string) {
	// Banner
	banner := figure.NewFigure("KODO", "o8", true)
	color.Cyan(banner.String())

	fmt.Println()

	fmt.Println(color.GreenString("======================================"))

	if s.config.Flags.Config != ".kodo" && s.config.Flags.Config != "./.kodo" {
		fmt.Println(color.YellowString("▶ Config Path: %s", s.config.Flags.Config))
	}

	if s.config.Flags.Investor {
		fmt.Println(color.BlueString("▶ Investor Mode: Enabled"))
	}

	fmt.Println(color.GreenString("▶ Running at: %s", url))

	fmt.Println(color.GreenString("======================================"))
}

func (s *Server) handleInvestor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"investor": s.config.Flags.Investor})
}

func (s *Server) handleItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.scanner.Rescan()
	json.NewEncoder(w).Encode(s.scanner.GetItems())
}

func (s *Server) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		s.logger.Error("Invalid JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.logger.Info("Updating item status", zap.Int("id", updateReq.ID), zap.String("new_status", updateReq.Status))

	// Find the item
	var targetItem *Item
	for _, item := range s.scanner.GetItems() {
		if item.ID == updateReq.ID {
			targetItem = item
			break
		}
	}

	if targetItem == nil {
		s.logger.Error("Item not found", zap.Int("id", updateReq.ID))
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	s.logger.Info("Found item", zap.String("file", targetItem.File), zap.Int("line", targetItem.Line), zap.String("current_status", string(targetItem.Status)))

	// Update todo status using scanner methods
	var err error
	switch updateReq.Status {
	case "done":
		err = s.scanner.MarkAsDone(targetItem)
	case "in-progress":
		err = s.scanner.MarkAsInProgress(targetItem)
	case "todo":
		err = s.scanner.MarkAsUndone(targetItem)
	default:
		err = s.scanner.MarkAsUndone(targetItem)
	}

	if err != nil {
		s.logger.Error("Failed to update status", zap.Int("id", targetItem.ID), zap.String("status", updateReq.Status), zap.Error(err))
		http.Error(w, fmt.Sprintf("Failed to update status: %v", err), http.StatusInternalServerError)
		return
	}

	// Save stats after successful update
	if err := s.scanner.GetTracker().SaveStats(); err != nil {
		s.logger.Warn("Failed to save stats after item update", zap.Error(err))
	}

	s.logger.Info("Successfully updated item status", zap.Int("id", targetItem.ID), zap.String("new_status", string(targetItem.Status)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"item":   targetItem,
	})
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	s.scanner.Rescan()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  s.scanner.GetItemsLength(),
	})
}

func (s *Server) handleOpenFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		File string `json:"file"`
		Line int    `json:"line"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get absolute path
	wd, _ := os.Getwd()
	fullPath := filepath.Join(wd, req.File)

	// Try to open in various IDEs/editors
	opened := s.tryOpenInIDE(fullPath, req.Line)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"opened": opened,
	})
}

func (s *Server) handleGetContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		File string `json:"file"`
		Line int    `json:"line"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	context := s.getCodeContext(req.File, req.Line)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(context)
}

func (s *Server) getCodeContext(filePath string, itemLine int) map[string]interface{} {
	wd, _ := os.Getwd()
	fullPath := filepath.Join(wd, filePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return map[string]interface{}{
			"error": "Could not open file",
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		lines = append(lines, scanner.Text())
	}

	if itemLine > len(lines) || itemLine < 1 {
		return map[string]interface{}{
			"error": "Invalid line number",
		}
	}

	// Get context around the TODO (±10 lines)
	bound := 40
	start := itemLine - 5
	if start < 0 {
		start = 0
	}
	end := itemLine + bound
	if end > len(lines) {
		end = len(lines)
	}

	contextLines := make([]map[string]interface{}, 0)
	for i := start; i < end; i++ {
		contextLines = append(contextLines, map[string]interface{}{
			"number":  i + 1,
			"content": lines[i],
			"isTodo":  i+1 == itemLine,
		})
	}

	return map[string]interface{}{
		"file":     filePath,
		"itemLine": itemLine,
		"lines":    contextLines,
	}
}

// indentationLevel counts leading spaces (tabs count as 4 spaces).
func (s *Server) indentationLevel(line string) int {
	line = strings.ReplaceAll(line, "\t", "    ")
	return len(line) - len(strings.TrimLeft(line, " "))
}

// contains checks if a slice contains a string
func (s *Server) contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// cross-platform browser opener
func (s *Server) openBrowser(url string) error {
	var cmdName string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmdName = "open" // macOS: open default browser
		args = []string{url}
	case "windows":
		cmdName = "cmd" // Windows: start default browser
		args = []string{"/C", "start", url}
	default: // Linux
		cmdName = "xdg-open" // Linux: open default browser
		args = []string{url}
	}

	return exec.Command(cmdName, args...).Start()
}

func (s *Server) tryOpenInIDE(filePath string, line int) bool {
	// List of IDE commands to try (in order of preference)
	var commands [][]string

	switch runtime.GOOS {
	case "darwin": // macOS
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)}, // VS Code
			{"subl", fmt.Sprintf("%s:%d", filePath, line)},       // Sublime Text
			{"atom", fmt.Sprintf("%s:%d", filePath, line)},       // Atom
			{"vim", fmt.Sprintf("+%d", line), filePath},          // Vim
			{"nvim", fmt.Sprintf("+%d", line), filePath},         // Neovim
			{"open", filePath}, // Default app
		}
	case "windows":
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)}, // VS Code
			{"notepad++", fmt.Sprintf("-n%d", line), filePath},   // Notepad++
			{"notepad", filePath},                                // Notepad
		}
	default: // Linux/Unix
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)}, // VS Code
			{"subl", fmt.Sprintf("%s:%d", filePath, line)},       // Sublime Text
			{"atom", fmt.Sprintf("%s:%d", filePath, line)},       // Atom
			{"vim", fmt.Sprintf("+%d", line), filePath},          // Vim
			{"nvim", fmt.Sprintf("+%d", line), filePath},         // Neovim
			{"gedit", fmt.Sprintf("+%d", line), filePath},        // Gedit
			{"xdg-open", filePath},                               // Default app
		}
	}

	// Try each command until one works
	for _, cmd := range commands {
		if len(cmd) > 0 {
			err := exec.Command(cmd[0], cmd[1:]...).Start()
			if err == nil {
				s.logger.Info("file opened", zap.String("filename", filePath), zap.Int("line", line), zap.String("command", cmd[0]))
				return true
			}
		}
	}

	s.logger.Error("could not open %s:%d in any IDE", zap.String("filename", filePath), zap.Int("line", line))
	return false
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		// Return current stats
		stats := s.scanner.GetProjectStats()
		json.NewEncoder(w).Encode(stats)
	} else if r.Method == "POST" {
		// Force refresh stats
		s.scanner.Rescan()
		stats := s.scanner.GetProjectStats()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "refreshed",
			"stats":  stats,
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleStatsHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	history := s.scanner.GetBranchHistory()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
		"count":   len(history),
	})
}

func (s *Server) handleStatsCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	comparison := s.scanner.CompareWithPrevious()
	json.NewEncoder(w).Encode(comparison)
}

func (s *Server) handleStatsCleanup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tracker := s.scanner.GetTracker()
	if err := tracker.CleanupOldStats(); err != nil {
		s.logger.Error("Failed to cleanup old stats", zap.Error(err))
		http.Error(w, "Failed to cleanup stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Old stats cleaned up",
	})
}

func (s *Server) handleProjectInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get current items stats
	items := s.scanner.GetItems()
	categories := s.scanner.GetItemsByCategory()

	// Get project stats
	projectStats := s.scanner.GetProjectStats()

	// Get recent history
	history := s.scanner.GetBranchHistory()
	recentHistory := history
	if len(history) > 10 {
		recentHistory = history[len(history)-10:]
	}

	response := map[string]interface{}{
		"current_items":  len(items),
		"categories":     categories,
		"project_stats":  projectStats,
		"recent_history": recentHistory,
		"git_info": map[string]string{
			"branch": s.scanner.GetTracker().getGitBranch(),
			"commit": s.scanner.GetTracker().getGitCommitShort(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleStatsItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tracker := s.scanner.GetTracker()
	analysis := tracker.GetTaskItemsAnalysis()
	json.NewEncoder(w).Encode(analysis)
}

func (s *Server) handleStatsItemsByFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tracker := s.scanner.GetTracker()
	fileGroups := tracker.GetItemsByFile()
	json.NewEncoder(w).Encode(fileGroups)
}

// Handler for item trends over time
func (s *Server) handleStatsTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tracker := s.scanner.GetTracker()
	trends := tracker.GetItemTrends()
	json.NewEncoder(w).Encode(trends)
}

// Handler for recent item changes
func (s *Server) handleStatsChanges(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tracker := s.scanner.GetTracker()
	changes := tracker.getRecentItemChanges()
	json.NewEncoder(w).Encode(changes)
}

func (s *Server) handleProjectFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			http.Error(w, "Cannot get working directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	search := r.URL.Query().Get("search")

	var files []ProjectFile
	var err error
	if search != "" {
		files, err = searchFilesRecursive(dir, search)
	} else {
		files, err = scanDirOneLevel(dir)
	}
	if err != nil {
		http.Error(w, "Cannot read directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Get all notes or filter by query parameters
		category := r.URL.Query().Get("category")
		tag := r.URL.Query().Get("tag")
		folderId := r.URL.Query().Get("folderId")

		notes, err := s.getNotes(category, tag, folderId)
		if err != nil {
			s.logger.Error("Failed to get notes", zap.Error(err))
			http.Error(w, "Failed to get notes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"notes": notes,
			"count": len(notes),
		})

	case "POST":
		// Create new note
		var noteReq struct {
			Title    string   `json:"title"`
			Content  string   `json:"content"`
			Tags     []string `json:"tags"`
			Category string   `json:"category"`
			FolderID *int     `json:"folderId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
			s.logger.Error("Invalid JSON for note creation", zap.Error(err))
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if noteReq.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		// Use history-enabled method
		note, err := s.createNoteWithHistory(noteReq.Title, noteReq.Content, noteReq.Tags, noteReq.Category, noteReq.FolderID)
		if err != nil {
			s.logger.Error("Failed to create note", zap.Error(err))
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		s.logger.Info("Note created with history tracking", zap.Int("id", note.ID), zap.String("title", note.Title))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"note":   note,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleNoteUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq struct {
		ID       int      `json:"id"`
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Tags     []string `json:"tags"`
		Category string   `json:"category"`
		Pinned   bool     `json:"pinned"`
		FolderID *int     `json:"folderId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		s.logger.Error("Invalid JSON for note update", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if updateReq.ID <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	if updateReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Use history-enabled method
	note, err := s.updateNoteWithHistory(updateReq.ID, updateReq.Title, updateReq.Content, updateReq.Tags, updateReq.Category, updateReq.Pinned, updateReq.FolderID)
	if err != nil {
		s.logger.Error("Failed to update note", zap.Int("id", updateReq.ID), zap.Error(err))
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Note updated with history tracking", zap.Int("id", note.ID), zap.String("title", note.Title))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"note":   note,
	})
}

func (s *Server) handleNoteDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteReq struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		s.logger.Error("Invalid JSON for note deletion", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if deleteReq.ID <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Use history-enabled method
	err := s.deleteNoteWithHistory(deleteReq.ID)
	if err != nil {
		s.logger.Error("Failed to delete note", zap.Int("id", deleteReq.ID), zap.Error(err))
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Note deleted with history tracking", zap.Int("id", deleteReq.ID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Note deleted successfully",
	})
}

// Add these handler methods for folders
func (s *Server) handleFolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		folders, err := s.getFolders()
		if err != nil {
			s.logger.Error("Failed to get folders", zap.Error(err))
			http.Error(w, "Failed to get folders", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"folders": folders,
			"count":   len(folders),
		})

	case "POST":
		var folderReq struct {
			Name     string `json:"name"`
			ParentID *int   `json:"parentId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&folderReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if folderReq.Name == "" {
			http.Error(w, "Folder name is required", http.StatusBadRequest)
			return
		}

		folder, err := s.createFolder(folderReq.Name, folderReq.ParentID)
		if err != nil {
			s.logger.Error("Failed to create folder", zap.Error(err))
			http.Error(w, "Failed to create folder", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"folder": folder,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories := s.getCategories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
		"count":      len(categories),
	})
}

func (s *Server) handleNoteStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.getNoteStats()
	if err != nil {
		s.logger.Error("Failed to get note stats", zap.Error(err))
		http.Error(w, "Failed to get note stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleFolderUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		ParentID *int   `json:"parentId"`
		Expanded *bool  `json:"expanded"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if updateReq.ID <= 0 {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := s.updateFolder(updateReq.ID, updateReq.Name, updateReq.ParentID, updateReq.Expanded)
	if err != nil {
		s.logger.Error("Failed to update folder", zap.Int("id", updateReq.ID), zap.Error(err))
		if err.Error() == "folder not found" {
			http.Error(w, "Folder not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update folder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"folder": folder,
	})
}

func (s *Server) handleFolderDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteReq struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if deleteReq.ID <= 0 {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	err := s.deleteFolder(deleteReq.ID)
	if err != nil {
		s.logger.Error("Failed to delete folder", zap.Int("id", deleteReq.ID), zap.Error(err))
		if err.Error() == "folder not found" {
			http.Error(w, "Folder not found", http.StatusNotFound)
			return
		}
		if err.Error() == "folder has notes" {
			http.Error(w, "Cannot delete folder that contains notes", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to delete folder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Folder deleted successfully",
	})
}

func (s *Server) handleFolderTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tree, err := s.getFolderTree()
	if err != nil {
		s.logger.Error("Failed to get folder tree", zap.Error(err))
		http.Error(w, "Failed to get folder tree", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tree":  tree,
		"count": len(tree),
	})
}

func (s *Server) handleNoteSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	folderIdStr := r.URL.Query().Get("folderId")

	var folderId *int
	if folderIdStr != "" {
		if id, err := strconv.Atoi(folderIdStr); err == nil {
			folderId = &id
		}
	}

	notes, err := s.searchNotes(query, category, folderId)
	if err != nil {
		s.logger.Error("Failed to search notes", zap.Error(err))
		http.Error(w, "Failed to search notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notes": notes,
		"count": len(notes),
		"query": query,
	})
}

func (s *Server) handleMoveNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var moveReq struct {
		NoteIds        []int `json:"noteIds"`
		TargetFolderId *int  `json:"targetFolderId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(moveReq.NoteIds) == 0 {
		http.Error(w, "No notes specified", http.StatusBadRequest)
		return
	}

	// Use history-enabled method
	err := s.moveNotesToFolderWithHistory(moveReq.NoteIds, moveReq.TargetFolderId)
	if err != nil {
		s.logger.Error("Failed to move notes", zap.Error(err))
		if err.Error() == "target folder not found" {
			http.Error(w, "Target folder not found", http.StatusNotFound)
			return
		}
		if err.Error() == "no notes found to move" {
			http.Error(w, "No notes found to move", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to move notes", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Notes moved with history tracking", zap.Ints("noteIds", moveReq.NoteIds), zap.Any("targetFolderId", moveReq.TargetFolderId))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Moved %d notes", len(moveReq.NoteIds)),
	})
}

func (s *Server) handleExportNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	folderIdStr := r.URL.Query().Get("folderId")

	var folderId *int
	if folderIdStr != "" {
		if id, err := strconv.Atoi(folderIdStr); err == nil {
			folderId = &id
		}
	}

	export, err := s.exportNotes(folderId, category)
	if err != nil {
		s.logger.Error("Failed to export notes", zap.Error(err))
		http.Error(w, "Failed to export notes", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("notes-export-%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	json.NewEncoder(w).Encode(export)
}

// handleNoteTags returns all unique tags from notes
func (s *Server) handleNoteTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tags, err := s.getAllTags()
	if err != nil {
		s.logger.Error("Failed to get tags", zap.Error(err))
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tags":  tags,
		"count": len(tags),
	})
}

// Add these handlers to your server.go file

// handleNoteHistory retrieves note history with filtering
func (s *Server) handleNoteHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for filtering
	var filter NoteHistoryFilter

	if noteIdStr := r.URL.Query().Get("noteId"); noteIdStr != "" {
		if noteId, err := strconv.Atoi(noteIdStr); err == nil {
			filter.NoteID = &noteId
		}
	}

	if actionStr := r.URL.Query().Get("action"); actionStr != "" {
		action := NoteHistoryAction(actionStr)
		filter.Action = &action
	}

	if author := r.URL.Query().Get("author"); author != "" {
		filter.Author = &author
	}

	if branch := r.URL.Query().Get("branch"); branch != "" {
		filter.GitBranch = &branch
	}

	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		if since, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			filter.Since = &since
		}
	}

	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		if until, err := time.Parse(time.RFC3339, untilStr); err == nil {
			filter.Until = &until
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 50 // Default limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	history, err := s.getNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get note history", zap.Error(err))
		http.Error(w, "Failed to get note history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
		"count":   len(history),
		"filter":  filter,
	})
}

// handleNoteHistoryStats retrieves note history statistics
func (s *Server) handleNoteHistoryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.getNoteHistoryStats()
	if err != nil {
		s.logger.Error("Failed to get note history stats", zap.Error(err))
		http.Error(w, "Failed to get note history stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleNoteWithHistory retrieves a note with its complete history
func (s *Server) handleNoteWithHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	noteIdStr := r.URL.Query().Get("id")
	if noteIdStr == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}

	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Get the note
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		s.logger.Error("Failed to load note storage", zap.Error(err))
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	var note *Note
	for _, n := range storage.Notes {
		if n.ID == noteId {
			note = &n
			break
		}
	}

	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Get note history
	filter := NoteHistoryFilter{NoteID: &noteId}
	history, err := s.getNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get note history", zap.Error(err))
		http.Error(w, "Failed to get note history", http.StatusInternalServerError)
		return
	}

	noteWithHistory := NoteWithHistory{
		Note:    *note,
		History: history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"note":          noteWithHistory,
		"history_count": len(history),
	})
}

// handleAuthorActivity retrieves activity summary for a specific author
func (s *Server) handleAuthorActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	author := r.URL.Query().Get("author")
	if author == "" {
		http.Error(w, "Author parameter is required", http.StatusBadRequest)
		return
	}

	// Get activity for specific author
	filter := NoteHistoryFilter{Author: &author, Limit: 100}
	history, err := s.getNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get author activity", zap.Error(err))
		http.Error(w, "Failed to get author activity", http.StatusInternalServerError)
		return
	}

	// Aggregate stats
	actionCount := make(map[NoteHistoryAction]int)
	noteCount := make(map[int]bool)
	dayCount := make(map[string]int)

	for _, entry := range history {
		actionCount[entry.Action]++
		noteCount[entry.NoteID] = true
		day := entry.Timestamp.Format("2006-01-02")
		dayCount[day]++
	}

	activity := map[string]interface{}{
		"author":         author,
		"total_actions":  len(history),
		"notes_affected": len(noteCount),
		"by_action":      actionCount,
		"by_day":         dayCount,
		"recent_history": history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

// handleBranchActivity retrieves activity summary for a specific git branch
func (s *Server) handleBranchActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	branch := r.URL.Query().Get("branch")
	if branch == "" {
		http.Error(w, "Branch parameter is required", http.StatusBadRequest)
		return
	}

	// Get activity for specific branch
	filter := NoteHistoryFilter{GitBranch: &branch, Limit: 100}
	history, err := s.getNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get branch activity", zap.Error(err))
		http.Error(w, "Failed to get branch activity", http.StatusInternalServerError)
		return
	}

	// Aggregate stats
	actionCount := make(map[NoteHistoryAction]int)
	authorCount := make(map[string]int)
	noteCount := make(map[int]bool)

	for _, entry := range history {
		actionCount[entry.Action]++
		authorCount[entry.Author]++
		noteCount[entry.NoteID] = true
	}

	activity := map[string]interface{}{
		"branch":         branch,
		"total_actions":  len(history),
		"notes_affected": len(noteCount),
		"by_action":      actionCount,
		"by_author":      authorCount,
		"recent_history": history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

// handleCleanupHistory removes old history entries
func (s *Server) handleCleanupHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OlderThanDays int `json:"older_than_days"`
		KeepMinimum   int `json:"keep_minimum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.OlderThanDays <= 0 {
		req.OlderThanDays = 90 // Default to 90 days
	}
	if req.KeepMinimum <= 0 {
		req.KeepMinimum = 10 // Keep at least 10 entries per note
	}

	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		s.logger.Error("Failed to load note storage", zap.Error(err))
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -req.OlderThanDays)

	// Group history by note ID
	noteHistory := make(map[int][]NoteHistoryEntry)
	for _, entry := range storage.History {
		noteHistory[entry.NoteID] = append(noteHistory[entry.NoteID], entry)
	}

	// Keep recent entries and minimum number per note
	var filteredHistory []NoteHistoryEntry
	removedCount := 0

	for _, entries := range noteHistory {
		// Sort entries by timestamp descending (newest first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})

		keptCount := 0
		for _, entry := range entries {
			// Always keep minimum number of entries per note
			// Also keep entries newer than cutoff
			if keptCount < req.KeepMinimum || entry.Timestamp.After(cutoff) {
				filteredHistory = append(filteredHistory, entry)
				keptCount++
			} else {
				removedCount++
			}
		}
	}

	storage.History = filteredHistory

	if err := s.saveEnhancedNoteStorage(storage); err != nil {
		s.logger.Error("Failed to save cleaned history", zap.Error(err))
		http.Error(w, "Failed to save changes", http.StatusInternalServerError)
		return
	}

	s.logger.Info("History cleanup completed", zap.Int("removed_entries", removedCount), zap.Int("remaining_entries", len(filteredHistory)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"removed_entries":   removedCount,
		"remaining_entries": len(filteredHistory),
		"message":           fmt.Sprintf("Removed %d old history entries", removedCount),
	})
}
