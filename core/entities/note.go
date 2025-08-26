package entities

import (
	"time"
)

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
	Status    *string   `json:"status,omitempty"`

	GitHubIssueNumber *int       `json:"github_issue_number,omitempty"`
	GitHubIssueURL    *string    `json:"github_issue_url,omitempty"`
	GitHubLastSync    *time.Time `json:"github_last_sync,omitempty"`
	GitHubSyncHash    *string    `json:"github_sync_hash,omitempty"`
}

type Folder struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
	Expanded bool   `json:"expanded"`
}

type Category struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type GitHubSyncConfig struct {
	Enabled         bool       `json:"enabled"`
	Token           string     `json:"token"`
	ImportNewIssues bool       `json:"import_new_issues"`
	SyncLabels      bool       `json:"sync_labels"`
	SyncAssignees   bool       `json:"sync_assignees"`
	AutoCloseIssues bool       `json:"auto_close_issues"`
	ExcludedLabels  []string   `json:"excluded_labels"`
	SyncInterval    int        `json:"sync_interval"`
	LastSync        *time.Time `json:"last_sync,omitempty"`
}

type SyncStats struct {
	LastSyncTime      time.Time `json:"last_sync_time"`
	TotalNotesSynced  int       `json:"total_notes_synced"`
	TotalIssuesSynced int       `json:"total_issues_synced"`
	NotesToIssues     int       `json:"notes_to_issues"`
	IssuesToNotes     int       `json:"issues_to_notes"`
	ConflictsResolved int       `json:"conflicts_resolved"`
	SyncErrors        int       `json:"sync_errors"`
}

type ConflictResolution string

const (
	ResolveManually  ConflictResolution = "manual"
	PreferLocal      ConflictResolution = "local"
	PreferRemote     ConflictResolution = "remote"
	PreferMostRecent ConflictResolution = "recent"
)

type ExtendedSettings struct {
	GitHubSync         GitHubSyncConfig   `json:"github_sync"`
	ConflictResolution ConflictResolution `json:"conflict_resolution"`
	AutoSync           bool               `json:"auto_sync"`
	SyncOnStartup      bool               `json:"sync_on_startup"`
}

type IssueTemplate struct {
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Labels    []string          `json:"labels"`
	Assignees []string          `json:"assignees"`
	Milestone *string           `json:"milestone,omitempty"`
	Variables map[string]string `json:"variables"`
}

type SyncEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"`
	NoteID      *int                   `json:"note_id,omitempty"`
	IssueNumber *int                   `json:"issue_number,omitempty"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Error       *string                `json:"error,omitempty"`
}

type NoteMetadata struct {
	Priority      *int                   `json:"priority,omitempty"`
	Difficulty    *string                `json:"difficulty,omitempty"`
	EstimatedTime *string                `json:"estimated_time,omitempty"`
	ActualTime    *string                `json:"actual_time,omitempty"`
	Dependencies  []int                  `json:"dependencies,omitempty"`
	References    []string               `json:"references,omitempty"`
	CustomFields  map[string]interface{} `json:"custom_fields,omitempty"`
}

type EnhancedNote struct {
	Note
	Metadata *NoteMetadata `json:"metadata,omitempty"`
}

type GitHubIssueMapping struct {
	NoteID           int                `json:"note_id"`
	IssueNumber      int                `json:"issue_number"`
	Repository       string             `json:"repository"`
	LastSyncHash     string             `json:"last_sync_hash"`
	SyncDirection    string             `json:"sync_direction"`
	ConflictStrategy ConflictResolution `json:"conflict_strategy"`
	SyncEnabled      bool               `json:"sync_enabled"`
	CustomMappings   map[string]string  `json:"custom_mappings"`
	LastSyncTime     time.Time          `json:"last_sync_time"`
	SyncErrors       []string           `json:"sync_errors,omitempty"`
}

type SyncProfile struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Repository       string             `json:"repository"`
	CategoryFilter   []string           `json:"category_filter,omitempty"`
	TagFilter        []string           `json:"tag_filter,omitempty"`
	FolderFilter     []int              `json:"folder_filter,omitempty"`
	IssueTemplate    *IssueTemplate     `json:"issue_template,omitempty"`
	LabelMapping     map[string]string  `json:"label_mapping,omitempty"`
	ConflictStrategy ConflictResolution `json:"conflict_strategy"`
	SyncInterval     int                `json:"sync_interval"`
	Enabled          bool               `json:"enabled"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
