package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prodemmi/kodo/core/services"
	"go.uber.org/zap"
)

type HistoryHandler struct {
	logger          *zap.Logger
	scannerService  *services.ScannerService
	historyService  *services.HistoryService
	settingsService *services.SettingsService
}

func NewHistoryHandler(logger *zap.Logger,
	scannerService *services.ScannerService,
	historyService *services.HistoryService,
	settingsService *services.SettingsService) *HistoryHandler {
	return &HistoryHandler{
		logger:          logger,
		scannerService:  scannerService,
		historyService:  historyService,
		settingsService: settingsService,
	}
}

func (s *HistoryHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		history := s.historyService.GetProjectStats(s.settingsService)
		json.NewEncoder(w).Encode(history)
	} else if r.Method == "POST" {
		s.scannerService.Rescan()
		history := s.historyService.GetProjectStats(s.settingsService)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "refreshed",
			"history": history,
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *HistoryHandler) HandleStatsHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	history := s.historyService.GetBranchHistory()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
		"count":   len(history),
	})
}

func (s *HistoryHandler) HandleStatsCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	comparison := s.historyService.CompareWithPrevious(s.settingsService)
	json.NewEncoder(w).Encode(comparison)
}

func (s *HistoryHandler) HandleStatsCleanup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	historyService := s.historyService
	if err := historyService.CleanupOldStats(); err != nil {
		s.logger.Error("Failed to cleanup old history", zap.Error(err))
		http.Error(w, "Failed to cleanup history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Old history cleaned up",
	})
}

func (s *HistoryHandler) HandleStatsItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	historyService := s.historyService
	analysis := historyService.GetTaskItemsAnalysis(s.settingsService)
	json.NewEncoder(w).Encode(analysis)
}

func (s *HistoryHandler) HandleStatsItemsByFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	historyService := s.historyService
	fileGroups := historyService.GetItemsByFile(s.settingsService)
	json.NewEncoder(w).Encode(fileGroups)
}

func (s *HistoryHandler) HandleStatsTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	historyService := s.historyService
	trends := historyService.GetItemTrends(s.settingsService)
	json.NewEncoder(w).Encode(trends)
}

func (s *HistoryHandler) HandleStatsChanges(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	historyService := s.historyService
	changes := historyService.GetRecentItemChanges()
	json.NewEncoder(w).Encode(changes)
}
