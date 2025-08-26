package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/prodemmi/kodo/core/entities"
	"github.com/prodemmi/kodo/core/services"
	"go.uber.org/zap"
)

type NoteHandler struct {
	logger        *zap.Logger
	noteService   *services.NoteService
	remoteService *services.RemoteService
}

func NewNoteHandler(logger *zap.Logger,
	noteService *services.NoteService, remoteService *services.RemoteService) *NoteHandler {
	return &NoteHandler{
		logger:        logger,
		noteService:   noteService,
		remoteService: remoteService,
	}
}

func (s *NoteHandler) HandleNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":

		category := r.URL.Query().Get("category")
		tag := r.URL.Query().Get("tag")
		folderId := r.URL.Query().Get("folderId")

		notes, err := s.noteService.GetNotes(category, tag, folderId)
		if err != nil {
			s.logger.Error("Failed to get notes", zap.Error(err))
			http.Error(w, "Failed to get notes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"notes": notes,
			"count": len(notes),
		})

	case "POST":

		var noteReq struct {
			Title    string   `json:"title"`
			Content  string   `json:"content"`
			Tags     []string `json:"tags"`
			Category string   `json:"category"`
			FolderID *int     `json:"folderId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
			s.logger.Error("Invalid JSON for note creation", zap.Error(err))
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if noteReq.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		note, err := s.noteService.CreateNoteWithHistory(noteReq.Title, noteReq.Content, noteReq.Tags, noteReq.Category, noteReq.FolderID, nil)
		if err != nil {
			s.logger.Error("Failed to create note", zap.Error(err))
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		s.logger.Info("Note created with history tracking", zap.Int("id", note.ID), zap.String("title", note.Title))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"note":   note,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *NoteHandler) HandleNoteUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq struct {
		ID       int      `json:"id"`
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Tags     []string `json:"tags"`
		Category string   `json:"category"`
		Pinned   bool     `json:"pinned"`
		FolderID *int     `json:"folderId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		s.logger.Error("Invalid JSON for note update", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if updateReq.ID <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	if updateReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	note, err := s.noteService.UpdateNoteWithHistory(updateReq.ID, updateReq.Title, updateReq.Content, updateReq.Tags, updateReq.Category, updateReq.Pinned, updateReq.FolderID)
	if err != nil {
		s.logger.Error("Failed to update note", zap.Int("id", updateReq.ID), zap.Error(err))
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Note updated with history tracking", zap.Int("id", note.ID), zap.String("title", note.Title))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"note":   note,
	})
}

func (s *NoteHandler) HandleNoteDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteReq struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		s.logger.Error("Invalid JSON for note deletion", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if deleteReq.ID <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	err := s.noteService.DeleteNoteWithHistory(deleteReq.ID)
	if err != nil {
		s.logger.Error("Failed to delete note", zap.Int("id", deleteReq.ID), zap.Error(err))
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Note deleted with history tracking", zap.Int("id", deleteReq.ID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Note deleted successfully",
	})
}

