package core

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Scanner struct {
	Todos  []*Item
	todoID int

	tracker *ProjectTracker
}

func NewScanner(config Config, logger *zap.Logger) *Scanner {
	scanner := &Scanner{}
	scanner.tracker = NewProjectTracker(config, scanner, logger)
	return scanner
}

func (s *Scanner) GetItems() []*Item {
	return s.Todos
}

func (s *Scanner) GetItemsLength() int {
	return len(s.Todos)
}

func (s *Scanner) Rescan() {
	s.ScanTodos()

	// Save stats after scanning
	if err := s.tracker.SaveStats(); err != nil {
		// Don't fail the scan if stats saving fails, just log it
		fmt.Printf("Warning: Failed to save project stats: %v\n", err)
	}
}

func (s *Scanner) ScanTodos() {
	s.Todos = []*Item{}
	s.todoID = 0

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	// Create comprehensive pattern for all item types
	itemTypes := []string{
		"TODO", "FIXME", "BUG", "NOTE", "REFACTOR", "OPTIMIZE", "CLEANUP", "DEPRECATED",
		"FEATURE", "FEAT", "ENHANCE", "DOC", "TEST", "EXAMPLE", "SECURITY", "COMPLIANCE",
		"DEBT", "ARCHITECTURE", "ARCH", "CONFIG", "DEPLOY", "MONITOR", "QUESTION", "IDEA",
		"REVIEW", "HACK",
	}

	// Build dynamic pattern
	typePattern := strings.Join(itemTypes, "|")

	// Match first line of a TODO item
	todoPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s):\s*(.+)?`, typePattern))
	// Match continuation description lines
	descPattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*(.+)`)
	// Match DONE comments
	donePattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*DONE\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2})\s+by\s+(.+?)(\s*-->)?$`)
	// Match IN PROGRESS comments
	inProgressPattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*IN\s+PROGRESS\s+from\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2})\s+by\s+(.+?)(\s*-->)?$`)

	err = filepath.Walk(wd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		textExts := map[string]bool{
			".go": true, ".js": true, ".ts": true, ".html": true, ".css": true,
			".py": true, ".java": true, ".cpp": true, ".c": true, ".h": true,
			".rb": true, ".php": true, ".vue": true, ".jsx": true, ".tsx": true,
			".sql": true, ".yml": true, ".yaml": true, ".json": true, ".md": true,
			".rs": true, ".kt": true, ".swift": true, ".dart": true, ".scala": true,
		}
		if !textExts[ext] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if matches := todoPattern.FindStringSubmatch(line); len(matches) > 0 {
				s.todoID++
				itemType := ItemType(strings.ToUpper(matches[2]))
				title := strings.TrimSpace(matches[3])
				todoStartLine := lineNum

				// Collect following description lines and check for status comments
				var descriptions []string
				var history []StatusHistory
				currentStatus := StatusTodo

				// Get current user for initial creation
				currentUser := s.getCurrentUser()

				for scanner.Scan() {
					nextLine := scanner.Text()
					lineNum++

					// Check if this is a DONE comment
					if doneMatches := donePattern.FindStringSubmatch(nextLine); len(doneMatches) > 0 {
						if parsedTime, err := time.Parse("2006-01-02 15:04", doneMatches[2]); err == nil {
							history = append(history, StatusHistory{
								Status:    StatusDone,
								Timestamp: parsedTime,
								User:      strings.TrimSpace(doneMatches[3]),
							})
							currentStatus = StatusDone
						}
					} else if inProgressMatches := inProgressPattern.FindStringSubmatch(nextLine); len(inProgressMatches) > 0 {
						// Check if this is an IN PROGRESS comment
						if parsedTime, err := time.Parse("2006-01-02 15:04", inProgressMatches[2]); err == nil {
							history = append(history, StatusHistory{
								Status:    StatusInProgress,
								Timestamp: parsedTime,
								User:      strings.TrimSpace(inProgressMatches[3]),
							})
							// Only set as in-progress if not already done
							if currentStatus != StatusDone {
								currentStatus = StatusInProgress
							}
						}
					} else if descMatches := descPattern.FindStringSubmatch(nextLine); len(descMatches) > 0 {
						desc := strings.TrimSpace(descMatches[2])
						// Skip status lines that didn't match patterns
						upperDesc := strings.ToUpper(desc)
						if !strings.HasPrefix(upperDesc, "DONE ") &&
							!strings.HasPrefix(upperDesc, "IN PROGRESS ") {
							descriptions = append(descriptions, desc)
						}
					} else {
						// Not a comment line, break
						break
					}
				}

				relPath, _ := filepath.Rel(wd, path)

				// Create new item with proper history
				item := &Item{
					ID:          s.todoID,
					Type:        itemType,
					Title:       title,
					Description: strings.Join(descriptions, "\n"),
					File:        relPath,
					Line:        todoStartLine,
					Status:      currentStatus,
					Priority:    GetPriorityFromType(itemType),
					CreatedAt:   time.Now(), // We don't know actual creation time
					UpdatedAt:   time.Now(),
					CurrentUser: currentUser,
					History:     history,
				}

				// If no history was found, add initial creation entry
				if len(history) == 0 {
					item.History = []StatusHistory{{
						Status:    StatusTodo,
						Timestamp: item.CreatedAt,
						User:      currentUser,
					}}
				}

				s.Todos = append(s.Todos, item)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

// MarkAsInProgress marks the item as in progress, updates file comment
func (s *Scanner) MarkAsInProgress(item *Item) error {
	currentUser := s.getCurrentUser()

	// Update the item status first
	item.SetStatus(StatusInProgress, currentUser)

	// Then update the file comment
	return s.updateStatusComment(item)
}

// MarkAsDone marks the item as done, updates file comment
func (s *Scanner) MarkAsDone(item *Item) error {
	currentUser := s.getCurrentUser()

	// Update the item status first
	item.SetStatus(StatusDone, currentUser)

	// Then update the file comment
	return s.updateStatusComment(item)
}

// MarkAsUndone marks the item as undone, removes status comments
func (s *Scanner) MarkAsUndone(item *Item) error {
	currentUser := s.getCurrentUser()

	// Update the item status first
	item.SetStatus(StatusTodo, currentUser)

	// Then update the file comment
	return s.updateStatusComment(item)
}

// updateStatusComment - COMPLETELY REWRITTEN for better reliability
func (s *Scanner) updateStatusComment(item *Item) error {
	// Get working directory to construct full path
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	fullPath := filepath.Join(wd, item.File)

	// Read the file
	lines, err := s.readFileLines(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", fullPath, err)
	}

	if item.Line > len(lines) || item.Line < 1 {
		return fmt.Errorf("invalid line number: %d (file has %d lines)", item.Line, len(lines))
	}

	// Get comment prefix for this file type
	prefix := s.getCommentPrefix(item.File)

	// Patterns to match existing status comments
	donePattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*DONE\s+.*$`)
	inProgressPattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*IN\s+PROGRESS\s+.*$`)

	// Find the TODO comment block
	todoLineIndex := item.Line - 1 // Convert to 0-based index

	// Find the end of the comment block
	endIndex := todoLineIndex
	for i := todoLineIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if s.isCommentLine(line, prefix) &&
			!donePattern.MatchString(lines[i]) &&
			!inProgressPattern.MatchString(lines[i]) {
			endIndex = i
		} else if donePattern.MatchString(lines[i]) || inProgressPattern.MatchString(lines[i]) {
			// This is a status comment, include it for potential removal
			endIndex = i
		} else {
			// Not a comment line, end of block
			break
		}
	}

	// Build the new lines array
	var newLines []string

	// Add all lines before the TODO block
	newLines = append(newLines, lines[:todoLineIndex]...)

	// Add the original TODO comment and description lines (but not status comments)
	for i := todoLineIndex; i <= endIndex; i++ {
		line := lines[i]
		// Skip existing status comments
		if !donePattern.MatchString(line) && !inProgressPattern.MatchString(line) {
			newLines = append(newLines, line)
		}
	}

	// Add new status comment based on current status
	switch item.Status {
	case StatusDone:
		doneAt := item.GetDoneAt()
		doneBy := item.GetDoneBy()
		if doneAt != nil && doneBy != "" {
			statusComment := s.formatStatusComment(prefix, "DONE", doneAt.Format("2006-01-02 15:04"), doneBy)
			newLines = append(newLines, statusComment)
		}

	case StatusInProgress:
		inProgressAt := item.GetInProgressAt()
		inProgressBy := item.GetInProgressBy()
		if inProgressAt != nil && inProgressBy != "" {
			statusComment := s.formatStatusComment(prefix, "IN PROGRESS from", inProgressAt.Format("2006-01-02 15:04"), inProgressBy)
			newLines = append(newLines, statusComment)
		}

	case StatusTodo:
		// No status comment needed for TODO
	}

	// Add all remaining lines after the comment block
	if endIndex+1 < len(lines) {
		newLines = append(newLines, lines[endIndex+1:]...)
	}

	// Write the file back
	return s.writeFileLines(fullPath, newLines)
}

