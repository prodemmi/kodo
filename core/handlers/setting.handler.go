package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prodemmi/kodo/core/services"
	"go.uber.org/zap"
)

type SettingHandler struct {
	logger          *zap.Logger
	settingsService *services.SettingsService
	scannerService  *services.ScannerService
}

func NewSettingHandler(logger *zap.Logger, settingsService *services.SettingsService,
	scannerService *services.ScannerService) *SettingHandler {
	return &SettingHandler{
		logger:          logger,
		settingsService: settingsService,
		scannerService:  scannerService,
	}
}

func (s *SettingHandler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		settings := s.settingsService.LoadSettings()
		json.NewEncoder(w).Encode(settings)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *SettingHandler) HandleSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		s.logger.Error("Invalid JSON for settings update", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.logger.Info("Updating settings", zap.Any("updates", updateReq))
	oldSettings := s.settingsService.LoadSettings()

	updatedSettings, err := s.settingsService.UpdatePartialSettings(updateReq)
	if err != nil {
		s.logger.Error("Failed to update settings", zap.Error(err))
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	err = s.scannerService.UpdateOldStatuses(oldSettings, updatedSettings)
	if err != nil {
		s.logger.Error("Failed to update settings", zap.Error(err))
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Settings updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"settings": updatedSettings,
		"message":  "Settings updated successfully",
	})
}