func (s *NoteHandler) HandleMoveNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var moveReq struct {
		NoteIds        []int `json:"noteIds"`
		TargetFolderId *int  `json:"targetFolderId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(moveReq.NoteIds) == 0 {
		http.Error(w, "No notes specified", http.StatusBadRequest)
		return
	}

	err := s.noteService.MoveNotesToFolderWithHistory(moveReq.NoteIds, moveReq.TargetFolderId)
	if err != nil {
		s.logger.Error("Failed to move notes", zap.Error(err))
		if err.Error() == "target folder not found" {
			http.Error(w, "Target folder not found", http.StatusNotFound)
			return
		}
		if err.Error() == "no notes found to move" {
			http.Error(w, "No notes found to move", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to move notes", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Notes moved with history tracking", zap.Ints("noteIds", moveReq.NoteIds), zap.Any("targetFolderId", moveReq.TargetFolderId))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Moved %d notes", len(moveReq.NoteIds)),
	})
}

func (s *NoteHandler) HandleExportNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	folderIdStr := r.URL.Query().Get("folderId")

	var folderId *int
	if folderIdStr != "" {
		if id, err := strconv.Atoi(folderIdStr); err == nil {
			folderId = &id
		}
	}

	export, err := s.noteService.ExportNotes(folderId, category)
	if err != nil {
		s.logger.Error("Failed to export notes", zap.Error(err))
		http.Error(w, "Failed to export notes", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("notes-export-%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	json.NewEncoder(w).Encode(export)
}

func (s *NoteHandler) HandleNoteTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tags, err := s.noteService.GetAllTags()
	if err != nil {
		s.logger.Error("Failed to get tags", zap.Error(err))
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tags":  tags,
		"count": len(tags),
	})
}

func (s *NoteHandler) HandleNoteHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var filter entities.NoteHistoryFilter

	if noteIdStr := r.URL.Query().Get("noteId"); noteIdStr != "" {
		if noteId, err := strconv.Atoi(noteIdStr); err == nil {
			filter.NoteID = &noteId
		}
	}

	if actionStr := r.URL.Query().Get("action"); actionStr != "" {
		action := entities.NoteHistoryAction(actionStr)
		filter.Action = &action
	}

	if author := r.URL.Query().Get("author"); author != "" {
		filter.Author = &author
	}

	if branch := r.URL.Query().Get("branch"); branch != "" {
		filter.GitBranch = &branch
	}

	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		if since, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			filter.Since = &since
		}
	}

	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		if until, err := time.Parse(time.RFC3339, untilStr); err == nil {
			filter.Until = &until
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 50
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	history, err := s.noteService.GetNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get note history", zap.Error(err))
		http.Error(w, "Failed to get note history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
		"count":   len(history),
		"filter":  filter,
	})
}

func (s *NoteHandler) HandleNoteHistoryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := s.noteService.GetNoteHistoryStats()
	if err != nil {
		s.logger.Error("Failed to get note history history", zap.Error(err))
		http.Error(w, "Failed to get note history history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (s *NoteHandler) HandleNoteWithHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	noteIdStr := r.URL.Query().Get("id")
	if noteIdStr == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}

	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	storage, err := s.noteService.LoadEnhancedNoteStorage()
	if err != nil {
		s.logger.Error("Failed to load note storage", zap.Error(err))
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	var note *entities.Note
	for _, n := range storage.Notes {
		if n.ID == noteId {
			note = &n
			break
		}
	}

	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	filter := entities.NoteHistoryFilter{NoteID: &noteId}
	history, err := s.noteService.GetNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get note history", zap.Error(err))
		http.Error(w, "Failed to get note history", http.StatusInternalServerError)
		return
	}

	noteWithHistory := entities.NoteWithHistory{
		Note:    *note,
		History: history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"note":          noteWithHistory,
		"history_count": len(history),
	})
}

func (s *NoteHandler) HandleAuthorActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	author := r.URL.Query().Get("author")
	if author == "" {
		http.Error(w, "Author parameter is required", http.StatusBadRequest)
		return
	}

	filter := entities.NoteHistoryFilter{Author: &author, Limit: 100}
	history, err := s.noteService.GetNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get author activity", zap.Error(err))
		http.Error(w, "Failed to get author activity", http.StatusInternalServerError)
		return
	}

	actionCount := make(map[entities.NoteHistoryAction]int)
	noteCount := make(map[int]bool)
	dayCount := make(map[string]int)

	for _, entry := range history {
		actionCount[entry.Action]++
		noteCount[entry.NoteID] = true
		day := entry.Timestamp.Format("2006-01-02")
		dayCount[day]++
	}

	activity := map[string]interface{}{
		"author":         author,
		"total_actions":  len(history),
		"notes_affected": len(noteCount),
		"by_action":      actionCount,
		"by_day":         dayCount,
		"recent_history": history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

func (s *NoteHandler) HandleBranchActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	branch := r.URL.Query().Get("branch")
	if branch == "" {
		http.Error(w, "Branch parameter is required", http.StatusBadRequest)
		return
	}

	filter := entities.NoteHistoryFilter{GitBranch: &branch, Limit: 100}
	history, err := s.noteService.GetNoteHistory(filter)
	if err != nil {
		s.logger.Error("Failed to get branch activity", zap.Error(err))
		http.Error(w, "Failed to get branch activity", http.StatusInternalServerError)
		return
	}

	actionCount := make(map[entities.NoteHistoryAction]int)
	authorCount := make(map[string]int)
	noteCount := make(map[int]bool)

	for _, entry := range history {
		actionCount[entry.Action]++
		authorCount[entry.Author]++
		noteCount[entry.NoteID] = true
	}

	activity := map[string]interface{}{
		"branch":         branch,
		"total_actions":  len(history),
		"notes_affected": len(noteCount),
		"by_action":      actionCount,
		"by_author":      authorCount,
		"recent_history": history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

func (s *NoteHandler) HandleCleanupHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OlderThanDays int `json:"older_than_days"`
		KeepMinimum   int `json:"keep_minimum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.OlderThanDays <= 0 {
		req.OlderThanDays = 90
	}
	if req.KeepMinimum <= 0 {
		req.KeepMinimum = 10
	}

	storage, err := s.noteService.LoadEnhancedNoteStorage()
	if err != nil {
		s.logger.Error("Failed to load note storage", zap.Error(err))
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -req.OlderThanDays)

	noteHistory := make(map[int][]entities.NoteHistoryEntry)
	for _, entry := range storage.History {
		noteHistory[entry.NoteID] = append(noteHistory[entry.NoteID], entry)
	}

	var filteredHistory []entities.NoteHistoryEntry
	removedCount := 0

	for _, entries := range noteHistory {

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})

		keptCount := 0
		for _, entry := range entries {

			if keptCount < req.KeepMinimum || entry.Timestamp.After(cutoff) {
				filteredHistory = append(filteredHistory, entry)
				keptCount++
			} else {
				removedCount++
			}
		}
	}

	storage.History = filteredHistory

	if err := s.noteService.SaveEnhancedNoteStorage(storage); err != nil {
		s.logger.Error("Failed to save cleaned history", zap.Error(err))
		http.Error(w, "Failed to save changes", http.StatusInternalServerError)
		return
	}

	s.logger.Info("History cleanup completed", zap.Int("removed_entries", removedCount), zap.Int("remaining_entries", len(filteredHistory)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"removed_entries":   removedCount,
		"remaining_entries": len(filteredHistory),
		"message":           fmt.Sprintf("Removed %d old history entries", removedCount),
	})
}

