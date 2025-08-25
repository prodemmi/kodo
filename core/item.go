package core

import (
	"time"
)

// ItemType represents the category/purpose of the TODO item
type ItemType string

// ItemStatus represents the possible states of a TODO item
type ItemStatus string

// ItemPriority represents the priority levels
type ItemPriority string

// StatusHistory tracks when status changes occurred
type StatusHistory struct {
	Status    ItemStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
	User      string     `json:"user"`
}

// Item represents a TODO/FIXME/BUG/NOTE item found in code
type Item struct {
	ID          int          `json:"id"`
	Type        ItemType     `json:"type"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	File        string       `json:"file"`
	Line        int          `json:"line"`
	Status      ItemStatus   `json:"status"`
	Priority    ItemPriority `json:"priority"`
	IsDone      bool         `json:"is_done"`
	DoneAt      *time.Time   `json:"done_at"`
	DoneBy      *string      `json:"done_by"`

	// Track status changes over time
	History []StatusHistory `json:"history,omitempty"`

	// Computed fields for backward compatibility and convenience
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CurrentUser string    `json:"current_user,omitempty"`
}

// Helper methods for better usability
func (i *Item) GetIsDone() bool {
	return i.IsDone
}

// SetStatus changes the status and adds to history
func (i *Item) SetStatus(status ItemStatus, user string) {
	if i.Status == status {
		return // No change needed
	}

	now := time.Now()
	i.Status = status
	i.UpdatedAt = now
	i.CurrentUser = user

	// Add to history
	i.History = append(i.History, StatusHistory{
		Status:    status,
		Timestamp: now,
		User:      user,
	})
}

// GetFullTitle returns type and title combined
func (i *Item) GetFullTitle() string {
	return string(i.Type) + ": " + i.Title
}

// NewItem creates a new Item with initial status
func NewItem(id int, itemType ItemType, itemStatus ItemStatus, itemPriority ItemPriority, title, description, file string, line int, user string) *Item {
	now := time.Now()

	item := &Item{
		ID:          id,
		Type:        itemType,
		Title:       title,
		Description: description,
		File:        file,
		Line:        line,
		Status:      itemStatus,
		Priority:    itemPriority,
		CreatedAt:   now,
		UpdatedAt:   now,
		CurrentUser: user,
		History: []StatusHistory{{
			Status:    itemStatus,
			Timestamp: now,
			User:      user,
		}},
	}

	return item
}
