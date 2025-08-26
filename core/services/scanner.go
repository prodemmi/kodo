package services

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
	"github.com/prodemmi/kodo/core/entities"
	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
)

const LargeFileSize = 1 * 1024 * 1024

type ScannerService struct {
	Items  []*entities.Item
	todoID int

	itemHistoryService *HistoryService
	settings           *SettingsService
}

func NewScannerService(config *entities.Config, settings *SettingsService, itemHistoryService *HistoryService, logger *zap.Logger) *ScannerService {
	scannerService := &ScannerService{
		itemHistoryService: itemHistoryService,
		settings:           settings,
	}
	return scannerService
}

func (s *ScannerService) GetItems() []*entities.Item {
	return s.Items
}

func (s *ScannerService) GetItemsLength() int {
	return len(s.Items)
}

func (s *ScannerService) Rescan() error {
	s.ScanTodos()

	if err := s.itemHistoryService.SaveStats(s.GetItems(), s.settings); err != nil {
		return err
	}

	return nil
}

func (s *ScannerService) ScanTodos() {
	s.Items = []*entities.Item{}
	s.todoID = 0

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	settings := s.settings.LoadSettings()

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

	typePattern := strings.Join(itemTypes, "|")
	priorityPatternString := strings.Join(itemPriorities, "|")
	noneStartItemIdentifiersPattern := strings.Join(noneStartItemIdentifiers, "|")

	itemPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s):\s*(.+)?`, typePattern))

	descPattern := regexp.MustCompile(`^\s*(//|#|--|\<!--)\s*(.+)`)

	priorityPattern := regexp.MustCompile(fmt.Sprintf(`^\s*(//|#|--|\<!--)\s*(%s)`, priorityPatternString))

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

		scannerService := bufio.NewScanner(file)
		lineNum := 0
		for scannerService.Scan() {
			lineNum++
			line := scannerService.Text()

			if matches := itemPattern.FindStringSubmatch(line); len(matches) > 0 {
				s.todoID++
				itemType := entities.ItemType(matches[2])
				title := strings.TrimSpace(matches[3])
				todoStartLine := lineNum

				var descriptions []string
				var history []entities.StatusHistory
				currentStatus := entities.ItemStatus(firstColumn.ID)
				currentPriority := entities.ItemPriority("LOW")

				currentUser := s.getCurrentUser()

				for scannerService.Scan() {
					nextLine := scannerService.Text()
					lineNum++

					if noneStartMatches := noneStartItemPattern.FindStringSubmatch(nextLine); len(noneStartMatches) > 0 {
						if parsedTime, err := time.Parse("2006-01-02 15:04", noneStartMatches[3]); err == nil {
							status := entities.ItemStatus(strcase.SnakeCase(strings.TrimSpace(noneStartMatches[2])))
							history = append(history, entities.StatusHistory{
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
						break
					}
				}

				relPath, _ := filepath.Rel(wd, path)

				item := &entities.Item{
					ID:          s.todoID,
					Type:        itemType,
					Title:       title,
					Description: strings.Join(descriptions, "\n"),
					File:        relPath,
					Line:        todoStartLine,
					Status:      currentStatus,
					Priority:    currentPriority,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					CurrentUser: currentUser,
					History:     history,
				}

				if len(history) == 0 {
					item.History = []entities.StatusHistory{{
						Status:    entities.ItemStatus(firstColumn.ID),
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

func (s *ScannerService) UpdateItemStatus(item *entities.Item, targetColumnID string) error {
	currentUser := s.getCurrentUser()
	settings := s.settings.LoadSettings()

	var targetColumn *entities.KanbanColumn
	assignablePatterns := make(map[string]interface{})
	statusColumns := make(map[string]interface{})

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

	item.Status = entities.ItemStatus(targetColumn.Name)

	var newComment string
	if targetColumn.AutoAssignPattern == nil {
		timestamp := time.Now().Format("2006-01-02 15:04")
		newComment = fmt.Sprintf("// %s %s by %s", item.Status, timestamp, currentUser)
	} else {
		newComment = ""
	}

	return s.updateStatusCommentDynamic(item, newComment, assignablePatterns, statusColumns)
}

func (s *ScannerService) updateStatusCommentDynamic(item *entities.Item, newComment string, assignablePatterns, statusColumns map[string]interface{}) error {
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

func (s *ScannerService) getCurrentUser() string {
	if user, err := getGitUserName(); err == nil && user != "" {
		return user
	}

	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

func (s *ScannerService) isCommentLine(line, prefix string) bool {
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

func (s *ScannerService) readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scannerService := bufio.NewScanner(file)
	for scannerService.Scan() {
		lines = append(lines, scannerService.Text())
	}
	return lines, scannerService.Err()
}

func (s *ScannerService) writeFileLines(filename string, lines []string) error {
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

func (s *ScannerService) getCommentPrefix(filename string) string {
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

func (s *ScannerService) GetItemsByType(itemType entities.ItemType) []*entities.Item {
	var filtered []*entities.Item
	for _, item := range s.Items {
		if item.Type == itemType {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *ScannerService) GetItemsByStatus(status entities.ItemStatus) []*entities.Item {
	var filtered []*entities.Item
	for _, item := range s.Items {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *ScannerService) GetItemsByPriority(priority entities.ItemPriority) []*entities.Item {
	var filtered []*entities.Item
	for _, item := range s.Items {
		if item.Priority == priority {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *ScannerService) GetItemsByCategory() map[string][]*entities.Item {
	categories := make(map[string][]*entities.Item)
	for _, item := range s.Items {
		category := string(item.Type)
		categories[category] = append(categories[category], item)
	}
	return categories
}

func (s *ScannerService) UpdateOldStatuses(oldSettings, newSettings *entities.Settings) error {

	oldIds := make(map[int]string)
	for i, col := range oldSettings.KanbanColumns {
		oldIds[i] = col.ID
	}

	newIds := make(map[int]string)
	for i, col := range newSettings.KanbanColumns {
		newIds[i] = col.ID
	}

	renamed := make(map[string]string)
	for id, oldId := range oldIds {
		if newId, ok := newIds[id]; ok && oldId != newId {
			renamed[oldId] = newId
		}
	}

	fmt.Println("oldIds", oldIds)
	fmt.Println("newIds", newIds)
	fmt.Println("renamed", renamed)

	if len(renamed) == 0 {
		return nil
	}

	for _, item := range s.Items {
		if newName, ok := renamed[string(item.Status)]; ok {
			item.Status = entities.ItemStatus(newName)

			for i, h := range item.History {
				if string(h.Status) == renamed[string(item.Status)] {
					item.History[i].Status = item.Status
				}
			}
		}
	}

	return s.itemHistoryService.SaveStats(s.GetItems(), s.settings)
}
