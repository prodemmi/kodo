package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prodemmi/kodo/core/entities"
	"go.uber.org/zap"
)

// SettingsService handles settings persistence and management
type SettingsService struct {
	config       *entities.Config
	logger       *zap.Logger
	projectDir   string
	settingsFile string
}

// NewSettingsService creates a new settings manager
func NewSettingsService(config *entities.Config, logger *zap.Logger) *SettingsService {
	wd, _ := os.Getwd()
	projectDir := wd

	return &SettingsService{
		config:       config,
		logger:       logger,
		projectDir:   projectDir,
		settingsFile: filepath.Join(projectDir, config.Flags.Config, "settings.json"),
	}
}

// GetDefaultSettings returns the default application settings
func (sm *SettingsService) GetDefaultSettings() *entities.Settings {
	defaultAutoAssignPattern := "TODO|FIXME"
	return &entities.Settings{
		KanbanColumns: []entities.KanbanColumn{
			{ID: "todo", Name: "TODO", Color: "dark", AutoAssignPattern: &defaultAutoAssignPattern},
			{ID: "in_progress", Name: "IN PROGRESS", Color: "blue"},
			{ID: "done", Name: "DONE", Color: "green"},
		},
		PriorityPatterns: entities.PriorityPatterns{
			Low:    "LOW",
			Medium: "MEDIUM",
			High:   "HIGH",
		},
		CodeScanSettings: entities.CodeScanConfig{
			ExcludeDirectories: []string{
				"node_modules",
				".git",
				".idea",
				".vscode",
				".cache",
				".next",
				"dist",
				"build",
				"out",
				"public",
				"vendor",
				"target",
				"tmp",
				"logs",
				"coverage",
			},

			ExcludeFiles: []string{
				"*.min.js",
				"*.min.css",
				"*.bundle.js",
				"*.png",
				"*.jpg",
				"*.jpeg",
				"*.gif",
				"*.webp",
				"*.svg",
				"*.ico",
				"*.map",
				"*.lock",
				"package-lock.json",
				"yarn.lock",
				"pnpm-lock.yaml",
				".gitignore",
				".dockerignore",
				"README.md",
				"LICENSE",
				"*.env",
				"*.local",
				"*.log",
				"*.db",
			},

			SyncEnabled: false,
		},
		GithubAuth: entities.GithubAuth{
			Token: "",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// LoadSettings loads settings from file or returns defaults
func (sm *SettingsService) LoadSettings() *entities.Settings {
	data, err := os.ReadFile(sm.settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			sm.logger.Info("Settings file not found, using defaults")
			return sm.GetDefaultSettings()
		}
		sm.logger.Error("Failed to read settings file", zap.Error(err))
		return sm.GetDefaultSettings()
	}

	var settings entities.Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		sm.logger.Error("Failed to unmarshal settings", zap.Error(err))
		return sm.GetDefaultSettings()
	}

	// Validate and fix settings if needed
	settings = *sm.validateSettings(&settings)

	return &settings
}

// SaveSettings saves settings to file
func (sm *SettingsService) SaveSettings(settings *entities.Settings) error {
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
func (sm *SettingsService) validateSettings(settings *entities.Settings) *entities.Settings {
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
func (sm *SettingsService) UpdatePartialSettings(updates map[string]interface{}) (*entities.Settings, error) {
	// Load current settings
	settings := sm.LoadSettings()

	// Apply updates
	if kanbanColumns, ok := updates["kanban_columns"]; ok {
		if columnsData, ok := kanbanColumns.([]interface{}); ok {
			var columns []entities.KanbanColumn
			for _, col := range columnsData {
				if colMap, ok := col.(map[string]interface{}); ok {
					column := entities.KanbanColumn{}
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
						column.AutoAssignPattern = &pattern
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
func (sm *SettingsService) GetSettingsSummary() map[string]interface{} {
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
