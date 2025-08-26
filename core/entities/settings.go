package entities

import "time"

type Settings struct {
	KanbanColumns    []KanbanColumn   `json:"kanban_columns"`
	PriorityPatterns PriorityPatterns `json:"priority_patterns"`

	CodeScanSettings CodeScanConfig `json:"code_scan_settings"`
	GithubAuth       GithubAuth     `json:"github_auth"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KanbanColumn struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Color             string  `json:"color"`
	AutoAssignPattern *string `json:"auto_assign_pattern,omitempty"`
}

type PriorityPatterns struct {
	Low    string `json:"low"`
	Medium string `json:"medium"`
	High   string `json:"high"`
}

type CodeScanConfig struct {
	ExcludeDirectories []string `json:"exclude_directories"`
	ExcludeFiles       []string `json:"exclude_files"`
	SyncEnabled        bool     `json:"sync_enabled"`
}

type GithubAuth struct {
	Token string `json:"token"`
}
