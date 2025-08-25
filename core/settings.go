package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Settings represents the application settings structure
type Settings struct {
	// Kanban Settings
	KanbanColumns    []KanbanColumn   `json:"kanban_columns"`
	PriorityPatterns PriorityPatterns `json:"priority_patterns"`

	// Code Scan Settings
	CodeScanSettings CodeScanConfig `json:"code_scan_settings"`
	GithubAuth       GithubAuth     `json:"github_auth"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KanbanColumn represents a kanban board column
type KanbanColumn struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Color             string `json:"color"`
	AutoAssignPattern string `json:"auto_assign_pattern,omitempty"`
}

// PriorityPatterns represents priority patterns for task detection
type PriorityPatterns struct {
	Low    string `json:"low"`
	Medium string `json:"medium"`
	High   string `json:"high"`
}

// CodeScanConfig represents code scanning settings
type CodeScanConfig struct {
	ExcludeDirectories []string `json:"exclude_directories"`
	ExcludeFiles       []string `json:"exclude_files"`
	SyncEnabled        bool     `json:"sync_enabled"`
}

// GithubAuth represents GitHub authentication settings
type GithubAuth struct {
	Token string `json:"token"`
}

// SettingsManager handles settings persistence and management
type SettingsManager struct {
	config       Config
	logger       *zap.Logger
	projectDir   string
	settingsFile string
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager(config Config, logger *zap.Logger) *SettingsManager {
	wd, _ := os.Getwd()
	projectDir := wd

	return &SettingsManager{
		config:       config,
		logger:       logger,
		projectDir:   projectDir,
		settingsFile: filepath.Join(projectDir, config.Flags.Config, "settings.json"),
	}
}

// GetDefaultSettings returns the default application settings
func (sm *SettingsManager) GetDefaultSettings() *Settings {
	return &Settings{
		KanbanColumns: []KanbanColumn{
			{ID: "todo", Name: "Todo", Color: "blue", AutoAssignPattern: "TODO"},
			{ID: "progress", Name: "In Progress", Color: "orange"},
			{ID: "done", Name: "Done", Color: "dark"},
		},
		PriorityPatterns: PriorityPatterns{
			Low:    "low:",
			Medium: "medium:",
			High:   "high:",
		},
		CodeScanSettings: CodeScanConfig{
			ExcludeDirectories: []string{
				"node_modules",
				".git",
				"dist",
				"build",
				"vendor",
				".next",
				"target",
			},
			ExcludeFiles: []string{
				"*.min.js",
				"*.min.css",
				"*.map",
				"package-lock.json",
				"yarn.lock",
			},
			SyncEnabled: false,
		},
		GithubAuth: GithubAuth{
			Token: "",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// LoadSettings loads settings from file or returns defaults
func (sm *SettingsManager) LoadSettings() *Settings {
	data, err := os.ReadFile(sm.settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			sm.logger.Info("Settings file not found, using defaults")
			return sm.GetDefaultSettings()
		}
		sm.logger.Error("Failed to read settings file", zap.Error(err))
		return sm.GetDefaultSettings()
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		sm.logger.Error("Failed to unmarshal settings", zap.Error(err))
		return sm.GetDefaultSettings()
	}

	// Validate and fix settings if needed
	settings = *sm.validateSettings(&settings)

	return &settings
}

// SaveSettings saves settings to file
func (sm *SettingsManager) SaveSettings(settings *Settings) error {
	// Validate settings before saving
	settings = sm.validateSettings(settings)
	settings.UpdatedAt = time.Now()

	// Ensure CreatedAt is set
	if settings.CreatedAt.IsZero() {
		settings.CreatedAt = time.Now()
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %v", err)
	}

	if err := os.WriteFile(sm.settingsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %v", err)
	}

	sm.logger.Info("Settings saved successfully",
		zap.String("file", sm.settingsFile),
		zap.Int("kanban_columns", len(settings.KanbanColumns)),
		zap.Bool("sync_enabled", settings.CodeScanSettings.SyncEnabled))

	return nil
}

// validateSettings ensures settings are valid and complete
func (sm *SettingsManager) validateSettings(settings *Settings) *Settings {
	// Ensure we have at least the basic kanban columns
	if len(settings.KanbanColumns) == 0 {
		settings.KanbanColumns = sm.GetDefaultSettings().KanbanColumns
	}

	// Ensure we have default exclude directories if none set
	if len(settings.CodeScanSettings.ExcludeDirectories) == 0 {
		settings.CodeScanSettings.ExcludeDirectories = sm.GetDefaultSettings().CodeScanSettings.ExcludeDirectories
	}

	// Ensure we have default exclude files if none set
	if len(settings.CodeScanSettings.ExcludeFiles) == 0 {
		settings.CodeScanSettings.ExcludeFiles = sm.GetDefaultSettings().CodeScanSettings.ExcludeFiles
	}

	return settings
}

// UpdatePartialSettings updates only the provided fields in settings
func (sm *SettingsManager) UpdatePartialSettings(updates map[string]interface{}) (*Settings, error) {
	// Load current settings
	settings := sm.LoadSettings()

	// Apply updates
	if kanbanColumns, ok := updates["kanban_columns"]; ok {
		if columnsData, ok := kanbanColumns.([]interface{}); ok {
			var columns []KanbanColumn
			for _, col := range columnsData {
				if colMap, ok := col.(map[string]interface{}); ok {
					column := KanbanColumn{}
					if id, ok := colMap["id"].(string); ok {
						column.ID = id
					}
					if name, ok := colMap["name"].(string); ok {
						column.Name = name
					}
					if color, ok := colMap["color"].(string); ok {
						column.Color = color
					}
					if pattern, ok := colMap["auto_assign_pattern"].(string); ok {
						column.AutoAssignPattern = pattern
					}
					columns = append(columns, column)
				}
			}
			settings.KanbanColumns = columns
		}
	}

	if priority_patterns, ok := updates["priority_patterns"]; ok {
		if prioritiesMap, ok := priority_patterns.(map[string]interface{}); ok {
			if low, ok := prioritiesMap["low"].(string); ok {
				settings.PriorityPatterns.Low = low
			}
			if medium, ok := prioritiesMap["medium"].(string); ok {
				settings.PriorityPatterns.Medium = medium
			}
			if high, ok := prioritiesMap["high"].(string); ok {
				settings.PriorityPatterns.High = high
			}
		}
	}

	if code_scan_settings, ok := updates["code_scan_settings"]; ok {
		if cssMap, ok := code_scan_settings.(map[string]interface{}); ok {
			if excludeDirs, ok := cssMap["exclude_directories"].([]interface{}); ok {
				var dirs []string
				for _, dir := range excludeDirs {
					if dirStr, ok := dir.(string); ok {
						dirs = append(dirs, dirStr)
					}
				}
				settings.CodeScanSettings.ExcludeDirectories = dirs
			}
			if exclude_files, ok := cssMap["exclude_files"].([]interface{}); ok {
				var files []string
				for _, file := range exclude_files {
					if fileStr, ok := file.(string); ok {
						files = append(files, fileStr)
					}
				}
				settings.CodeScanSettings.ExcludeFiles = files
			}
			if sync_enabled, ok := cssMap["sync_enabled"].(bool); ok {
				settings.CodeScanSettings.SyncEnabled = sync_enabled
			}
		}
	}

	if github_auth, ok := updates["github_auth"]; ok {
		if gaMap, ok := github_auth.(map[string]interface{}); ok {
			if token, ok := gaMap["token"].(string); ok {
				settings.GithubAuth.Token = token
			}
		}
	}

	// Save updated settings
	if err := sm.SaveSettings(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// GetSettingsSummary returns a summary of current settings
func (sm *SettingsManager) GetSettingsSummary() map[string]interface{} {
	settings := sm.LoadSettings()

	return map[string]interface{}{
		"kanban_columns_count": len(settings.KanbanColumns),
		"sync_enabled":         settings.CodeScanSettings.SyncEnabled,
		"exclude_directories":  len(settings.CodeScanSettings.ExcludeDirectories),
		"exclude_files":        len(settings.CodeScanSettings.ExcludeFiles),
		"has_github_token":     settings.GithubAuth.Token != "",
		"created_at":           settings.CreatedAt,
		"updated_at":           settings.UpdatedAt,
	}
}

// Helper method to check if slice contains string
func (sm *SettingsManager) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