func (s *NoteHandler) HandleNoteSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	folderIdStr := r.URL.Query().Get("folderId")

	var folderId *int
	if folderIdStr != "" {
		if id, err := strconv.Atoi(folderIdStr); err == nil {
			folderId = &id
		}
	}

	notes, err := s.noteService.SearchNotes(query, category, folderId)
	if err != nil {
		s.logger.Error("Failed to search notes", zap.Error(err))
		http.Error(w, "Failed to search notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notes": notes,
		"count": len(notes),
		"query": query,
	})
}

func (s *NoteHandler) handleNoteStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := s.noteService.GetNoteStats()
	if err != nil {
		s.logger.Error("Failed to get note history", zap.Error(err))
		http.Error(w, "Failed to get note history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (s *NoteHandler) HandleSyncNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := s.remoteService.SyncIssuesWithNotes()
	if err != nil {
		s.logger.Error("Failed to sybc notes", zap.Error(err))
		http.Error(w, "Failed to sybc notes", http.StatusInternalServerError)
		return
	}

	s.noteService.SaveNoteStorage(s.noteService)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (s *NoteHandler) HandleFolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		folders, err := s.noteService.GetFolders()
		if err != nil {
			s.logger.Error("Failed to get folders", zap.Error(err))
			http.Error(w, "Failed to get folders", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"folders": folders,
			"count":   len(folders),
		})

	case "POST":
		var folderReq struct {
			Name     string `json:"name"`
			ParentID *int   `json:"parentId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&folderReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if folderReq.Name == "" {
			http.Error(w, "Folder name is required", http.StatusBadRequest)
			return
		}

		folder, err := s.noteService.CreateFolder(folderReq.Name, folderReq.ParentID)
		if err != nil {
			s.logger.Error("Failed to create folder", zap.Error(err))
			http.Error(w, "Failed to create folder", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"folder": folder,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *NoteHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories := s.noteService.GetCategories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
		"count":      len(categories),
	})
}

func (s *NoteHandler) HandleFolderUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateReq struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		ParentID *int   `json:"parentId"`
		Expanded *bool  `json:"expanded"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if updateReq.ID <= 0 {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := s.noteService.UpdateFolder(updateReq.ID, updateReq.Name, updateReq.ParentID, updateReq.Expanded)
	if err != nil {
		s.logger.Error("Failed to update folder", zap.Int("id", updateReq.ID), zap.Error(err))
		if err.Error() == "folder not found" {
			http.Error(w, "Folder not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update folder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"folder": folder,
	})
}

func (s *NoteHandler) HandleFolderDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteReq struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if deleteReq.ID <= 0 {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	err := s.noteService.DeleteFolder(deleteReq.ID)
	if err != nil {
		s.logger.Error("Failed to delete folder", zap.Int("id", deleteReq.ID), zap.Error(err))
		if err.Error() == "folder not found" {
			http.Error(w, "Folder not found", http.StatusNotFound)
			return
		}
		if err.Error() == "folder has notes" {
			http.Error(w, "Cannot delete folder that contains notes", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to delete folder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Folder deleted successfully",
	})
}

func (s *NoteHandler) HandleFolderTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tree, err := s.noteService.GetFolderTree()
	if err != nil {
		s.logger.Error("Failed to get folder tree", zap.Error(err))
		http.Error(w, "Failed to get folder tree", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tree":  tree,
		"count": len(tree),
	})
}
