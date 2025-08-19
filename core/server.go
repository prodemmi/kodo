package core

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

type Server struct {
	scanner     *Scanner
	staticFiles embed.FS

	logger *zap.Logger
}

func NewServer(logger *zap.Logger, staticFiles embed.FS, scanner *Scanner) *Server {
	return &Server{
		staticFiles: staticFiles,
		scanner:     scanner,
		logger:      logger,
	}
}

func (s *Server) StartServer() {
	// Middleware برای CORS
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

	// فایل‌های استاتیک
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

	port := 8080
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Println("Server running at", url)
	s.logger.Info("found items in the project", zap.Int("length", s.scanner.GetItemsLength()))

	s.scanner.LoadExistingStats()

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
	start := itemLine - (bound + 1)
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
