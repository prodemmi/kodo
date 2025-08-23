package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Note storage structure
type NoteStorage struct {
	Notes   []Note   `json:"notes"`
	Folders []Folder `json:"folders"`
	NextID  int      `json:"next_id"`
}

// ensureNotesDir creates the notes directory if it doesn't exist
func (s *Server) ensureNotesDir() error {
	wd, _ := os.Getwd()
	notesDir := filepath.Join(wd, s.config.Flags.Config)

	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %v", err)
	}

	return nil
}

// getNotesFilePath returns the path to the notes storage file
func (s *Server) getNotesFilePath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, s.config.Flags.Config, "notes.json")
}

// loadNoteStorage loads the note storage from file
func (s *Server) loadNoteStorage() (*NoteStorage, error) {
	if err := s.ensureNotesDir(); err != nil {
		return nil, err
	}

	filePath := s.getNotesFilePath()

	// If file doesn't exist, create empty storage
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage := &NoteStorage{
			Notes:   []Note{},
			Folders: []Folder{},
			NextID:  1,
		}
		return storage, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %v", err)
	}

	var storage NoteStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notes: %v", err)
	}

	return &storage, nil
}

// saveNoteStorage saves the note storage to file
func (s *Server) saveNoteStorage(storage *NoteStorage) error {
	if err := s.ensureNotesDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %v", err)
	}

	filePath := s.getNotesFilePath()
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %v", err)
	}

	return nil
}