// Helper method to format status comments properly
func (s *Scanner) formatStatusComment(prefix, statusText, timestamp, user string) string {
	switch prefix {
	case "<!--":
		return fmt.Sprintf("<!-- %s %s by %s -->", statusText, timestamp, user)
	default:
		return fmt.Sprintf("%s %s %s by %s", prefix, statusText, timestamp, user)
	}
}

// getCurrentUser gets the current user from git config or fallback
func (s *Scanner) getCurrentUser() string {
	if user, err := getGitUserName(); err == nil && user != "" {
		return user
	}
	// Fallback to system user
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

// Helper method to check if a line is a comment
func (s *Scanner) isCommentLine(line, prefix string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}

	switch prefix {
	case "<!--":
		return strings.HasPrefix(line, "<!--") || strings.Contains(line, "<!--")
	case "//", "#", "--":
		return strings.HasPrefix(line, prefix)
	default:
		return strings.HasPrefix(line, prefix)
	}
}

// Helper method to read file lines
func (s *Scanner) readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Helper method to write file lines
func (s *Scanner) writeFileLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := w.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

// Helper method to get comment prefix
func (s *Scanner) getCommentPrefix(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".swift", ".jsx", ".tsx", ".rs", ".kt", ".scala":
		return "//"
	case ".py", ".sh", ".rb", ".yml", ".yaml":
		return "#"
	case ".sql":
		return "--"
	case ".html", ".xml", ".vue":
		return "<!--"
	default:
		return "//"
	}
}

