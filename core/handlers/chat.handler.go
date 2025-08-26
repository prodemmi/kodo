package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/prodemmi/kodo/core/entities"
	"go.uber.org/zap"
)

type ChatHandler struct {
	logger *zap.Logger
}

func NewChatHandler(logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		logger: logger,
	}
}

func (s *ChatHandler) HandleProjectFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			http.Error(w, "Cannot get working directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	search := r.URL.Query().Get("search")

	var files []entities.ProjectFile
	var err error
	if search != "" {
		files, err = entities.SearchFilesRecursive(dir, search)
	} else {
		files, err = entities.ScanDirOneLevel(dir)
	}
	if err != nil {
		http.Error(w, "Cannot read directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