// getGitAuthor gets the git author name and email
func (s *Server) getGitAuthor() string {
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

// getGitInfo gets current git branch and commit
func (s *Server) getGitInfo() (string, string) {
	// Get branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, branchErr := cmd.Output()
	branch := "unknown"
	if branchErr == nil {
		branch = strings.TrimSpace(string(branchOutput))
	}

	// Get commit
	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	commitOutput, commitErr := cmd.Output()
	commit := "unknown"
	if commitErr == nil {
		commit = strings.TrimSpace(string(commitOutput))
	}

	return branch, commit
}

// createNote creates a new note
func (s *Server) createNote(title, content string, tags []string, category string, folderId *int) (*Note, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	author := s.getGitAuthor()
	branch, commit := s.getGitInfo()

	note := Note{
		ID:        storage.NextID,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      tags,
		Category:  category,
		FolderID:  folderId,
		GitBranch: &branch,
		GitCommit: &commit,
	}

	storage.Notes = append(storage.Notes, note)
	storage.NextID++

	if err := s.saveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &note, nil
}

// getNotes retrieves notes with optional filtering
func (s *Server) getNotes(category, tag, folderId string) ([]Note, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	var filteredNotes []Note

	for _, note := range storage.Notes {
		// Filter by category
		if category != "" && note.Category != category {
			continue
		}

		// Filter by tag
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

		// Filter by folder ID
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

	// Sort by UpdatedAt descending (most recent first)
	sort.Slice(filteredNotes, func(i, j int) bool {
		return filteredNotes[i].UpdatedAt.After(filteredNotes[j].UpdatedAt)
	})

	return filteredNotes, nil
}

// updateNote updates an existing note
func (s *Server) updateNote(id int, title, content string, tags []string, category string, folderId *int) (*Note, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	// Find the note
	noteIndex := -1
	for i, note := range storage.Notes {
		if note.ID == id {
			noteIndex = i
			break
		}
	}

	if noteIndex == -1 {
		return nil, errors.New("note not found")
	}

	// Update the note
	branch, commit := s.getGitInfo()
	storage.Notes[noteIndex].Title = title
	storage.Notes[noteIndex].Content = content
	storage.Notes[noteIndex].Tags = tags
	storage.Notes[noteIndex].Category = category
	storage.Notes[noteIndex].FolderID = folderId
	storage.Notes[noteIndex].UpdatedAt = time.Now()
	storage.Notes[noteIndex].GitBranch = &branch
	storage.Notes[noteIndex].GitCommit = &commit

	if err := s.saveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &storage.Notes[noteIndex], nil
}

// deleteNote deletes a note by ID
func (s *Server) deleteNote(id int) error {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return err
	}

	// Find and remove the note
	noteIndex := -1
	for i, note := range storage.Notes {
		if note.ID == id {
			noteIndex = i
			break
		}
	}

	if noteIndex == -1 {
		return errors.New("note not found")
	}

	// Remove the note from slice
	storage.Notes = append(storage.Notes[:noteIndex], storage.Notes[noteIndex+1:]...)

	return s.saveNoteStorage(storage)
}

// createFolder creates a new folder
func (s *Server) createFolder(name string, parentId *int) (*Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	folder := Folder{
		ID:       storage.NextID,
		Name:     name,
		ParentID: parentId,
		Expanded: true, // New folders are expanded by default
	}

	storage.Folders = append(storage.Folders, folder)
	storage.NextID++

	if err := s.saveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &folder, nil
}

// getFolders retrieves all folders
func (s *Server) getFolders() ([]Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	// Sort folders by name
	sort.Slice(storage.Folders, func(i, j int) bool {
		return storage.Folders[i].Name < storage.Folders[j].Name
	})

	return storage.Folders, nil
}

// getCategories returns predefined categories
func (s *Server) getCategories() []Category {
	return []Category{
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

// getNoteStats provides statistics about notes
func (s *Server) getNoteStats() (map[string]interface{}, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_notes":   len(storage.Notes),
		"total_folders": len(storage.Folders),
		"by_category":   make(map[string]int),
		"by_author":     make(map[string]int),
		"by_month":      make(map[string]int),
		"recent_notes":  []Note{},
	}

	// Calculate statistics
	categoryCount := make(map[string]int)
	authorCount := make(map[string]int)
	monthCount := make(map[string]int)

	var recentNotes []Note
	cutoff := time.Now().AddDate(0, 0, -7) // Last 7 days

	for _, note := range storage.Notes {
		// Count by category
		if note.Category != "" {
			categoryCount[note.Category]++
		}

		// Count by author
		if note.Author != "" {
			authorCount[note.Author]++
		}

		// Count by month
		monthKey := note.CreatedAt.Format("2006-01")
		monthCount[monthKey]++

		// Collect recent notes
		if note.CreatedAt.After(cutoff) {
			recentNotes = append(recentNotes, note)
		}
	}

	// Sort recent notes by creation time
	sort.Slice(recentNotes, func(i, j int) bool {
		return recentNotes[i].CreatedAt.After(recentNotes[j].CreatedAt)
	})

	// Limit to 10 recent notes
	if len(recentNotes) > 10 {
		recentNotes = recentNotes[:10]
	}

	stats["by_category"] = categoryCount
	stats["by_author"] = authorCount
	stats["by_month"] = monthCount
	stats["recent_notes"] = recentNotes

	return stats, nil
}

// deleteFolderRecursive works in-memory (no reload/save inside)
func (s *Server) deleteFolderRecursive(storage *NoteStorage, id int) error {
	// Find the folder
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

	// Delete notes in this folder
	newNotes := storage.Notes[:0]
	for _, note := range storage.Notes {
		if note.FolderID == nil || *note.FolderID != id {
			newNotes = append(newNotes, note)
		}
	}
	storage.Notes = newNotes

	// Delete subfolders recursively
	for _, folder := range storage.Folders {
		if folder.ParentID != nil && *folder.ParentID == id {
			if err := s.deleteFolderRecursive(storage, folder.ID); err != nil {
				return err
			}
		}
	}

	// Remove the folder itself
	storage.Folders = append(storage.Folders[:folderIndex], storage.Folders[folderIndex+1:]...)

	return nil
}

// updateFolder updates an existing folder
func (s *Server) updateFolder(id int, name string, parentId *int, expanded *bool) (*Folder, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	// Find the folder
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

	// Update the folder
	if name != "" {
		storage.Folders[folderIndex].Name = name
	}
	if expanded != nil {
		storage.Folders[folderIndex].Expanded = *expanded
	}
	storage.Folders[folderIndex].ParentID = parentId

	if err := s.saveNoteStorage(storage); err != nil {
		return nil, err
	}

	return &storage.Folders[folderIndex], nil
}

// deleteFolder deletes a folder by ID and all its notes and subfolders recursively
func (s *Server) deleteFolder(id int) error {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return err
	}

	if err := s.deleteFolderRecursive(storage, id); err != nil {
		return err
	}

	return s.saveNoteStorage(storage)
}

// getFolderTree returns folders organized as a tree structure
func (s *Server) getFolderTree() ([]map[string]interface{}, error) {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	// Create a map for quick folder lookup
	folderMap := make(map[int]*Folder)
	for i := range storage.Folders {
		folderMap[storage.Folders[i].ID] = &storage.Folders[i]
	}

	// Build the tree structure
	var rootFolders []map[string]interface{}

	for _, folder := range storage.Folders {
		if folder.ParentID == nil {
			// This is a root folder
			rootFolders = append(rootFolders, s.buildFolderNode(folder, folderMap, storage.Notes))
		}
	}

	// Sort root folders by name
	sort.Slice(rootFolders, func(i, j int) bool {
		return rootFolders[i]["name"].(string) < rootFolders[j]["name"].(string)
	})

	return rootFolders, nil
}

// buildFolderNode recursively builds a folder node with its children and note count
func (s *Server) buildFolderNode(folder Folder, folderMap map[int]*Folder, notes []Note) map[string]interface{} {
	// Count notes in this folder
	noteCount := 0
	for _, note := range notes {
		if note.FolderID != nil && *note.FolderID == folder.ID {
			noteCount++
		}
	}

	// Find children folders
	var children []map[string]interface{}
	for _, otherFolder := range folderMap {
		if otherFolder.ParentID != nil && *otherFolder.ParentID == folder.ID {
			children = append(children, s.buildFolderNode(*otherFolder, folderMap, notes))
		}
	}

	// Sort children by name
	sort.Slice(children, func(i, j int) bool {
		return children[i]["name"].(string) < children[j]["name"].(string)
	})

	// Calculate total note count (including children)
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

// moveNotesToFolder moves notes from one folder to another
func (s *Server) moveNotesToFolder(noteIds []int, targetFolderId *int) error {
	storage, err := s.loadNoteStorage()
	if err != nil {
		return err
	}

	// Validate target folder exists if specified
	if targetFolderId != nil {
		folderExists := false
		for _, folder := range storage.Folders {
			if folder.ID == *targetFolderId {
				folderExists = true
				break
			}
		}
		if !folderExists {
			return errors.New("target folder not found")
		}
	}

	// Update notes
	updatedCount := 0
	branch, commit := s.getGitInfo()

	for i := range storage.Notes {
		for _, noteId := range noteIds {
			if storage.Notes[i].ID == noteId {
				storage.Notes[i].FolderID = targetFolderId
				storage.Notes[i].UpdatedAt = time.Now()
				storage.Notes[i].GitBranch = &branch
				storage.Notes[i].GitCommit = &commit
				updatedCount++
				break
			}
		}
	}

	if updatedCount == 0 {
		return errors.New("no notes found to move")
	}

	return s.saveNoteStorage(storage)
}

// searchNotes searches notes by title, content, or tags
func (s *Server) searchNotes(query string, category string, folderId *int) ([]Note, error) {
	if query == "" {
		return s.getNotes(category, "", "")
	}

	storage, err := s.loadNoteStorage()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var matchingNotes []Note

	for _, note := range storage.Notes {
		// Filter by category first
		if category != "" && note.Category != category {
			continue
		}

		// Filter by folder if specified
		if folderId != nil {
			if note.FolderID == nil || *note.FolderID != *folderId {
				continue
			}
		}

		// Search in title
		if strings.Contains(strings.ToLower(note.Title), query) {
			matchingNotes = append(matchingNotes, note)
			continue
		}

		// Search in content
		if strings.Contains(strings.ToLower(note.Content), query) {
			matchingNotes = append(matchingNotes, note)
			continue
		}

		// Search in tags
		for _, tag := range note.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matchingNotes = append(matchingNotes, note)
				break
			}
		}
	}

	// Sort by relevance (title matches first, then by update time)
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

// exportNotes exports notes to JSON format
func (s *Server) exportNotes(folderId *int, category string) (map[string]interface{}, error) {
	var notes []Note
	var err error

	if folderId != nil || category != "" {
		folderIdStr := ""
		if folderId != nil {
			folderIdStr = strconv.Itoa(*folderId)
		}
		notes, err = s.getNotes(category, "", folderIdStr)
	} else {
		notes, err = s.getNotes("", "", "")
	}

	if err != nil {
		return nil, err
	}

	folders, err := s.getFolders()
	if err != nil {
		return nil, err
	}

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"notes":       notes,
		"folders":     folders,
		"categories":  s.getCategories(),
		"total":       len(notes),
	}

	return export, nil
}

