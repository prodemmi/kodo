package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prodemmi/kodo/core/entities"
	"go.uber.org/zap"
)

type NoteService struct {
	Notes   []entities.Note   `json:"notes"`
	Folders []entities.Folder `json:"folders"`
	NextID  int               `json:"next_id"`

	logger *zap.Logger
	config *entities.Config
}

func NewNoteService(config *entities.Config, logger *zap.Logger) *NoteService {
	return &NoteService{
		Notes:   []entities.Note{},
		Folders: []entities.Folder{},
		NextID:  1,
		logger:  logger,
		config:  config,
	}
}

func (s *NoteService) ensureNotesDir() error {
	wd, _ := os.Getwd()
	notesDir := filepath.Join(wd, s.config.Flags.Config)

	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %v", err)
	}

	return nil
}

func (s *NoteService) getNotesFilePath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, s.config.Flags.Config, "notes.json")
}

func (s *NoteService) loadNoteStorage() (*NoteService, error) {
	if err := s.ensureNotesDir(); err != nil {
		return nil, err
	}

	filePath := s.getNotesFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage := &NoteService{
			Notes:   []entities.Note{},
			Folders: []entities.Folder{},
			NextID:  1,
		}
		return storage, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %v", err)
	}

	var storage NoteService
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notes: %v", err)
	}

	return &storage, nil
}

func (s *NoteService) SaveNoteStorage(storage *NoteService) error {
	if err := s.ensureNotesDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %v", err)
	}

	filePath := s.getNotesFilePath()
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %v", err)
	}

	return nil
}

func (s *NoteService) getGitAuthor() string {
	cmd := exec.Command("git", "config", "--get", "user.name")
	nameOutput, nameErr := cmd.Output()

	cmd = exec.Command("git", "config", "--get", "user.email")
	emailOutput, emailErr := cmd.Output()

	var author string
	if nameErr == nil && len(nameOutput) > 0 {
		author = strings.TrimSpace(string(nameOutput))
	}

	if emailErr == nil && len(emailOutput) > 0 {
		email := strings.TrimSpace(string(emailOutput))
		if author != "" {
			author = fmt.Sprintf("%s <%s>", author, email)
		} else {
			author = email
		}
	}

	if author == "" {
		author = "Unknown"
	}

	return author
}

func (s *NoteService) getGitInfo() (string, string) {

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, branchErr := cmd.Output()
	branch := "unknown"
	if branchErr == nil {
		branch = strings.TrimSpace(string(branchOutput))
	}

	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	commitOutput, commitErr := cmd.Output()
	commit := "unknown"
	if commitErr == nil {
		commit = strings.TrimSpace(string(commitOutput))
	}

	return branch, commit
}

