package core

import "time"

// Note struct
type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Tags      []string  `json:"tags"`
	Category  string    `json:"category"`
	FolderID  *int      `json:"folderId"`
	GitBranch *string   `json:"gitBranch,omitempty"`
	GitCommit *string   `json:"gitCommit,omitempty"`
}

// Folder struct
type Folder struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parentId"`
	Expanded bool   `json:"expanded"`
}

// Category struct
type Category struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