// getAllTags returns all unique tags used in notes
func (s *Server) getAllTags() ([]string, error) {
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

	// Sort tags alphabetically
	sort.Strings(tags)

	return tags, nil
}

// loadEnhancedNoteStorage loads the enhanced note storage with history
func (s *Server) loadEnhancedNoteStorage() (*EnhancedNoteStorage, error) {
	if err := s.ensureNotesDir(); err != nil {
		return nil, err
	}

	filePath := s.getNotesFilePath()

	// If file doesn't exist, create empty storage
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage := &EnhancedNoteStorage{
			Notes:         []Note{},
			Folders:       []Folder{},
			History:       []NoteHistoryEntry{},
			NextID:        1,
			NextHistoryID: 1,
		}
		return storage, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %v", err)
	}

	// Try to unmarshal as enhanced storage first
	var enhancedStorage EnhancedNoteStorage
	if err := json.Unmarshal(data, &enhancedStorage); err == nil && enhancedStorage.NextHistoryID > 0 {
		return &enhancedStorage, nil
	}

	// Fallback to old format and migrate
	var oldStorage NoteStorage
	if err := json.Unmarshal(data, &oldStorage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notes: %v", err)
	}

	// Migrate to enhanced format
	enhancedStorage = EnhancedNoteStorage{
		Notes:         oldStorage.Notes,
		Folders:       oldStorage.Folders,
		History:       []NoteHistoryEntry{},
		NextID:        oldStorage.NextID,
		NextHistoryID: 1,
	}

	return &enhancedStorage, nil
}