func getGitUserName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Additional utility methods for the Scanner

// GetItemsByType returns items filtered by type
func (s *Scanner) GetItemsByType(itemType ItemType) []*Item {
	var filtered []*Item
	for _, item := range s.Todos {
		if item.Type == itemType {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByStatus returns items filtered by status
func (s *Scanner) GetItemsByStatus(status ItemStatus) []*Item {
	var filtered []*Item
	for _, item := range s.Todos {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByPriority returns items filtered by priority
func (s *Scanner) GetItemsByPriority(priority ItemPriority) []*Item {
	var filtered []*Item
	for _, item := range s.Todos {
		if item.Priority == priority {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByCategory returns items grouped by category
func (s *Scanner) GetItemsByCategory() map[string][]*Item {
	categories := make(map[string][]*Item)
	for _, item := range s.Todos {
		category := GetTypeCategory(item.Type)
		categories[category] = append(categories[category], item)
	}
	return categories
}

// GetStats returns summary statistics
func (s *Scanner) GetStats() map[string]int {
	stats := map[string]int{
		"total":       len(s.Todos),
		"todo":        0,
		"in_progress": 0,
		"done":        0,
		"high":        0,
		"medium":      0,
		"low":         0,
	}

	for _, item := range s.Todos {
		switch item.Status {
		case StatusTodo:
			stats["todo"]++
		case StatusInProgress:
			stats["in_progress"]++
		case StatusDone:
			stats["done"]++
		}

		switch item.Priority {
		case PriorityHigh:
			stats["high"]++
		case PriorityMedium:
			stats["medium"]++
		case PriorityLow:
			stats["low"]++
		}
	}

	return stats
}

func (s *Scanner) GetTracker() *ProjectTracker {
	return s.tracker
}

func (s *Scanner) GetProjectStats() map[string]interface{} {
	return s.tracker.GetStatsSummary()
}

func (s *Scanner) GetBranchHistory() []BranchSnapshot {
	return s.tracker.GetBranchHistory()
}

func (s *Scanner) CompareWithPrevious() map[string]interface{} {
	return s.tracker.CompareWithPreviousCommit()
}

func (s *Scanner) LoadExistingStats() {
	if err := s.tracker.Initialize(); err != nil {
		fmt.Printf("Warning: Failed to initialize project tracker: %v\n", err)
		return
	}

	s.tracker.LoadStats()
}
