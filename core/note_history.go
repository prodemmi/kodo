package core

import (
	"time"
)

// NoteHistoryAction represents the type of action performed on a note
type NoteHistoryAction string

const (
	ActionCreated NoteHistoryAction = "created"
	ActionUpdated NoteHistoryAction = "updated"
	ActionDeleted NoteHistoryAction = "deleted"
	ActionMoved   NoteHistoryAction = "moved"
	ActionTagged  NoteHistoryAction = "tagged"
)

// NoteHistoryEntry represents a single history entry for a note
type NoteHistoryEntry struct {
	ID        int               `json:"id"`
	NoteID    int               `json:"note_id"`
	Action    NoteHistoryAction `json:"action"`
	Author    string            `json:"author"`
	Timestamp time.Time         `json:"timestamp"`
	GitBranch *string           `json:"git_branch,omitempty"`
	GitCommit *string           `json:"git_commit,omitempty"`

	// Change details
	Changes  map[string]interface{} `json:"changes,omitempty"`
	OldValue interface{}            `json:"old_value,omitempty"`
	NewValue interface{}            `json:"new_value,omitempty"`

	// Additional context
	Message  string                 `json:"message,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NoteHistoryFilter represents filters for querying note history
type NoteHistoryFilter struct {
	NoteID    *int               `json:"note_id,omitempty"`
	Action    *NoteHistoryAction `json:"action,omitempty"`
	Author    *string            `json:"author,omitempty"`
	GitBranch *string            `json:"git_branch,omitempty"`
	Since     *time.Time         `json:"since,omitempty"`
	Until     *time.Time         `json:"until,omitempty"`
	Limit     int                `json:"limit,omitempty"`
	Offset    int                `json:"offset,omitempty"`
}

// NoteHistoryStats represents statistics about note history
type NoteHistoryStats struct {
	TotalEntries    int                       `json:"total_entries"`
	ByAction        map[NoteHistoryAction]int `json:"by_action"`
	ByAuthor        map[string]int            `json:"by_author"`
	ByBranch        map[string]int            `json:"by_branch"`
	ByDay           map[string]int            `json:"by_day"`
	MostActiveNotes []NoteActivitySummary     `json:"most_active_notes"`
	RecentActivity  []NoteHistoryEntry        `json:"recent_activity"`
}

// NoteActivitySummary represents activity summary for a note
type NoteActivitySummary struct {
	NoteID      int       `json:"note_id"`
	NoteTitle   string    `json:"note_title"`
	ActionCount int       `json:"action_count"`
	LastAction  time.Time `json:"last_action"`
	Authors     []string  `json:"authors"`
}

// Enhanced Note structure with history tracking
type NoteWithHistory struct {
	Note
	History []NoteHistoryEntry `json:"history,omitempty"`
}

// Update the existing NoteStorage to include history
type EnhancedNoteStorage struct {
	Notes         []Note             `json:"notes"`
	Folders       []Folder           `json:"folders"`
	History       []NoteHistoryEntry `json:"history"`
	NextID        int                `json:"next_id"`
	NextHistoryID int                `json:"next_history_id"`
}
