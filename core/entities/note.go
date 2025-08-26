package entities

import (
	"time"
)

// Note represents a note with GitHub integration
type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Category  string    `json:"category"`
	FolderID  *int      `json:"folder_id,omitempty"`
	GitBranch *string   `json:"git_branch,omitempty"`
	GitCommit *string   `json:"git_commit,omitempty"`
	Pinned    bool      `json:"pinned"`
	Status    *string   `json:"status,omitempty"` // open, closed, etc.

	// GitHub Integration fields
	GitHubIssueNumber *int       `json:"github_issue_number,omitempty"`
	GitHubIssueURL    *string    `json:"github_issue_url,omitempty"`
	GitHubLastSync    *time.Time `json:"github_last_sync,omitempty"`
	GitHubSyncHash    *string    `json:"github_sync_hash,omitempty"` // Hash to detect changes
}

// Folder represents a folder for organizing notes
type Folder struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
	Expanded bool   `json:"expanded"`
}

// Category represents a note category
type Category struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GitHub Sync Configuration
type GitHubSyncConfig struct {
	Enabled         bool       `json:"enabled"`
	Token           string     `json:"token"`
	ImportNewIssues bool       `json:"import_new_issues"`
	SyncLabels      bool       `json:"sync_labels"`
	SyncAssignees   bool       `json:"sync_assignees"`
	AutoCloseIssues bool       `json:"auto_close_issues"`
	ExcludedLabels  []string   `json:"excluded_labels"`
	SyncInterval    int        `json:"sync_interval"` // in minutes
	LastSync        *time.Time `json:"last_sync,omitempty"`
}

// Sync Statistics
type SyncStats struct {
	LastSyncTime      time.Time `json:"last_sync_time"`
	TotalNotesSynced  int       `json:"total_notes_synced"`
	TotalIssuesSynced int       `json:"total_issues_synced"`
	NotesToIssues     int       `json:"notes_to_issues"`
	IssuesToNotes     int       `json:"issues_to_notes"`
	ConflictsResolved int       `json:"conflicts_resolved"`
	SyncErrors        int       `json:"sync_errors"`
}

// Conflict Resolution Strategy
type ConflictResolution string

const (
	ResolveManually  ConflictResolution = "manual"
	PreferLocal      ConflictResolution = "local"
	PreferRemote     ConflictResolution = "remote"
	PreferMostRecent ConflictResolution = "recent"
)

// Extended Settings to include GitHub sync
type ExtendedSettings struct {
	GitHubSync         GitHubSyncConfig   `json:"github_sync"`
	ConflictResolution ConflictResolution `json:"conflict_resolution"`
	AutoSync           bool               `json:"auto_sync"`
	SyncOnStartup      bool               `json:"sync_on_startup"`
}

// Issue Template for creating standardized GitHub issues
type IssueTemplate struct {
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Labels    []string          `json:"labels"`
	Assignees []string          `json:"assignees"`
	Milestone *string           `json:"milestone,omitempty"`
	Variables map[string]string `json:"variables"` // For template substitution
}

// Sync Event for logging and monitoring
type SyncEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"` // "sync_start", "sync_complete", "conflict", "error"
	NoteID      *int                   `json:"note_id,omitempty"`
	IssueNumber *int                   `json:"issue_number,omitempty"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Error       *string                `json:"error,omitempty"`
}

// Note Metadata for extended properties
type NoteMetadata struct {
	Priority      *int                   `json:"priority,omitempty"`       // 1-5 scale
	Difficulty    *string                `json:"difficulty,omitempty"`     // easy, medium, hard
	EstimatedTime *string                `json:"estimated_time,omitempty"` // 1h, 2d, etc.
	ActualTime    *string                `json:"actual_time,omitempty"`
	Dependencies  []int                  `json:"dependencies,omitempty"` // Note IDs
	References    []string               `json:"references,omitempty"`   // URLs, file paths, etc.
	CustomFields  map[string]interface{} `json:"custom_fields,omitempty"`
}

// Enhanced Note structure with metadata
type EnhancedNote struct {
	Note
	Metadata *NoteMetadata `json:"metadata,omitempty"`
}

// GitHub Issue Mapping for complex sync scenarios
type GitHubIssueMapping struct {
	NoteID           int                `json:"note_id"`
	IssueNumber      int                `json:"issue_number"`
	Repository       string             `json:"repository"` // owner/repo format
	LastSyncHash     string             `json:"last_sync_hash"`
	SyncDirection    string             `json:"sync_direction"` // "bidirectional", "note_to_issue", "issue_to_note"
	ConflictStrategy ConflictResolution `json:"conflict_strategy"`
	SyncEnabled      bool               `json:"sync_enabled"`
	CustomMappings   map[string]string  `json:"custom_mappings"` // field mappings
	LastSyncTime     time.Time          `json:"last_sync_time"`
	SyncErrors       []string           `json:"sync_errors,omitempty"`
}

// Sync Profile for different sync configurations
type SyncProfile struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Repository       string             `json:"repository"`
	CategoryFilter   []string           `json:"category_filter,omitempty"` // Only sync these categories
	TagFilter        []string           `json:"tag_filter,omitempty"`      // Only sync notes with these tags
	FolderFilter     []int              `json:"folder_filter,omitempty"`   // Only sync notes in these folders
	IssueTemplate    *IssueTemplate     `json:"issue_template,omitempty"`
	LabelMapping     map[string]string  `json:"label_mapping,omitempty"` // note_tag -> github_label
	ConflictStrategy ConflictResolution `json:"conflict_strategy"`
	SyncInterval     int                `json:"sync_interval"` // minutes
	Enabled          bool               `json:"enabled"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