func (s *NoteService) GetNotes(category, tag, folderId string) ([]entities.Note, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	var filteredNotes []entities.Note

	for _, note := range storage.Notes {

		if category != "" && note.Category != category {
			continue
		}

		if tag != "" {
			hasTag := false
			for _, noteTag := range note.Tags {
				if noteTag == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		if folderId != "" {
			folderIdInt, err := strconv.Atoi(folderId)
			if err != nil {
				continue
			}
			if note.FolderID == nil || *note.FolderID != folderIdInt {
				continue
			}
		}

		filteredNotes = append(filteredNotes, note)
	}

	sort.Slice(filteredNotes, func(i, j int) bool {
		return filteredNotes[i].UpdatedAt.After(filteredNotes[j].UpdatedAt)
	})

	return filteredNotes, nil
}

func (s *NoteService) CreateFolder(name string, parentId *int) (*entities.Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	folder := entities.Folder{
		ID:       storage.NextID,
		Name:     name,
		ParentID: parentId,
		Expanded: true,
	}

	storage.Folders = append(storage.Folders, folder)
	storage.NextID++

	if err := s.SaveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &folder, nil
}

func (s *NoteService) GetFolders() ([]entities.Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	sort.Slice(storage.Folders, func(i, j int) bool {
		return storage.Folders[i].Name < storage.Folders[j].Name
	})

	return storage.Folders, nil
}

func (s *NoteService) GetCategories() []entities.Category {
	return []entities.Category{
		{Value: "general", Label: "General"},
		{Value: "development", Label: "Development"},
		{Value: "documentation", Label: "Documentation"},
		{Value: "meeting", Label: "Meeting"},
		{Value: "idea", Label: "Idea"},
		{Value: "bug", Label: "Bug"},
		{Value: "feature", Label: "Feature"},
		{Value: "review", Label: "Review"},
		{Value: "research", Label: "Research"},
		{Value: "personal", Label: "Personal"},
	}
}

func (s *NoteService) GetNoteStats() (map[string]interface{}, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	history := map[string]interface{}{
		"total_notes":   len(storage.Notes),
		"total_folders": len(storage.Folders),
		"by_category":   make(map[string]int),
		"by_author":     make(map[string]int),
		"by_month":      make(map[string]int),
		"recent_notes":  []entities.Note{},
	}

	categoryCount := make(map[string]int)
	authorCount := make(map[string]int)
	monthCount := make(map[string]int)

	var recentNotes []entities.Note
	cutoff := time.Now().AddDate(0, 0, -7)

	for _, note := range storage.Notes {

		if note.Category != "" {
			categoryCount[note.Category]++
		}

		if note.Author != "" {
			authorCount[note.Author]++
		}

		monthKey := note.CreatedAt.Format("2006-01")
		monthCount[monthKey]++

		if note.CreatedAt.After(cutoff) {
			recentNotes = append(recentNotes, note)
		}
	}

	sort.Slice(recentNotes, func(i, j int) bool {
		return recentNotes[i].CreatedAt.After(recentNotes[j].CreatedAt)
	})

	if len(recentNotes) > 10 {
		recentNotes = recentNotes[:10]
	}

	history["by_category"] = categoryCount
	history["by_author"] = authorCount
	history["by_month"] = monthCount
	history["recent_notes"] = recentNotes

	return history, nil
}

func (s *NoteService) deleteFolderRecursive(storage *NoteService, id int) error {

	folderIndex := -1
	for i, folder := range storage.Folders {
		if folder.ID == id {
			folderIndex = i
			break
		}
	}
	if folderIndex == -1 {
		return errors.New("folder not found")
	}

	newNotes := storage.Notes[:0]
	for _, note := range storage.Notes {
		if note.FolderID == nil || *note.FolderID != id {
			newNotes = append(newNotes, note)
		}
	}
	storage.Notes = newNotes

	for _, folder := range storage.Folders {
		if folder.ParentID != nil && *folder.ParentID == id {
			if err := s.deleteFolderRecursive(storage, folder.ID); err != nil {
				return err
			}
		}
	}

	storage.Folders = append(storage.Folders[:folderIndex], storage.Folders[folderIndex+1:]...)

	return nil
}

func (s *NoteService) UpdateFolder(id int, name string, parentId *int, expanded *bool) (*entities.Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	folderIndex := -1
	for i, folder := range storage.Folders {
		if folder.ID == id {
			folderIndex = i
			break
		}
	}

	if folderIndex == -1 {
		return nil, errors.New("folder not found")
	}

	if name != "" {
		storage.Folders[folderIndex].Name = name
	}
	if expanded != nil {
		storage.Folders[folderIndex].Expanded = *expanded
	}
	storage.Folders[folderIndex].ParentID = parentId

	if err := s.SaveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &storage.Folders[folderIndex], nil
}

func (s *NoteService) DeleteFolder(id int) error {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return err
	}

	if err := s.deleteFolderRecursive(storage, id); err != nil {
		return err
	}

	return s.SaveNoteStorage(storage)
}

func (s *NoteService) GetFolderTree() ([]map[string]interface{}, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	folderMap := make(map[int]*entities.Folder)
	for i := range storage.Folders {
		folderMap[storage.Folders[i].ID] = &storage.Folders[i]
	}

	var rootFolders []map[string]interface{}

	for _, folder := range storage.Folders {
		if folder.ParentID == nil {

			rootFolders = append(rootFolders, s.buildFolderNode(folder, folderMap, storage.Notes))
		}
	}

	sort.Slice(rootFolders, func(i, j int) bool {
		return rootFolders[i]["name"].(string) < rootFolders[j]["name"].(string)
	})

	return rootFolders, nil
}

func (s *NoteService) buildFolderNode(folder entities.Folder, folderMap map[int]*entities.Folder, notes []entities.Note) map[string]interface{} {

	noteCount := 0
	for _, note := range notes {
		if note.FolderID != nil && *note.FolderID == folder.ID {
			noteCount++
		}
	}

	var children []map[string]interface{}
	for _, otherFolder := range folderMap {
		if otherFolder.ParentID != nil && *otherFolder.ParentID == folder.ID {
			children = append(children, s.buildFolderNode(*otherFolder, folderMap, notes))
		}
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i]["name"].(string) < children[j]["name"].(string)
	})

	totalNoteCount := noteCount
	for _, child := range children {
		totalNoteCount += child["totalNoteCount"].(int)
	}

	return map[string]interface{}{
		"id":             folder.ID,
		"name":           folder.Name,
		"parentId":       folder.ParentID,
		"expanded":       folder.Expanded,
		"noteCount":      noteCount,
		"totalNoteCount": totalNoteCount,
		"children":       children,
		"hasChildren":    len(children) > 0,
	}
}

