package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/prodemmi/kodo/core/entities"
	"github.com/prodemmi/kodo/core/services"
	"go.uber.org/zap"
)

type ItemHandler struct {
	logger          *zap.Logger
	scannerService  *services.ScannerService
	historyService  *services.HistoryService
	settingsService *services.SettingsService
}

func NewItemHandler(logger *zap.Logger,
	scannerService *services.ScannerService,
	historyService *services.HistoryService,
	settingsService *services.SettingsService) *ItemHandler {
	return &ItemHandler{
		logger:          logger,
		scannerService:  scannerService,
		historyService:  historyService,
		settingsService: settingsService,
	}
}

func (s *ItemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := s.scannerService.Rescan()

	if err != nil {
		s.logger.Error("Failed to get items", zap.Error(err))
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(s.scannerService.GetItems())
}

func (s *ItemHandler) HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
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

	var targetItem *entities.Item
	for _, item := range s.scannerService.GetItems() {
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

	err := s.scannerService.UpdateItemStatus(targetItem, updateReq.Status)

	if err != nil {
		s.logger.Error("Failed to update status", zap.Int("id", targetItem.ID), zap.String("status", updateReq.Status), zap.Error(err))
		http.Error(w, fmt.Sprintf("Failed to update status: %v", err), http.StatusInternalServerError)
		return
	}

	if err := s.historyService.SaveStats(s.scannerService.GetItems(), s.settingsService); err != nil {
		s.logger.Warn("Failed to save history after item update", zap.Error(err))
	}

	s.logger.Info("Successfully updated item status", zap.Int("id", targetItem.ID), zap.String("new_status", string(targetItem.Status)))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"item":   targetItem,
	})
}

func (s *ItemHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	_ = s.scannerService.Rescan()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  s.scannerService.GetItemsLength(),
	})
}

func (s *ItemHandler) HandleOpenFile(w http.ResponseWriter, r *http.Request) {
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

	wd, _ := os.Getwd()
	fullPath := filepath.Join(wd, req.File)

	opened := s.tryOpenInIDE(fullPath, req.Line)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"opened": opened,
	})
}

func (s *ItemHandler) HandleGetContext(w http.ResponseWriter, r *http.Request) {
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
	_ = json.NewEncoder(w).Encode(context)
}

func (s *ItemHandler) getCodeContext(filePath string, itemLine int) map[string]interface{} {
	wd, _ := os.Getwd()
	fullPath := filepath.Join(wd, filePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return map[string]interface{}{
			"error": "Could not open file",
		}
	}
	defer func() {
		_ = file.Close()
	}()

	var lines []string
	scannerService := bufio.NewScanner(file)
	lineNum := 0

	for scannerService.Scan() {
		lineNum++
		lines = append(lines, scannerService.Text())
	}

	if itemLine > len(lines) || itemLine < 1 {
		return map[string]interface{}{
			"error": "Invalid line number",
		}
	}

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

func (s *ItemHandler) tryOpenInIDE(filePath string, line int) bool {

	var commands [][]string

	switch runtime.GOOS {
	case "darwin":
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)},
			{"subl", fmt.Sprintf("%s:%d", filePath, line)},
			{"atom", fmt.Sprintf("%s:%d", filePath, line)},
			{"vim", fmt.Sprintf("+%d", line), filePath},
			{"nvim", fmt.Sprintf("+%d", line), filePath},
			{"open", filePath},
		}
	case "windows":
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)},
			{"notepad++", fmt.Sprintf("-n%d", line), filePath},
			{"notepad", filePath},
		}
	default:
		commands = [][]string{
			{"code", "-g", fmt.Sprintf("%s:%d", filePath, line)},
			{"subl", fmt.Sprintf("%s:%d", filePath, line)},
			{"atom", fmt.Sprintf("%s:%d", filePath, line)},
			{"vim", fmt.Sprintf("+%d", line), filePath},
			{"nvim", fmt.Sprintf("+%d", line), filePath},
			{"gedit", fmt.Sprintf("+%d", line), filePath},
			{"xdg-open", filePath},
		}
	}

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