// saveEnhancedNoteStorage saves the enhanced note storage with history
func (s *Server) saveEnhancedNoteStorage(storage *EnhancedNoteStorage) error {
	if err := s.ensureNotesDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %v", err)
	}

	filePath := s.getNotesFilePath()
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %v", err)
	}

	return nil
}

// addNoteHistoryEntry adds a new history entry
func (s *Server) addNoteHistoryEntry(storage *EnhancedNoteStorage, noteID int, action NoteHistoryAction, changes map[string]interface{}, oldValue, newValue interface{}, message string) {
	author := s.getGitAuthor()
	branch, commit := s.getGitInfo()

	entry := NoteHistoryEntry{
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

// Enhanced note creation with history
func (s *Server) createNoteWithHistory(title, content string, tags []string, category string, folderId *int) (*Note, error) {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	author := s.getGitAuthor()
	branch, commit := s.getGitInfo()

	note := Note{
		ID:        storage.NextID,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      tags,
		Category:  category,
		FolderID:  folderId,
		GitBranch: &branch,
		GitCommit: &commit,
	}

	storage.Notes = append(storage.Notes, note)
	storage.NextID++

	// Add history entry
	changes := map[string]interface{}{
		"title":     title,
		"category":  category,
		"tags":      tags,
		"folder_id": folderId,
	}
	s.addNoteHistoryEntry(storage, note.ID, ActionCreated, changes, nil, note, "Note created")

	if err := s.saveEnhancedNoteStorage(storage); err != nil {
		return nil, err
	}

	return &note, nil
}

// Enhanced note update with history
func (s *Server) updateNoteWithHistory(id int, title, content string, tags []string, category string, pinned bool, folderId *int) (*Note, error) {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	// Find the note
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

	// Track what changed
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

	// Update the note
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

	// Add history entry only if something changed
	if len(changes) > 0 {
		s.addNoteHistoryEntry(storage, id, ActionUpdated, changes, oldNote, storage.Notes[noteIndex], "Note updated")
	}

	if err := s.saveEnhancedNoteStorage(storage); err != nil {
		return nil, err
	}

	return &storage.Notes[noteIndex], nil
}

// Enhanced note deletion with history
func (s *Server) deleteNoteWithHistory(id int) error {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return err
	}

	// Find and remove the note
	noteIndex := -1
	var deletedNote Note
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

	// Add history entry before deletion
	changes := map[string]interface{}{
		"deleted_note": map[string]interface{}{
			"title":    deletedNote.Title,
			"category": deletedNote.Category,
			"tags":     deletedNote.Tags,
		},
	}
	s.addNoteHistoryEntry(storage, id, ActionDeleted, changes, deletedNote, nil, "Note deleted")

	// Remove the note from slice
	storage.Notes = append(storage.Notes[:noteIndex], storage.Notes[noteIndex+1:]...)

	return s.saveEnhancedNoteStorage(storage)
}

// Enhanced note move with history
func (s *Server) moveNotesToFolderWithHistory(noteIds []int, targetFolderId *int) error {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return err
	}

	// Validate target folder exists if specified
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

	// Update notes and track history
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

				// Add history entry
				changes := map[string]interface{}{
					"folder_id": map[string]interface{}{
						"from": oldFolderId,
						"to":   targetFolderId,
					},
				}
				s.addNoteHistoryEntry(storage, noteId, ActionMoved, changes, oldFolderId, targetFolderId, "Note moved to different folder")

				updatedCount++
				break
			}
		}
	}

	if updatedCount == 0 {
		return fmt.Errorf("no notes found to move")
	}

	return s.saveEnhancedNoteStorage(storage)
}