func (s *NoteService) SearchNotes(query string, category string, folderId *int) ([]entities.Note, error) {
	if query == "" {
		return s.GetNotes(category, "", "")
	}

	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var matchingNotes []entities.Note

	for _, note := range storage.Notes {

		if category != "" && note.Category != category {
			continue
		}

		if folderId != nil {
			if note.FolderID == nil || *note.FolderID != *folderId {
				continue
			}
		}

		if strings.Contains(strings.ToLower(note.Title), query) {
			matchingNotes = append(matchingNotes, note)
			continue
		}

		if strings.Contains(strings.ToLower(note.Content), query) {
			matchingNotes = append(matchingNotes, note)
			continue
		}

		for _, tag := range note.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matchingNotes = append(matchingNotes, note)
				break
			}
		}
	}

	sort.Slice(matchingNotes, func(i, j int) bool {
		titleMatchI := strings.Contains(strings.ToLower(matchingNotes[i].Title), query)
		titleMatchJ := strings.Contains(strings.ToLower(matchingNotes[j].Title), query)

		if titleMatchI && !titleMatchJ {
			return true
		}
		if !titleMatchI && titleMatchJ {
			return false
		}

		return matchingNotes[i].UpdatedAt.After(matchingNotes[j].UpdatedAt)
	})

	return matchingNotes, nil
}

func (s *NoteService) ExportNotes(folderId *int, category string) (map[string]interface{}, error) {
	var notes []entities.Note
	var err error

	if folderId != nil || category != "" {
		folderIdStr := ""
		if folderId != nil {
			folderIdStr = strconv.Itoa(*folderId)
		}
		notes, err = s.GetNotes(category, "", folderIdStr)
	} else {
		notes, err = s.GetNotes("", "", "")
	}

	if err != nil {
		return nil, err
	}

	folders, err := s.GetFolders()
	if err != nil {
		return nil, err
	}

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"notes":       notes,
		"folders":     folders,
		"categories":  s.GetCategories(),
		"total":       len(notes),
	}

	return export, nil
}

