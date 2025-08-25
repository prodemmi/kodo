package core

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
)

const LargeFileSize = 1 * 1024 * 1024

type Scanner struct {
	Items  []*Item
	todoID int

	tracker  *ProjectTracker
	settings *SettingsManager
}

func NewScanner(config Config, settings *SettingsManager, tracker *ProjectTracker, logger *zap.Logger) *Scanner {
	scanner := &Scanner{
		tracker:  tracker,
		settings: settings,
	}
	return scanner
}

func (s *Scanner) GetItems() []*Item {
	return s.Items
}

func (s *Scanner) GetItemsLength() int {
	return len(s.Items)
}

func (s *Scanner) Rescan() error {
	s.ScanTodos()

	// Save stats after scanning
	if err := s.tracker.SaveStats(s.GetItems(), s.settings); err != nil {
		return err
	}

	return nil
}

func (s *Scanner) ScanTodos() {
	s.Items = []*Item{}
	s.todoID = 0

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	settings := s.settings.LoadSettings()
	// Create comprehensive pattern for all item types
	itemTypes := []string{}
	noneStartItemIdentifiers := []string{}

	for _, setting := range settings.KanbanColumns {
		if pt := setting.AutoAssignPattern; pt != nil {
			for _, pattern := range strings.Split(*pt, "|") {
				if pattern != "" {
					itemTypes = append(itemTypes, pattern)
				}
			}
		} else {
			noneStartItemIdentifiers = append(noneStartItemIdentifiers, setting.Name)
		}
	}

	itemPriorities := []string{}
	itemPriorities = append(itemPriorities, settings.PriorityPatterns.Low)
	itemPriorities = append(itemPriorities, settings.PriorityPatterns.Medium)
	itemPriorities = append(itemPriorities, settings.PriorityPatterns.High)

	// Build dynamic pattern
	typePattern := strings.Join(itemTypes, "|")
	priorityPatternString := strings.Join(itemPriorities, "|")
	noneStartItemIdentifiersPattern := strings.Join(noneStartItemIdentifiers, "|")
	// Match first line of a TODO item
	itemPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s):\s*(.+)?`, typePattern))
	// Match continuation description lines
	descPattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*(.+)`)
	// Match continuation description lines
	priorityPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s)`, priorityPatternString))
	// Match DONE comments
	noneStartItemPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s)\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2})\s+by\s+(.+?)(\s*-->)?$`, noneStartItemIdentifiersPattern))

	firstColumn := settings.KanbanColumns[0]

	err = filepath.Walk(wd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if slices.Contains(settings.CodeScanSettings.ExcludeDirectories, name) {
				return filepath.SkipDir
			}
			return nil
		} else {
			for _, excludeFilePattern := range settings.CodeScanSettings.ExcludeFiles {
				relPath, _ := filepath.Rel(wd, path)
				matched, err := doublestar.Match(excludeFilePattern, relPath)
				if err != nil {
					return err
				}
				if matched {
					return nil
				}
			}
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

			if matches := itemPattern.FindStringSubmatch(line); len(matches) > 0 {
				s.todoID++
				itemType := ItemType(matches[2])
				title := strings.TrimSpace(matches[3])
				todoStartLine := lineNum

				// Collect following description lines and check for status comments
				var descriptions []string
				var history []StatusHistory
				currentStatus := ItemStatus(firstColumn.ID)
				currentPriority := ItemPriority("LOW")

				// Get current user for initial creation
				currentUser := s.getCurrentUser()

				for scanner.Scan() {
					nextLine := scanner.Text()
					lineNum++

					// Check if this is a none start comment
					if noneStartMatches := noneStartItemPattern.FindStringSubmatch(nextLine); len(noneStartMatches) > 0 {
						if parsedTime, err := time.Parse("2006-01-02 15:04", noneStartMatches[3]); err == nil {
							status := ItemStatus(strcase.SnakeCase(strings.TrimSpace(noneStartMatches[2])))
							history = append(history, StatusHistory{
								Status:    status,
								Timestamp: parsedTime,
								User:      strings.TrimSpace(noneStartMatches[4]),
							})
							currentStatus = status
						}
					} else if priorityMatches := priorityPattern.FindStringSubmatch(nextLine); len(priorityMatches) > 0 {
						pr := strings.TrimSpace(priorityMatches[2])
						switch pr {
						case settings.PriorityPatterns.Low:
							currentPriority = "LOW"
						case settings.PriorityPatterns.Medium:
							currentPriority = "MEDIUM"
						case settings.PriorityPatterns.High:
							currentPriority = "HIGH"
						}
					} else if descMatches := descPattern.FindStringSubmatch(nextLine); len(descMatches) > 0 {
						desc := strings.TrimSpace(descMatches[2])
						// Skip status lines that didn't match patterns
						upperDesc := strings.ToUpper(desc)
						isStatusLine := false
						isPriorityLine := false
						for _, kanbanCol := range settings.KanbanColumns {
							if strings.HasPrefix(upperDesc, strings.ToUpper(kanbanCol.Name)) {
								isStatusLine = true
								break
							}
						}

						if desc == settings.PriorityPatterns.Low && desc == settings.PriorityPatterns.Medium && desc == settings.PriorityPatterns.High {
							isStatusLine = true
							break
						}

						if !isPriorityLine && !isStatusLine {
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
					Priority:    currentPriority,
					CreatedAt:   time.Now(), // We don't know actual creation time
					UpdatedAt:   time.Now(),
					CurrentUser: currentUser,
					History:     history,
				}

				// If no history was found, add initial creation entry
				if len(history) == 0 {
					item.History = []StatusHistory{{
						Status:    ItemStatus(firstColumn.ID),
						Timestamp: item.CreatedAt,
						User:      currentUser,
					}}
				}

				s.Items = append(s.Items, item)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

func (s *Scanner) UpdateItemStatus(item *Item, targetColumnID string) error {
	currentUser := s.getCurrentUser()
	settings := s.settings.LoadSettings()

	// Find the target Kanban column
	var targetColumn *KanbanColumn
	assignablePatterns := make(map[string]interface{}) // patterns from AutoAssignPattern
	statusColumns := make(map[string]interface{})      // columns without AutoAssignPattern

	for _, col := range settings.KanbanColumns {
		if col.ID == targetColumnID {
			targetColumn = &col
		}
		if col.AutoAssignPattern != nil {
			for _, p := range strings.Split(*col.AutoAssignPattern, "|") {
				assignablePatterns[strings.TrimSpace(p)] = struct{}{}
			}
		} else {
			statusColumns[col.Name] = struct{}{}
		}
	}

	if targetColumn == nil {
		return fmt.Errorf("kanban column with ID '%s' not found", targetColumnID)
	}

	// Update item status
	item.Status = ItemStatus(targetColumn.Name)

	// Build new comment line
	var newComment string
	if targetColumn.AutoAssignPattern == nil {
		// Status-only column → add timestamp + user
		timestamp := time.Now().Format("2006-01-02 15:04")
		newComment = fmt.Sprintf("// %s %s by %s", item.Status, timestamp, currentUser)
	} else {
		// remove line
		newComment = ""
	}

	return s.updateStatusCommentDynamic(item, newComment, assignablePatterns, statusColumns)
}

// updateStatusCommentDynamic replaces old status lines dynamically
func (s *Scanner) updateStatusCommentDynamic(item *Item, newComment string, assignablePatterns, statusColumns map[string]interface{}) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	fullPath := filepath.Join(wd, item.File)

	lines, err := s.readFileLines(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", fullPath, err)
	}

	if item.Line < 1 || item.Line > len(lines) {
		return fmt.Errorf("invalid line number %d (file has %d lines)", item.Line, len(lines))
	}

	prefix := s.getCommentPrefix(item.File)

	// Build dynamic regex for existing status lines
	var patterns []string
	for p := range assignablePatterns {
		patterns = append(patterns, regexp.QuoteMeta(p))
	}
	for s := range statusColumns {
		patterns = append(patterns, regexp.QuoteMeta(s))
	}
	statusPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|<!--)\s*(%s)(:| .*)?`, strings.Join(patterns, "|")))

	todoIndex := item.Line - 1
	endIndex := todoIndex

	// Detect end of TODO comment block
	for i := todoIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if s.isCommentLine(line, prefix) && statusPattern.MatchString(line) {
			endIndex = i
		} else if s.isCommentLine(line, prefix) {
			endIndex = i
		} else {
			break
		}
	}

	// Build new file content
	newLines := append([]string{}, lines[:todoIndex+1]...)
	for i := todoIndex + 1; i <= endIndex; i++ {
		if !statusPattern.MatchString(lines[i]) {
			newLines = append(newLines, lines[i])
		}
	}

	if strings.TrimSpace(newComment) != "" {
		newLines = append(newLines, newComment)
	}

	if endIndex+1 < len(lines) {
		newLines = append(newLines, lines[endIndex+1:]...)
	}

	return s.writeFileLines(fullPath, newLines)
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
	for _, item := range s.Items {
		if item.Type == itemType {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByStatus returns items filtered by status
func (s *Scanner) GetItemsByStatus(status ItemStatus) []*Item {
	var filtered []*Item
	for _, item := range s.Items {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByPriority returns items filtered by priority
func (s *Scanner) GetItemsByPriority(priority ItemPriority) []*Item {
	var filtered []*Item
	for _, item := range s.Items {
		if item.Priority == priority {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetItemsByCategory returns items grouped by category
func (s *Scanner) GetItemsByCategory() map[string][]*Item {
	categories := make(map[string][]*Item)
	for _, item := range s.Items {
		category := string(item.Type)
		categories[category] = append(categories[category], item)
	}
	return categories
}

func (s *Scanner) UpdateOldStatuses(oldSettings, newSettings *Settings) error {
	// Build map of old ID -> old Name
	oldIds := make(map[int]string)
	for i, col := range oldSettings.KanbanColumns {
		oldIds[i] = col.ID
	}

	// Build map of new ID -> new Name
	newIds := make(map[int]string)
	for i, col := range newSettings.KanbanColumns {
		newIds[i] = col.ID
	}

	// Now check which IDs exist in both but with different names
	renamed := make(map[string]string) // oldName -> newName
	for id, oldId := range oldIds {
		if newId, ok := newIds[id]; ok && oldId != newId {
			renamed[oldId] = newId
		}
	}

	fmt.Println("oldIds", oldIds)
	fmt.Println("newIds", newIds)
	fmt.Println("renamed", renamed)

	if len(renamed) == 0 {
		return nil // nothing to update
	}

	// Walk over all items and update their status if it matches an old name
	for _, item := range s.Items {
		if newName, ok := renamed[string(item.Status)]; ok {
			item.Status = ItemStatus(newName)

			// also patch history if needed
			for i, h := range item.History {
				if string(h.Status) == renamed[string(item.Status)] {
					item.History[i].Status = item.Status
				}
			}
		}
	}

	return s.tracker.SaveStats(s.GetItems(), s.settings)
}
