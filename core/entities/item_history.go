package entities

import "time"

type ItemsHistory struct {
	ProjectPath    string           `json:"project_path"`
	LastScanAt     time.Time        `json:"last_scan_at"`
	GitBranch      string           `json:"git_branch"`
	GitCommit      string           `json:"git_commit"`
	GitCommitShort string           `json:"git_commit_short"`
	TotalItems     int              `json:"total_items"`
	ItemsByStatus  map[string]int   `json:"items_by_status"`
	ItemsByType    map[string]int   `json:"items_by_type"`
	ItemsByFile    map[string]int   `json:"items_by_file"`
	CurrentItems   []TaskItem       `json:"current_items"`
	BranchHistory  []BranchSnapshot `json:"branch_history,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type BranchSnapshot struct {
	Branch        string    `json:"branch"`
	Commit        string    `json:"commit"`
	CommitShort   string    `json:"commit_short"`
	CommitMessage string    `json:"commit_message"`
	Timestamp     time.Time `json:"timestamp"`
	History       ItemStats `json:"history"`
}

type ItemStats struct {
	Total      int            `json:"total"`
	ByStatus   map[string]int `json:"by_status"`
	ByType     map[string]int `json:"by_type"`
	ByPriority map[string]int `json:"by_priority"`
	Items      []TaskItem     `json:"items"`
}

type TaskItem struct {
	ID       int          `json:"id"`
	Type     ItemType     `json:"type"`
	Title    string       `json:"title"`
	File     string       `json:"file"`
	Line     int          `json:"line"`
	Status   ItemStatus   `json:"status"`
	Priority ItemPriority `json:"priority"`
	IsDone   bool         `json:"is_done"`
	DoneAt   *time.Time   `json:"done_at"`
	DoneBy   *string      `json:"done_by"`
	Hash     string       `json:"hash"`
}