func (s *NoteService) GetAllTags() ([]string, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, note := range storage.Notes {
		for _, tag := range note.Tags {
			if tag != "" {
				tagSet[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	sort.Strings(tags)

	return tags, nil
}

func (s *NoteService) LoadEnhancedNoteStorage() (*entities.EnhancedNoteStorage, error) {
	if err := s.ensureNotesDir(); err != nil {
		return nil, err
	}

	filePath := s.getNotesFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage := &entities.EnhancedNoteStorage{
			Notes:         []entities.Note{},
			Folders:       []entities.Folder{},
			History:       []entities.NoteHistoryEntry{},
			NextID:        1,
			NextHistoryID: 1,
		}
		return storage, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %v", err)
	}

	var enhancedStorage entities.EnhancedNoteStorage
	if err := json.Unmarshal(data, &enhancedStorage); err == nil && enhancedStorage.NextHistoryID > 0 {
		return &enhancedStorage, nil
	}

	var oldStorage NoteService
	if err := json.Unmarshal(data, &oldStorage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notes: %v", err)
	}

	enhancedStorage = entities.EnhancedNoteStorage{
		Notes:         oldStorage.Notes,
		Folders:       oldStorage.Folders,
		History:       []entities.NoteHistoryEntry{},
		NextID:        oldStorage.NextID,
		NextHistoryID: 1,
	}

	return &enhancedStorage, nil
}

func (s *NoteService) SaveEnhancedNoteStorage(storage *entities.EnhancedNoteStorage) error {
	if err := s.ensureNotesDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %v", err)
	}

	filePath := s.getNotesFilePath()
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %v", err)
	}

	return nil
}

func (s *NoteService) addNoteHistoryEntry(storage *entities.EnhancedNoteStorage, noteID int, action entities.NoteHistoryAction, changes map[string]interface{}, oldValue, newValue interface{}, message string) {
	author := s.getGitAuthor()
	branch, commit := s.getGitInfo()

	entry := entities.NoteHistoryEntry{
		ID:        storage.NextHistoryID,
		NoteID:    noteID,
		Action:    action,
		Author:    author,
		Timestamp: time.Now(),
		GitBranch: &branch,
		GitCommit: &commit,
		Changes:   changes,
		OldValue:  oldValue,
		NewValue:  newValue,
		Message:   message,
		Metadata:  make(map[string]interface{}),
	}

	storage.History = append(storage.History, entry)
	storage.NextHistoryID++
}

func (s *NoteService) CreateNoteWithHistory(title, content string, tags []string, category string, folderId *int, author *string) (*entities.Note, error) {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if author == nil {
		gitUser := s.getGitAuthor()
		author = &gitUser
	}
	branch, commit := s.getGitInfo()

	status := "open"
	note := entities.Note{
		ID:        storage.NextID,
		Title:     title,
		Content:   content,
		Author:    *author,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      tags,
		Category:  category,
		FolderID:  folderId,
		GitBranch: &branch,
		GitCommit: &commit,
		Pinned:    false,
		Status:    &status,
	}

	storage.Notes = append(storage.Notes, note)
	storage.NextID++

	changes := map[string]interface{}{
		"title":     title,
		"category":  category,
		"tags":      tags,
		"folder_id": folderId,
	}
	s.addNoteHistoryEntry(storage, note.ID, entities.ActionCreated, changes, nil, note, "Note created")

	if err := s.SaveEnhancedNoteStorage(storage); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *NoteService) UpdateNoteWithHistory(id int, title, content string, tags []string, category string, pinned bool, folderId *int) (*entities.Note, error) {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	noteIndex := -1
	for i, note := range storage.Notes {
		if note.ID == id {
			noteIndex = i
			break
		}
	}

	if noteIndex == -1 {
		return nil, fmt.Errorf("note not found")
	}

	oldNote := storage.Notes[noteIndex]

	changes := make(map[string]interface{})
	if oldNote.Title != title {
		changes["title"] = map[string]interface{}{"from": oldNote.Title, "to": title}
	}
	if oldNote.Content != content {
		changes["content"] = map[string]interface{}{"from": len(oldNote.Content), "to": len(content)}
	}
	if oldNote.Category != category {
		changes["category"] = map[string]interface{}{"from": oldNote.Category, "to": category}
	}
	if oldNote.Pinned != pinned {
		changes["pinned"] = map[string]interface{}{"from": oldNote.Pinned, "to": pinned}
	}
	if !equalSlices(oldNote.Tags, tags) {
		changes["tags"] = map[string]interface{}{"from": oldNote.Tags, "to": tags}
	}
	if !equalFolderID(oldNote.FolderID, folderId) {
		changes["folder_id"] = map[string]interface{}{"from": oldNote.FolderID, "to": folderId}
	}

	branch, commit := s.getGitInfo()
	storage.Notes[noteIndex].Title = title
	storage.Notes[noteIndex].Content = content
	storage.Notes[noteIndex].Tags = tags
	storage.Notes[noteIndex].Category = category
	storage.Notes[noteIndex].Pinned = pinned
	storage.Notes[noteIndex].FolderID = folderId
	storage.Notes[noteIndex].UpdatedAt = time.Now()
	storage.Notes[noteIndex].GitBranch = &branch
	storage.Notes[noteIndex].GitCommit = &commit

	if len(changes) > 0 {
		s.addNoteHistoryEntry(storage, id, entities.ActionUpdated, changes, oldNote, storage.Notes[noteIndex], "Note updated")
	}

	if err := s.SaveEnhancedNoteStorage(storage); err != nil {
		return nil, err
	}

	return &storage.Notes[noteIndex], nil
}

func (s *NoteService) DeleteNoteWithHistory(id int) error {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return err
	}

	noteIndex := -1
	var deletedNote entities.Note
	for i, note := range storage.Notes {
		if note.ID == id {
			noteIndex = i
			deletedNote = note
			break
		}
	}

	if noteIndex == -1 {
		return fmt.Errorf("note not found")
	}

	changes := map[string]interface{}{
		"deleted_note": map[string]interface{}{
			"title":    deletedNote.Title,
			"category": deletedNote.Category,
			"tags":     deletedNote.Tags,
		},
	}
	s.addNoteHistoryEntry(storage, id, entities.ActionDeleted, changes, deletedNote, nil, "Note deleted")

	storage.Notes = append(storage.Notes[:noteIndex], storage.Notes[noteIndex+1:]...)

	return s.SaveEnhancedNoteStorage(storage)
}

