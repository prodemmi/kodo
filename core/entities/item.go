package entities

import (
	"time"
)

type ItemType string

type ItemStatus string

type ItemPriority string

type StatusHistory struct {
	Status    ItemStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
	User      string     `json:"user"`
}

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

	History []StatusHistory `json:"history,omitempty"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CurrentUser string    `json:"current_user,omitempty"`
}

func (i *Item) GetIsDone() bool {
	return i.IsDone
}

func (i *Item) SetStatus(status ItemStatus, user string) {
	if i.Status == status {
		return
	}

	now := time.Now()
	i.Status = status
	i.UpdatedAt = now
	i.CurrentUser = user

	i.History = append(i.History, StatusHistory{
		Status:    status,
		Timestamp: now,
		User:      user,
	})
}

func (i *Item) GetFullTitle() string {
	return string(i.Type) + ": " + i.Title
}

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
