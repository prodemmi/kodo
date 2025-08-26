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

type HistoryService struct {
	config      *entities.Config
	logger      *zap.Logger
	kodoDir     string
	statsFile   string
	historyFile string
}

func NewHistoryService(config *entities.Config, logger *zap.Logger) *HistoryService {
	wd, _ := os.Getwd()
	kodoDir := filepath.Join(wd, config.Flags.Config)

	return &HistoryService{
		config:      config,
		logger:      logger,
		kodoDir:     kodoDir,
		statsFile:   filepath.Join(kodoDir, "items_history.json"),
		historyFile: filepath.Join(kodoDir, "branch_history.json"),
	}
}

func (pt *HistoryService) Initialize() error {

	if err := os.MkdirAll(pt.kodoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .kodo directory: %v", err)
	}

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

func (pt *HistoryService) generateItemHash(item *entities.Item) string {

	content := fmt.Sprintf("%s:%d:%s:%s", item.File, item.Line, item.Type, item.Title)
	hash := 0
	for _, char := range content {
		hash = hash*31 + int(char)
	}
	return fmt.Sprintf("%x", hash)
}

func (pt *HistoryService) GetTaskItemsAnalysis(settings *SettingsService) map[string]interface{} {
	history := pt.LoadStats()
	if history == nil {
		return map[string]interface{}{
			"error": "No history available",
		}
	}

	currentSettings := settings.LoadSettings()

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

	statusKeys := make(map[entities.ItemStatus]string)
	for _, col := range currentSettings.KanbanColumns {
		key := entities.ItemStatus(strings.ToUpper(col.Name))
		statusKeys[key] = strings.ToLower(strings.ReplaceAll(col.Name, " ", "_"))
	}

	for _, item := range history.CurrentItems {

		itemsByFile[item.File] = append(itemsByFile[item.File], item)

		itemsByType[string(item.Type)] = append(itemsByType[string(item.Type)], item)

		statusKey, ok := statusKeys[item.Status]
		if !ok {
			statusKey = strings.ToLower(string(item.Status))
		}
		itemsByStatus[statusKey] = append(itemsByStatus[statusKey], item)

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

func (pt *HistoryService) GetRecentItemChanges() map[string]interface{} {
	history := pt.GetBranchHistory()
	if len(history) < 2 {
		return map[string]interface{}{
			"message": "Not enough history for comparison",
		}
	}

	current := history[len(history)-1]
	previous := history[len(history)-2]

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

	for hash, item := range currentItems {
		if _, exists := previousItems[hash]; !exists {
			added = append(added, item)
		} else {

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

func (pt *HistoryService) GetItemsByFile(settings *SettingsService) map[string]interface{} {
	history := pt.LoadStats()
	if history == nil {
		return map[string]interface{}{
			"error": "No history available",
		}
	}

	currentSettings := settings.LoadSettings()

	statusKeys := make(map[entities.ItemStatus]string)
	for _, col := range currentSettings.KanbanColumns {
		key := entities.ItemStatus(strings.ToUpper(col.Name))
		statusKeys[key] = strings.ToLower(strings.ReplaceAll(col.Name, " ", "_"))
	}

	highPriorityKey := entities.ItemPriority(currentSettings.PriorityPatterns.High)

	fileGroups := make(map[string]map[string]interface{})

	for _, item := range history.CurrentItems {
		if _, exists := fileGroups[item.File]; !exists {

			statusCounts := map[string]int{}
			for _, k := range statusKeys {
				statusCounts[k] = 0
			}

			fileGroups[item.File] = map[string]interface{}{
				"items":         []entities.TaskItem{},
				"total":         0,
				"high_priority": 0,
				"status_counts": statusCounts,
			}
		}

		group := fileGroups[item.File]

		group["items"] = append(group["items"].([]entities.TaskItem), item)
		group["total"] = group["total"].(int) + 1

		statusKey, ok := statusKeys[item.Status]
		if ok {
			group["status_counts"].(map[string]int)[statusKey]++
		}

		if item.Priority == highPriorityKey {
			group["high_priority"] = group["high_priority"].(int) + 1
		}
	}

	return map[string]interface{}{
		"files":       fileGroups,
		"total_files": len(fileGroups),
	}
}

func (pt *HistoryService) GetItemTrends(settings *SettingsService) map[string]interface{} {
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

func (pt *HistoryService) SaveStats(items []*entities.Item, settings *SettingsService) error {
	if err := pt.Initialize(); err != nil {
		return err
	}

	history, err := pt.generateStats(items)
	if err != nil {
		return fmt.Errorf("failed to generate history: %v", err)
	}

	existingStats := pt.LoadStats()
	if existingStats != nil {
		history.CreatedAt = existingStats.CreatedAt
		history.BranchHistory = existingStats.BranchHistory
	} else {
		history.CreatedAt = time.Now()
	}

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

		if len(history.BranchHistory) > 50 {
			history.BranchHistory = history.BranchHistory[len(history.BranchHistory)-50:]
		}
	}

	history.UpdatedAt = time.Now()

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

func (pt *HistoryService) LoadStats() *entities.ItemsHistory {
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

func (pt *HistoryService) generateStats(items []*entities.Item) (*entities.ItemsHistory, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gitBranch := pt.GetGitBranch()
	gitCommit := pt.GetGitCommit()
	gitCommitShort := pt.GetGitCommitShort()

	itemsByStatus := make(map[string]int)
	itemsByType := make(map[string]int)
	itemsByFile := make(map[string]int)
	var taskItems []entities.TaskItem

	for _, item := range items {

		itemsByStatus[string(item.Status)]++

		itemsByType[string(item.Type)]++

		itemsByFile[item.File]++

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

func (pt *HistoryService) generateItemStats(items []*entities.Item, settings *SettingsService) entities.ItemStats {
	currentSettings := settings.LoadSettings()

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

		statusStr := string(item.Status)
		if _, ok := byStatus[statusStr]; ok {
			byStatus[statusStr]++
		} else {
			byStatus[statusStr] = 1
		}

		byType[string(item.Type)]++

		priorityStr := string(item.Priority)
		if _, ok := byPriority[priorityStr]; ok {
			byPriority[priorityStr]++
		} else {
			byPriority[priorityStr] = 1
		}

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

func (pt *HistoryService) GetProjectStats(settings *SettingsService) map[string]interface{} {
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

func (pt *HistoryService) GetBranchHistory() []entities.BranchSnapshot {
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

func (pt *HistoryService) GetGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git branch", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *HistoryService) GetGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *HistoryService) GetGitCommitShort() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		pt.logger.Debug("Failed to get git commit short", zap.Error(err))
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (pt *HistoryService) getCommitMessage(commit string) string {
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

	if lines := strings.Split(message, "\n"); len(lines) > 0 {
		return lines[0]
	}
	return message
}

func (pt *HistoryService) CompareWithPreviousCommit(settings *SettingsService) map[string]interface{} {
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

func (pt *HistoryService) CleanupOldStats() error {
	history := pt.LoadStats()
	if history == nil {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -30)
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

func (s *HistoryService) CompareWithPrevious(settings *SettingsService) map[string]interface{} {
	return s.CompareWithPreviousCommit(settings)
}