func (s *NoteService) MoveNotesToFolderWithHistory(noteIds []int, targetFolderId *int) error {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return err
	}

	if targetFolderId != nil {
		folderExists := false
		for _, folder := range storage.Folders {
			if folder.ID == *targetFolderId {
				folderExists = true
				break
			}
		}
		if !folderExists {
			return fmt.Errorf("target folder not found")
		}
	}

	updatedCount := 0
	branch, commit := s.getGitInfo()

	for i := range storage.Notes {
		for _, noteId := range noteIds {
			if storage.Notes[i].ID == noteId {
				oldFolderId := storage.Notes[i].FolderID

				storage.Notes[i].FolderID = targetFolderId
				storage.Notes[i].UpdatedAt = time.Now()
				storage.Notes[i].GitBranch = &branch
				storage.Notes[i].GitCommit = &commit

				changes := map[string]interface{}{
					"folder_id": map[string]interface{}{
						"from": oldFolderId,
						"to":   targetFolderId,
					},
				}
				s.addNoteHistoryEntry(storage, noteId, entities.ActionMoved, changes, oldFolderId, targetFolderId, "Note moved to different folder")

				updatedCount++
				break
			}
		}
	}

	if updatedCount == 0 {
		return fmt.Errorf("no notes found to move")
	}

	return s.SaveEnhancedNoteStorage(storage)
}

