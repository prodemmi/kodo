package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/prodemmi/kodo/core/entities"
	"go.uber.org/zap"
)

// ItemHistoryService handles project statistics and persistence
type ItemHistoryService struct {
	config      *entities.Config
	logger      *zap.Logger
	kodoDir     string
	statsFile   string
	historyFile string
}

// NewProjectTracker creates a new project itemHistoryService
func NewProjectTracker(config *entities.Config, logger *zap.Logger) *ItemHistoryService {
	wd, _ := os.Getwd()
	kodoDir := filepath.Join(wd, config.Flags.Config)

	return &ItemHistoryService{
		config:      config,
		logger:      logger,
		kodoDir:     kodoDir,
		statsFile:   filepath.Join(kodoDir, "items_history.json"),
		historyFile: filepath.Join(kodoDir, "branch_history.json"),
	}
}

// Initialize creates the .kodo directory and loads existing history
func (pt *ItemHistoryService) Initialize() error {
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

# Keep the history but ignore temporary data
!items_history.json
!branch_history.json
`
		if err := os.WriteFile(gitignoreFile, []byte(gitignoreContent), 0644); err != nil {
			pt.logger.Warn("Failed to create .gitignore", zap.Error(err))
		}
	}

	pt.logger.Info("Project itemHistoryService initialized", zap.String("kodo_dir", pt.kodoDir))
	return nil
}

// generateItemHash creates a unique hash for tracking item identity across commits
func (pt *ItemHistoryService) generateItemHash(item *entities.Item) string {
	// Create hash based on file, line, type, and title (content that shouldn't change)
	content := fmt.Sprintf("%s:%d:%s:%s", item.File, item.Line, item.Type, item.Title)
	// Simple hash - in production you might want to use crypto/sha256
	hash := 0
	for _, char := range content {
		hash = hash*31 + int(char)
	}
	return fmt.Sprintf("%x", hash)
}

// GetTaskItemsAnalysis provides detailed analysis of task items dynamically based on settings
func (pt *ItemHistoryService) GetTaskItemsAnalysis(settings *SettingsService) map[string]interface{} {
	history := pt.LoadStats()
	if history == nil {
		return map[string]interface{}{
			"error": "No history available",
		}
	}

	currentSettings := settings.LoadSettings()

	// Prepare dynamic priority key
	highPriorityKey := entities.ItemPriority(currentSettings.PriorityPatterns.High)

	analysis := map[string]interface{}{
		"total_items":     len(history.CurrentItems),
		"items_by_file":   make(map[string][]entities.TaskItem),
		"items_by_type":   make(map[string][]entities.TaskItem),
		"items_by_status": make(map[string][]entities.TaskItem),
		"high_priority":   []entities.TaskItem{},
		"recent_changes":  pt.GetRecentItemChanges(),
	}

	itemsByFile := analysis["items_by_file"].(map[string][]entities.TaskItem)
	itemsByType := analysis["items_by_type"].(map[string][]entities.TaskItem)
	itemsByStatus := analysis["items_by_status"].(map[string][]entities.TaskItem)
	highPriority := analysis["high_priority"].([]entities.TaskItem)

	// Build dynamic status keys from KanbanColumns
	statusKeys := make(map[entities.ItemStatus]string)
	for _, col := range currentSettings.KanbanColumns {
		key := entities.ItemStatus(strings.ToUpper(col.Name))
		statusKeys[key] = strings.ToLower(strings.ReplaceAll(col.Name, " ", "_"))
	}

	// Iterate items and group dynamically
	for _, item := range history.CurrentItems {
		// Group by file
		itemsByFile[item.File] = append(itemsByFile[item.File], item)

		// Group by type
		itemsByType[string(item.Type)] = append(itemsByType[string(item.Type)], item)

		// Group by status dynamically
		statusKey, ok := statusKeys[item.Status]
		if !ok {
			statusKey = strings.ToLower(string(item.Status)) // fallback
		}
		itemsByStatus[statusKey] = append(itemsByStatus[statusKey], item)

		// Collect high priority items dynamically
		if item.Priority == highPriorityKey {
			highPriority = append(highPriority, item)
		}
	}

	analysis["items_by_file"] = itemsByFile
	analysis["items_by_type"] = itemsByType
	analysis["items_by_status"] = itemsByStatus
	analysis["high_priority"] = highPriority

	return analysis
}

// getRecentItemChanges compares current items with previous snapshot
func (pt *ItemHistoryService) GetRecentItemChanges() map[string]interface{} {
	history := pt.GetBranchHistory()
	if len(history) < 2 {
		return map[string]interface{}{
			"message": "Not enough history for comparison",
		}
	}

	current := history[len(history)-1]
	previous := history[len(history)-2]

	// Create maps for easier lookup
	currentItems := make(map[string]entities.TaskItem)
	previousItems := make(map[string]entities.TaskItem)

	for _, item := range current.History.Items {
		currentItems[item.Hash] = item
	}

	for _, item := range previous.History.Items {
		previousItems[item.Hash] = item
	}

	var added []entities.TaskItem
	var removed []entities.TaskItem
	var statusChanged []map[string]interface{}

	// Find added items
	for hash, item := range currentItems {
		if _, exists := previousItems[hash]; !exists {
			added = append(added, item)
		} else {
			// Check for status changes
			prevItem := previousItems[hash]
			if item.Status != prevItem.Status {
				statusChanged = append(statusChanged, map[string]interface{}{
					"item":       item,
					"old_status": prevItem.Status,
					"new_status": item.Status,
				})
			}
		}
	}

	// Find removed items
	for hash, item := range previousItems {
		if _, exists := currentItems[hash]; !exists {
			removed = append(removed, item)
		}
	}

	return map[string]interface{}{
		"added":          added,
		"removed":        removed,
		"status_changed": statusChanged,
		"summary": map[string]int{
			"added":          len(added),
			"removed":        len(removed),
			"status_changed": len(statusChanged),
		},
	}
}

// GetItemsByFile returns items grouped by file with additional metadata
func (pt *ItemHistoryService) GetItemsByFile(settings *SettingsService) map[string]interface{} {
	history := pt.LoadStats()
	if history == nil {
		return map[string]interface{}{
			"error": "No history available",
		}
	}

	currentSettings := settings.LoadSettings()

	// Prepare dynamic status keys from KanbanColumns
	statusKeys := make(map[entities.ItemStatus]string)
	for _, col := range currentSettings.KanbanColumns {
		key := entities.ItemStatus(strings.ToUpper(col.Name))
		statusKeys[key] = strings.ToLower(strings.ReplaceAll(col.Name, " ", "_"))
	}

	// Prepare dynamic high priority key
	highPriorityKey := entities.ItemPriority(currentSettings.PriorityPatterns.High)

	fileGroups := make(map[string]map[string]interface{})

	for _, item := range history.CurrentItems {
		if _, exists := fileGroups[item.File]; !exists {
			// Initialize dynamic status counts
			statusCounts := map[string]int{}
			for _, k := range statusKeys {
				statusCounts[k] = 0
			}

			fileGroups[item.File] = map[string]interface{}{
				"items":         []entities.TaskItem{},
				"total":         0,
				"high_priority": 0,
				"status_counts": statusCounts, // dynamic status counters
			}
		}

		group := fileGroups[item.File]

		// Append item
		group["items"] = append(group["items"].([]entities.TaskItem), item)
		group["total"] = group["total"].(int) + 1

		// Update dynamic status count
		statusKey, ok := statusKeys[item.Status]
		if ok {
			group["status_counts"].(map[string]int)[statusKey]++
		}

		// Update high priority count dynamically
		if item.Priority == highPriorityKey {
			group["high_priority"] = group["high_priority"].(int) + 1
		}
	}

	return map[string]interface{}{
		"files":       fileGroups,
		"total_files": len(fileGroups),
	}
}

// GetItemTrends analyzes trends over time
func (pt *ItemHistoryService) GetItemTrends(settings *SettingsService) map[string]interface{} {
	history := pt.GetBranchHistory()
	if len(history) < 2 {
		return map[string]interface{}{
			"error": "Not enough history for trend analysis",
		}
	}

	currentSettings := settings.LoadSettings()
	kanbanCols := currentSettings.KanbanColumns
	if len(kanbanCols) == 0 {
		return map[string]interface{}{
			"error": "No Kanban columns configured",
		}
	}

	doneColumnID := kanbanCols[len(kanbanCols)-1].ID

	statusKeys := make(map[string]string)
	for _, col := range kanbanCols {
		statusKeys[col.ID] = col.Name
	}

	trends := map[string]interface{}{
		"timeline":        []map[string]interface{}{},
		"completion_rate": []map[string]interface{}{},
		"type_trends":     make(map[string][]map[string]interface{}),
	}

	for _, snapshot := range history {
		timelineEntry := map[string]interface{}{
			"timestamp": snapshot.Timestamp,
			"commit":    snapshot.CommitShort,
			"branch":    snapshot.Branch,
			"total":     snapshot.History.Total,
		}

		for statusID, name := range statusKeys {
			count := 0
			for _, item := range snapshot.History.Items {
				if string(item.Status) == statusID {
					count++
				} else if statusID == doneColumnID && item.IsDone {
					count++
				}
			}
			timelineEntry[name] = count
		}

		trends["timeline"] = append(trends["timeline"].([]map[string]interface{}), timelineEntry)

		doneCount := 0
		for _, item := range snapshot.History.Items {
			if item.IsDone {
				doneCount++
			}
		}

		completionRate := 0.0
		if snapshot.History.Total > 0 {
			completionRate = float64(doneCount) / float64(snapshot.History.Total) * 100
		}

		trends["completion_rate"] = append(trends["completion_rate"].([]map[string]interface{}), map[string]interface{}{
			"timestamp": snapshot.Timestamp,
			"commit":    snapshot.CommitShort,
			"rate":      completionRate,
		})

		// Track type trends
		for itemType, count := range snapshot.History.ByType {
			if trends["type_trends"].(map[string][]map[string]interface{})[itemType] == nil {
				trends["type_trends"].(map[string][]map[string]interface{})[itemType] = []map[string]interface{}{}
			}

			typeEntry := map[string]interface{}{
				"timestamp": snapshot.Timestamp,
				"commit":    snapshot.CommitShort,
				"count":     count,
			}
			trends["type_trends"].(map[string][]map[string]interface{})[itemType] = append(
				trends["type_trends"].(map[string][]map[string]interface{})[itemType],
				typeEntry,
			)
		}
	}

	return trends
}

// SaveStats saves current project statistics
func (pt *ItemHistoryService) SaveStats(items []*entities.Item, settings *SettingsService) error {
	if err := pt.Initialize(); err != nil {
		return err
	}

	history, err := pt.generateStats(items)
	if err != nil {
		return fmt.Errorf("failed to generate history: %v", err)
	}

	// Load existing history to preserve history
	existingStats := pt.LoadStats()
	if existingStats != nil {
		history.CreatedAt = existingStats.CreatedAt
		history.BranchHistory = existingStats.BranchHistory
	} else {
		history.CreatedAt = time.Now()
	}

	// Update the branch history if we're on a different commit
	if existingStats == nil || existingStats.GitCommit != history.GitCommit {
		snapshot := entities.BranchSnapshot{
			Branch:        history.GitBranch,
			Commit:        history.GitCommit,
			CommitShort:   history.GitCommitShort,
			CommitMessage: pt.getCommitMessage(history.GitCommit),
			Timestamp:     time.Now(),
			History:       pt.generateItemStats(items, settings),
		}

		history.BranchHistory = append(history.BranchHistory, snapshot)

		// Keep only last 50 snapshots to prevent file from growing too large
		if len(history.BranchHistory) > 50 {
			history.BranchHistory = history.BranchHistory[len(history.BranchHistory)-50:]
		}
	}

	history.UpdatedAt = time.Now()

	// Save to file
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	if err := os.WriteFile(pt.statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	pt.logger.Info("Project history saved",
		zap.String("file", pt.statsFile),
		zap.Int("total_items", history.TotalItems),
		zap.String("branch", history.GitBranch),
		zap.String("commit", history.GitCommitShort))

	return nil
}

// LoadStats loads project statistics from file
func (pt *ItemHistoryService) LoadStats() *entities.ItemsHistory {
	data, err := os.ReadFile(pt.statsFile)
	if err != nil {
		if !os.IsNotExist(err) {
			pt.logger.Warn("Failed to read history file", zap.Error(err))
		}
		return nil
	}

	var history entities.ItemsHistory
	if err := json.Unmarshal(data, &history); err != nil {
		pt.logger.Error("Failed to unmarshal history", zap.Error(err))
		return nil
	}

	return &history
}

// generateStats creates current project statistics
func (pt *ItemHistoryService) generateStats(items []*entities.Item) (*entities.ItemsHistory, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Get git information
	gitBranch := pt.GetGitBranch()
	gitCommit := pt.GetGitCommit()
	gitCommitShort := pt.GetGitCommitShort()

	// Generate item statistics
	itemsByStatus := make(map[string]int)
	itemsByType := make(map[string]int)
	itemsByFile := make(map[string]int)
	var taskItems []entities.TaskItem

	for _, item := range items {
		// Count by status
		itemsByStatus[string(item.Status)]++

		// Count by type
		itemsByType[string(item.Type)]++

		// Count by file
		itemsByFile[item.File]++

		// Create TaskItem for detailed tracking
		taskItem := entities.TaskItem{
			ID:       item.ID,
			Type:     item.Type,
			Title:    item.Title,
			File:     item.File,
			Line:     item.Line,
			Status:   item.Status,
			Priority: item.Priority,
			Hash:     pt.generateItemHash(item),
			IsDone:   item.IsDone,
			DoneAt:   item.DoneAt,
			DoneBy:   item.DoneBy,
		}
		taskItems = append(taskItems, taskItem)
	}

	return &entities.ItemsHistory{
		ProjectPath:    wd,
		LastScanAt:     time.Now(),
		GitBranch:      gitBranch,
		GitCommit:      gitCommit,
		GitCommitShort: gitCommitShort,
		TotalItems:     len(items),
		ItemsByStatus:  itemsByStatus,
		ItemsByType:    itemsByType,
		ItemsByFile:    itemsByFile,
		CurrentItems:   taskItems,
	}, nil
}

// generateItemStats creates ItemStats for history tracking
func (pt *ItemHistoryService) generateItemStats(items []*entities.Item, settings *SettingsService) entities.ItemStats {
	currentSettings := settings.LoadSettings()

	// Prepare dynamic status keys from KanbanColumns
	statusKeys := make(map[entities.ItemStatus]struct{})
	for _, col := range currentSettings.KanbanColumns {
		statusKeys[entities.ItemStatus(col.Name)] = struct{}{}
	}

	byStatus := make(map[string]int)
	byPriority := make(map[string]int)
	byType := make(map[string]int)

	history := entities.ItemStats{
		Total:      len(items),
		ByType:     byType,
		ByPriority: byPriority,
		Items:      make([]entities.TaskItem, 0, len(items)),
	}

	for _, item := range items {
		// Update status count dynamically
		statusStr := string(item.Status)
		if _, ok := byStatus[statusStr]; ok {
			byStatus[statusStr]++
		} else {
			byStatus[statusStr] = 1
		}

		// Update type count dynamically
		byType[string(item.Type)]++

		// Update priority count dynamically
		priorityStr := string(item.Priority)
		if _, ok := byPriority[priorityStr]; ok {
			byPriority[priorityStr]++
		} else {
			byPriority[priorityStr] = 1
		}

		// Add detailed TaskItem
		taskItem := entities.TaskItem{
			ID:       item.ID,
			Type:     item.Type,
			Title:    item.Title,
			File:     item.File,
			Line:     item.Line,
			Status:   item.Status,
			Priority: item.Priority,
			Hash:     pt.generateItemHash(item),
			IsDone:   item.IsDone,
			DoneAt:   item.DoneAt,
			DoneBy:   item.DoneBy,
		}
		history.Items = append(history.Items, taskItem)
	}

	history.ByStatus = byStatus
	history.ByPriority = byPriority

	return history
}

// GetProjectStats returns a summary of current history
func (pt *ItemHistoryService) GetProjectStats(settings *SettingsService) map[string]interface{} {
	history := pt.LoadStats()
	if history == nil {
		return map[string]interface{}{
			"error": "No history available",
		}
	}

	currentSettings := settings.LoadSettings()

	itemsByStatus := make(map[string]int)
	for _, col := range currentSettings.KanbanColumns {
		itemsByStatus[col.ID] = 0
	}

	total := history.TotalItems
	lastStatusID := currentSettings.KanbanColumns[len(currentSettings.KanbanColumns)-1].ID

	for _, item := range history.CurrentItems {
		if item.IsDone {
			itemsByStatus[lastStatusID]++
		} else {
			statusID := string(item.ID)
			if _, ok := itemsByStatus[statusID]; ok {
				itemsByStatus[statusID]++
			} else {
				itemsByStatus[statusID] = 1
			}
		}
	}

	progressPercent := 0.0
	if total > 0 {
		progressPercent = float64(itemsByStatus[lastStatusID]) / float64(total) * 100
	}

	return map[string]interface{}{
		"project_path":     history.ProjectPath,
		"last_scan":        history.LastScanAt,
		"git_branch":       history.GitBranch,
		"git_commit_short": history.GitCommitShort,
		"total_items":      total,
		"items_by_status":  itemsByStatus,
		"progress_percent": progressPercent,
		"items_by_type":    history.ItemsByType,
		"items_by_file":    history.ItemsByFile,
		"history_count":    len(history.BranchHistory),
		"created_at":       history.CreatedAt,
		"updated_at":       history.UpdatedAt,
	}
}

// GetBranchHistory returns the branch history
func (pt *ItemHistoryService) GetBranchHistory() []entities.BranchSnapshot {
	history := pt.LoadStats()
	if history == nil {
		return nil
	}

	slices.SortFunc(history.BranchHistory, func(a, b entities.BranchSnapshot) int {
		if a.Timestamp.Before(b.Timestamp) {
			return 1
		}
		if a.Timestamp.After(b.Timestamp) {
			return -1
		}
		return 0
	})

	return history.BranchHistory
}

// Git helper methods
func (pt *ItemHistoryService) GetGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git branch", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ItemHistoryService) GetGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ItemHistoryService) GetGitCommitShort() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit short", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *ItemHistoryService) getCommitMessage(commit string) string {
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

// CompareWithPreviousCommit compares current history with previous commit
func (pt *ItemHistoryService) CompareWithPreviousCommit(settings *SettingsService) map[string]interface{} {
	history := pt.GetBranchHistory()
	if len(history) < 2 {
		return map[string]interface{}{
			"error": "Not enough history for comparison",
		}
	}

	current := history[len(history)-1]
	previous := history[len(history)-2]

	currentSettings := settings.LoadSettings()
	kanbanCols := currentSettings.KanbanColumns
	if len(kanbanCols) == 0 {
		return map[string]interface{}{
			"error": "No Kanban columns configured",
		}
	}

	doneColumnID := kanbanCols[len(kanbanCols)-1].ID

	changes := make(map[string]int)
	for _, col := range kanbanCols {
		currentCount := 0
		previousCount := 0

		for _, item := range current.History.Items {
			if string(item.Status) == col.ID {
				currentCount++
			} else if col.ID == doneColumnID && item.IsDone {
				currentCount++
			}
		}

		for _, item := range previous.History.Items {
			if string(item.Status) == col.ID {
				previousCount++
			} else if col.ID == doneColumnID && item.IsDone {
				previousCount++
			}
		}

		changes[col.Name] = currentCount - previousCount
	}

	return map[string]interface{}{
		"current": map[string]interface{}{
			"commit":    current.CommitShort,
			"branch":    current.Branch,
			"timestamp": current.Timestamp,
			"history":   current.History,
		},
		"previous": map[string]interface{}{
			"commit":    previous.CommitShort,
			"branch":    previous.Branch,
			"timestamp": previous.Timestamp,
			"history":   previous.History,
		},
		"changes": changes,
	}
}

// CleanupOldStats removes old statistics (keeps last 30 days)
func (pt *ItemHistoryService) CleanupOldStats() error {
	history := pt.LoadStats()
	if history == nil {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -30) // 30 days ago
	var filteredHistory []entities.BranchSnapshot

	for _, snapshot := range history.BranchHistory {
		if snapshot.Timestamp.After(cutoff) {
			filteredHistory = append(filteredHistory, snapshot)
		}
	}

	if len(filteredHistory) != len(history.BranchHistory) {
		history.BranchHistory = filteredHistory
		history.UpdatedAt = time.Now()

		data, err := json.MarshalIndent(history, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(pt.statsFile, data, 0644); err != nil {
			return err
		}

		pt.logger.Info("Cleaned up old history",
			zap.Int("removed", len(history.BranchHistory)-len(filteredHistory)),
			zap.Int("remaining", len(filteredHistory)))
	}

	return nil
}

func (s *ItemHistoryService) CompareWithPrevious(settings *SettingsService) map[string]interface{} {
	return s.CompareWithPreviousCommit(settings)
}
