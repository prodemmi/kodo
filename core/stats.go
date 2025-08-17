package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ProjectStats represents the overall project statistics
type ProjectStats struct {
	ProjectPath    string           `json:"project_path"`
	LastScanAt     time.Time        `json:"last_scan_at"`
	GitBranch      string           `json:"git_branch"`
	GitCommit      string           `json:"git_commit"`
	GitCommitShort string           `json:"git_commit_short"`
	TotalItems     int              `json:"total_items"`
	ItemsByStatus  map[string]int   `json:"items_by_status"`
	ItemsByType    map[string]int   `json:"items_by_type"`
	ItemsByFile    map[string]int   `json:"items_by_file"`
	BranchHistory  []BranchSnapshot `json:"branch_history,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// BranchSnapshot represents stats at a specific git state
type BranchSnapshot struct {
	Branch        string    `json:"branch"`
	Commit        string    `json:"commit"`
	CommitShort   string    `json:"commit_short"`
	CommitMessage string    `json:"commit_message"`
	Timestamp     time.Time `json:"timestamp"`
	Stats         ItemStats `json:"stats"`
}

// ItemStats represents count statistics
type ItemStats struct {
	Total      int            `json:"total"`
	Todo       int            `json:"todo"`
	InProgress int            `json:"in_progress"`
	Done       int            `json:"done"`
	ByType     map[string]int `json:"by_type"`
	ByPriority map[string]int `json:"by_priority"`
}

// ProjectTracker handles project statistics and persistence
type ProjectTracker struct {
	scanner     *Scanner
	logger      *zap.Logger
	kodoDir     string
	statsFile   string
	historyFile string
}

// NewProjectTracker creates a new project tracker
func NewProjectTracker(scanner *Scanner, logger *zap.Logger) *ProjectTracker {
	wd, _ := os.Getwd()
	kodoDir := filepath.Join(wd, ".kodo")

	return &ProjectTracker{
		scanner:     scanner,
		logger:      logger,
		kodoDir:     kodoDir,
		statsFile:   filepath.Join(kodoDir, "project_stats.json"),
		historyFile: filepath.Join(kodoDir, "branch_history.json"),
	}
}

// Initialize creates the .kodo directory and loads existing stats
func (pt *ProjectTracker) Initialize() error {
	// Create .kodo directory if it doesn't exist
	if err := os.MkdirAll(pt.kodoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .kodo directory: %v", err)
	}

	// Create .gitignore if it doesn't exist
	gitignoreFile := filepath.Join(pt.kodoDir, ".gitignore")
	if _, err := os.Stat(gitignoreFile); os.IsNotExist(err) {
		gitignoreContent := `# Kodo temporary files
*.tmp
*.log

# Keep the stats but ignore temporary data
!project_stats.json
!branch_history.json
`
		if err := os.WriteFile(gitignoreFile, []byte(gitignoreContent), 0644); err != nil {
			pt.logger.Warn("Failed to create .gitignore", zap.Error(err))
		}
	}

	pt.logger.Info("Project tracker initialized", zap.String("kodo_dir", pt.kodoDir))
	return nil
}

// SaveStats saves current project statistics
func (pt *ProjectTracker) SaveStats() error {
	if err := pt.Initialize(); err != nil {
		return err
	}

	stats, err := pt.generateStats()
	if err != nil {
		return fmt.Errorf("failed to generate stats: %v", err)
	}

	// Load existing stats to preserve history
	existingStats := pt.LoadStats()
	if existingStats != nil {
		stats.CreatedAt = existingStats.CreatedAt
		stats.BranchHistory = existingStats.BranchHistory
	} else {
		stats.CreatedAt = time.Now()
	}

	// Update the branch history if we're on a different commit
	if existingStats == nil || existingStats.GitCommit != stats.GitCommit {
		snapshot := BranchSnapshot{
			Branch:        stats.GitBranch,
			Commit:        stats.GitCommit,
			CommitShort:   stats.GitCommitShort,
			CommitMessage: pt.getCommitMessage(stats.GitCommit),
			Timestamp:     time.Now(),
			Stats:         pt.generateItemStats(),
		}

		stats.BranchHistory = append(stats.BranchHistory, snapshot)

		// Keep only last 50 snapshots to prevent file from growing too large
		if len(stats.BranchHistory) > 50 {
			stats.BranchHistory = stats.BranchHistory[len(stats.BranchHistory)-50:]
		}
	}

	stats.UpdatedAt = time.Now()

	// Save to file
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %v", err)
	}

	if err := os.WriteFile(pt.statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write stats file: %v", err)
	}

	pt.logger.Info("Project stats saved",
		zap.String("file", pt.statsFile),
		zap.Int("total_items", stats.TotalItems),
		zap.String("branch", stats.GitBranch),
		zap.String("commit", stats.GitCommitShort))

	return nil
}

// LoadStats loads project statistics from file
func (pt *ProjectTracker) LoadStats() *ProjectStats {
	data, err := os.ReadFile(pt.statsFile)
	if err != nil {
		if !os.IsNotExist(err) {
			pt.logger.Warn("Failed to read stats file", zap.Error(err))
		}
		return nil
	}

	var stats ProjectStats
	if err := json.Unmarshal(data, &stats); err != nil {
		pt.logger.Error("Failed to unmarshal stats", zap.Error(err))
		return nil
	}

	return &stats
}

// generateStats creates current project statistics
func (pt *ProjectTracker) generateStats() (*ProjectStats, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Get git information
	gitBranch := pt.getGitBranch()
	gitCommit := pt.getGitCommit()
	gitCommitShort := pt.getGitCommitShort()

	// Generate item statistics
	items := pt.scanner.GetItems()
	itemsByStatus := make(map[string]int)
	itemsByType := make(map[string]int)
	itemsByFile := make(map[string]int)

	for _, item := range items {
		// Count by status
		itemsByStatus[string(item.Status)]++

		// Count by type
		itemsByType[string(item.Type)]++

		// Count by file
		itemsByFile[item.File]++
	}

	return &ProjectStats{
		ProjectPath:    wd,
		LastScanAt:     time.Now(),
		GitBranch:      gitBranch,
		GitCommit:      gitCommit,
		GitCommitShort: gitCommitShort,
		TotalItems:     len(items),
		ItemsByStatus:  itemsByStatus,
		ItemsByType:    itemsByType,
		ItemsByFile:    itemsByFile,
	}, nil
}

// generateItemStats creates ItemStats for history tracking
func (pt *ProjectTracker) generateItemStats() ItemStats {
	items := pt.scanner.GetItems()
	stats := ItemStats{
		Total:      len(items),
		ByType:     make(map[string]int),
		ByPriority: make(map[string]int),
	}

	for _, item := range items {
		switch item.Status {
		case StatusTodo:
			stats.Todo++
		case StatusInProgress:
			stats.InProgress++
		case StatusDone:
			stats.Done++
		}

		stats.ByType[string(item.Type)]++
		stats.ByPriority[string(item.Priority)]++
	}

	return stats
}

// GetStatsSummary returns a summary of current stats
func (pt *ProjectTracker) GetStatsSummary() map[string]interface{} {
	stats := pt.LoadStats()
	if stats == nil {
		return map[string]interface{}{
			"error": "No stats available",
		}
	}

	// Calculate progress percentage
	total := stats.TotalItems
	done := stats.ItemsByStatus["done"]
	inProgress := stats.ItemsByStatus["in-progress"]
	todo := stats.ItemsByStatus["todo"]

	var progressPercent float64
	if total > 0 {
		progressPercent = float64(done) / float64(total) * 100
	}

	return map[string]interface{}{
		"project_path":     stats.ProjectPath,
		"last_scan":        stats.LastScanAt,
		"git_branch":       stats.GitBranch,
		"git_commit_short": stats.GitCommitShort,
		"total_items":      total,
		"done":             done,
		"in_progress":      inProgress,
		"todo":             todo,
		"progress_percent": progressPercent,
		"items_by_type":    stats.ItemsByType,
		"items_by_file":    stats.ItemsByFile,
		"history_count":    len(stats.BranchHistory),
		"created_at":       stats.CreatedAt,
		"updated_at":       stats.UpdatedAt,
	}
}

// GetBranchHistory returns the branch history
func (pt *ProjectTracker) GetBranchHistory() []BranchSnapshot {
	stats := pt.LoadStats()
	if stats == nil {
		return nil
	}
	return stats.BranchHistory
}

// Git helper methods
func (pt *ProjectTracker) getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git branch", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ProjectTracker) getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ProjectTracker) getGitCommitShort() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit short", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ProjectTracker) getCommitMessage(commit string) string {
	if commit == "unknown" || commit == "" {
		return ""
	}

	cmd := exec.Command("git", "log", "--format=%B", "-n", "1", commit)
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get commit message", zap.Error(err))
		return ""
	}

	message := strings.TrimSpace(string(output))
	// Return only the first line (subject)
	if lines := strings.Split(message, "\n"); len(lines) > 0 {
		return lines[0]
	}
	return message
}

// CompareWithPreviousCommit compares current stats with previous commit
func (pt *ProjectTracker) CompareWithPreviousCommit() map[string]interface{} {
	history := pt.GetBranchHistory()
	if len(history) < 2 {
		return map[string]interface{}{
			"error": "Not enough history for comparison",
		}
	}

	current := history[len(history)-1]
	previous := history[len(history)-2]

	comparison := map[string]interface{}{
		"current": map[string]interface{}{
			"commit":    current.CommitShort,
			"branch":    current.Branch,
			"timestamp": current.Timestamp,
			"stats":     current.Stats,
		},
		"previous": map[string]interface{}{
			"commit":    previous.CommitShort,
			"branch":    previous.Branch,
			"timestamp": previous.Timestamp,
			"stats":     previous.Stats,
		},
		"changes": map[string]interface{}{
			"total":       current.Stats.Total - previous.Stats.Total,
			"done":        current.Stats.Done - previous.Stats.Done,
			"in_progress": current.Stats.InProgress - previous.Stats.InProgress,
			"todo":        current.Stats.Todo - previous.Stats.Todo,
		},
	}

	return comparison
}

// CleanupOldStats removes old statistics (keeps last 30 days)
func (pt *ProjectTracker) CleanupOldStats() error {
	stats := pt.LoadStats()
	if stats == nil {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -30) // 30 days ago
	var filteredHistory []BranchSnapshot

	for _, snapshot := range stats.BranchHistory {
		if snapshot.Timestamp.After(cutoff) {
			filteredHistory = append(filteredHistory, snapshot)
		}
	}

	if len(filteredHistory) != len(stats.BranchHistory) {
		stats.BranchHistory = filteredHistory
		stats.UpdatedAt = time.Now()

		data, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(pt.statsFile, data, 0644); err != nil {
			return err
		}

		pt.logger.Info("Cleaned up old stats",
			zap.Int("removed", len(stats.BranchHistory)-len(filteredHistory)),
			zap.Int("remaining", len(filteredHistory)))
	}

	return nil
}