func (s *NoteService) GetNoteHistory(filter entities.NoteHistoryFilter) ([]entities.NoteHistoryEntry, error) {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	var filteredHistory []entities.NoteHistoryEntry

	for _, entry := range storage.History {

		if filter.NoteID != nil && entry.NoteID != *filter.NoteID {
			continue
		}
		if filter.Action != nil && entry.Action != *filter.Action {
			continue
		}
		if filter.Author != nil && entry.Author != *filter.Author {
			continue
		}
		if filter.GitBranch != nil && (entry.GitBranch == nil || *entry.GitBranch != *filter.GitBranch) {
			continue
		}
		if filter.Since != nil && entry.Timestamp.Before(*filter.Since) {
			continue
		}
		if filter.Until != nil && entry.Timestamp.After(*filter.Until) {
			continue
		}

		filteredHistory = append(filteredHistory, entry)
	}

	sort.Slice(filteredHistory, func(i, j int) bool {
		return filteredHistory[i].Timestamp.After(filteredHistory[j].Timestamp)
	})

	if filter.Offset > 0 {
		if filter.Offset >= len(filteredHistory) {
			return []entities.NoteHistoryEntry{}, nil
		}
		filteredHistory = filteredHistory[filter.Offset:]
	}

	if filter.Limit > 0 && len(filteredHistory) > filter.Limit {
		filteredHistory = filteredHistory[:filter.Limit]
	}

	return filteredHistory, nil
}

func (s *NoteService) GetNoteHistoryStats() (*entities.NoteHistoryStats, error) {
	storage, err := s.LoadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	history := &entities.NoteHistoryStats{
		TotalEntries: len(storage.History),
		ByAction:     make(map[entities.NoteHistoryAction]int),
		ByAuthor:     make(map[string]int),
		ByBranch:     make(map[string]int),
		ByDay:        make(map[string]int),
	}

	noteActivity := make(map[int]*entities.NoteActivitySummary)

	for _, entry := range storage.History {

		history.ByAction[entry.Action]++

		history.ByAuthor[entry.Author]++

		if entry.GitBranch != nil {
			history.ByBranch[*entry.GitBranch]++
		}

		day := entry.Timestamp.Format("2006-01-02")
		history.ByDay[day]++

		if activity, exists := noteActivity[entry.NoteID]; exists {
			activity.ActionCount++
			if entry.Timestamp.After(activity.LastAction) {
				activity.LastAction = entry.Timestamp
			}

			found := false
			for _, author := range activity.Authors {
				if author == entry.Author {
					found = true
					break
				}
			}
			if !found {
				activity.Authors = append(activity.Authors, entry.Author)
			}
		} else {

			noteTitle := "Unknown"
			for _, note := range storage.Notes {
				if note.ID == entry.NoteID {
					noteTitle = note.Title
					break
				}
			}

			noteActivity[entry.NoteID] = &entities.NoteActivitySummary{
				NoteID:      entry.NoteID,
				NoteTitle:   noteTitle,
				ActionCount: 1,
				LastAction:  entry.Timestamp,
				Authors:     []string{entry.Author},
			}
		}
	}

	for _, activity := range noteActivity {
		history.MostActiveNotes = append(history.MostActiveNotes, *activity)
	}
	sort.Slice(history.MostActiveNotes, func(i, j int) bool {
		return history.MostActiveNotes[i].ActionCount > history.MostActiveNotes[j].ActionCount
	})

	if len(history.MostActiveNotes) > 10 {
		history.MostActiveNotes = history.MostActiveNotes[:10]
	}

	recentCount := 20
	if len(storage.History) < recentCount {
		recentCount = len(storage.History)
	}

	sortedHistory := make([]entities.NoteHistoryEntry, len(storage.History))
	copy(sortedHistory, storage.History)
	sort.Slice(sortedHistory, func(i, j int) bool {
		return sortedHistory[i].Timestamp.After(sortedHistory[j].Timestamp)
	})

	history.RecentActivity = sortedHistory[:recentCount]

	return history, nil
}

func (s *NoteService) getNoteByID(noteID int) *entities.Note {
	notes, _ := s.GetNotes("", "", "")
	for _, note := range notes {
		if note.ID == noteID {
			return &note
		}
	}

	return nil
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func equalFolderID(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