// getNoteHistory retrieves history for a specific note or all notes
func (s *Server) getNoteHistory(filter NoteHistoryFilter) ([]NoteHistoryEntry, error) {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	var filteredHistory []NoteHistoryEntry

	for _, entry := range storage.History {
		// Apply filters
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

	// Sort by timestamp descending (most recent first)
	sort.Slice(filteredHistory, func(i, j int) bool {
		return filteredHistory[i].Timestamp.After(filteredHistory[j].Timestamp)
	})

	// Apply limit and offset
	if filter.Offset > 0 {
		if filter.Offset >= len(filteredHistory) {
			return []NoteHistoryEntry{}, nil
		}
		filteredHistory = filteredHistory[filter.Offset:]
	}

	if filter.Limit > 0 && len(filteredHistory) > filter.Limit {
		filteredHistory = filteredHistory[:filter.Limit]
	}

	return filteredHistory, nil
}

// getNoteHistoryStats generates statistics about note history
func (s *Server) getNoteHistoryStats() (*NoteHistoryStats, error) {
	storage, err := s.loadEnhancedNoteStorage()
	if err != nil {
		return nil, err
	}

	stats := &NoteHistoryStats{
		TotalEntries: len(storage.History),
		ByAction:     make(map[NoteHistoryAction]int),
		ByAuthor:     make(map[string]int),
		ByBranch:     make(map[string]int),
		ByDay:        make(map[string]int),
	}

	// Count by various dimensions
	noteActivity := make(map[int]*NoteActivitySummary)

	for _, entry := range storage.History {
		// By action
		stats.ByAction[entry.Action]++

		// By author
		stats.ByAuthor[entry.Author]++

		// By branch
		if entry.GitBranch != nil {
			stats.ByBranch[*entry.GitBranch]++
		}

		// By day
		day := entry.Timestamp.Format("2006-01-02")
		stats.ByDay[day]++

		// Note activity
		if activity, exists := noteActivity[entry.NoteID]; exists {
			activity.ActionCount++
			if entry.Timestamp.After(activity.LastAction) {
				activity.LastAction = entry.Timestamp
			}
			// Add author if not already present
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
			// Find note title
			noteTitle := "Unknown"
			for _, note := range storage.Notes {
				if note.ID == entry.NoteID {
					noteTitle = note.Title
					break
				}
			}

			noteActivity[entry.NoteID] = &NoteActivitySummary{
				NoteID:      entry.NoteID,
				NoteTitle:   noteTitle,
				ActionCount: 1,
				LastAction:  entry.Timestamp,
				Authors:     []string{entry.Author},
			}
		}
	}

	// Convert note activity to slice and sort by activity count
	for _, activity := range noteActivity {
		stats.MostActiveNotes = append(stats.MostActiveNotes, *activity)
	}
	sort.Slice(stats.MostActiveNotes, func(i, j int) bool {
		return stats.MostActiveNotes[i].ActionCount > stats.MostActiveNotes[j].ActionCount
	})

	// Limit to top 10 most active notes
	if len(stats.MostActiveNotes) > 10 {
		stats.MostActiveNotes = stats.MostActiveNotes[:10]
	}

	// Get recent activity (last 20 entries)
	recentCount := 20
	if len(storage.History) < recentCount {
		recentCount = len(storage.History)
	}

	// Sort history by timestamp descending
	sortedHistory := make([]NoteHistoryEntry, len(storage.History))
	copy(sortedHistory, storage.History)
	sort.Slice(sortedHistory, func(i, j int) bool {
		return sortedHistory[i].Timestamp.After(sortedHistory[j].Timestamp)
	})

	stats.RecentActivity = sortedHistory[:recentCount]

	return stats, nil
}

// Helper functions
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
