package entities

import "time"

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
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Color             string  `json:"color"`
	AutoAssignPattern *string `json:"auto_assign_pattern,omitempty"`
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
