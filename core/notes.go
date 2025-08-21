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
	notesDir := filepath.Join(wd, ".kodo")

	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %v", err)
	}

	return nil
}

// getNotesFilePath returns the path to the notes storage file
func (s *Server) getNotesFilePath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, ".kodo", "notes.json")
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
