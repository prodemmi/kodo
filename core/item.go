package core

import (
	"encoding/json"
	"time"
)

// ItemStatus represents the possible states of a TODO item
type ItemStatus string

const (
	StatusTodo       ItemStatus = "todo"
	StatusInProgress ItemStatus = "in-progress"
	StatusDone       ItemStatus = "done"
)

// ItemPriority represents the priority levels
type ItemPriority string

const (
	PriorityLow    ItemPriority = "low"
	PriorityMedium ItemPriority = "medium"
	PriorityHigh   ItemPriority = "high"
)

// ItemType represents the category/purpose of the TODO item
type ItemType string

const (
	// Code Quality & Maintenance
	TypeRefactor   ItemType = "REFACTOR"   // Code needs restructuring
	TypeOptimize   ItemType = "OPTIMIZE"   // Performance improvements needed
	TypeCleanup    ItemType = "CLEANUP"    // Remove dead code, unused imports
	TypeDeprecated ItemType = "DEPRECATED" // Mark old code for removal

	// Bug Fixes & Issues
	TypeBug   ItemType = "BUG"   // Known bug that needs fixing
	TypeFixme ItemType = "FIXME" // Broken code that needs immediate attention

	// Features & Enhancements
	TypeTodo    ItemType = "TODO"    // General task or feature to implement
	TypeFeature ItemType = "FEATURE" // New functionality to add
	TypeEnhance ItemType = "ENHANCE" // Improve existing functionality

	// Documentation & Testing
	TypeDoc     ItemType = "DOC"     // Documentation needed
	TypeTest    ItemType = "TEST"    // Tests needed
	TypeExample ItemType = "EXAMPLE" // Code examples needed

	// Security & Compliance
	TypeSecurity   ItemType = "SECURITY"   // Security concern or improvement
	TypeCompliance ItemType = "COMPLIANCE" // Regulatory or standard compliance

	// Technical Debt & Architecture
	TypeDebt         ItemType = "DEBT"         // Technical debt that should be addressed
	TypeArchitecture ItemType = "ARCHITECTURE" // Architectural changes needed

	// Operations & Infrastructure
	TypeConfig  ItemType = "CONFIG"  // Configuration changes needed
	TypeDeploy  ItemType = "DEPLOY"  // Deployment related tasks
	TypeMonitor ItemType = "MONITOR" // Monitoring/logging improvements

	// General & Miscellaneous
	TypeNote     ItemType = "NOTE"     // Important information or reminder
	TypeQuestion ItemType = "QUESTION" // Something that needs clarification
	TypeIdea     ItemType = "IDEA"     // Future improvement ideas
	TypeReview   ItemType = "REVIEW"   // Code that needs review
)

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

	// Track status changes over time
	History []StatusHistory `json:"history,omitempty"`

	// Computed fields for backward compatibility and convenience
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CurrentUser string    `json:"current_user,omitempty"`
}

// Helper methods for better usability
func (i *Item) IsDone() bool {
	return i.Status == StatusDone
}

func (i *Item) IsInProgress() bool {
	return i.Status == StatusInProgress
}

func (i *Item) IsTodo() bool {
	return i.Status == StatusTodo
}

// GetDoneAt returns when the item was marked as done
func (i *Item) GetDoneAt() *time.Time {
	for j := len(i.History) - 1; j >= 0; j-- {
		if i.History[j].Status == StatusDone {
			return &i.History[j].Timestamp
		}
	}
	return nil
}

// GetDoneBy returns who marked the item as done
func (i *Item) GetDoneBy() string {
	for j := len(i.History) - 1; j >= 0; j-- {
		if i.History[j].Status == StatusDone {
			return i.History[j].User
		}
	}
	return ""
}

// GetInProgressAt returns when the item was marked as in-progress
func (i *Item) GetInProgressAt() *time.Time {
	for j := len(i.History) - 1; j >= 0; j-- {
		if i.History[j].Status == StatusInProgress {
			return &i.History[j].Timestamp
		}
	}
	return nil
}

// GetInProgressBy returns who marked the item as in-progress
func (i *Item) GetInProgressBy() string {
	for j := len(i.History) - 1; j >= 0; j-- {
		if i.History[j].Status == StatusInProgress {
			return i.History[j].User
		}
	}
	return ""
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

// GetPriorityFromType sets priority based on item type and real-world urgency
func GetPriorityFromType(itemType ItemType) ItemPriority {
	switch itemType {
	// Critical - needs immediate attention
	case TypeBug, TypeFixme, TypeSecurity, TypeCompliance:
		return PriorityHigh

	// Important - should be addressed soon
	case TypeRefactor, TypeDebt, TypeArchitecture, TypeDeploy:
		return PriorityMedium

	// Low priority - nice to have
	case TypeNote, TypeIdea, TypeQuestion, TypeExample, TypeDoc:
		return PriorityLow

	// Default medium priority for others
	default:
		return PriorityMedium
	}
}

// GetTypeCategory returns the broader category for grouping
func GetTypeCategory(itemType ItemType) string {
	switch itemType {
	case TypeBug, TypeFixme:
		return "Bug Fixes"
	case TypeRefactor, TypeOptimize, TypeCleanup, TypeDeprecated:
		return "Code Quality"
	case TypeTodo, TypeFeature, TypeEnhance:
		return "Features"
	case TypeDoc, TypeTest, TypeExample:
		return "Documentation & Testing"
	case TypeSecurity, TypeCompliance:
		return "Security & Compliance"
	case TypeDebt, TypeArchitecture:
		return "Technical Debt"
	case TypeConfig, TypeDeploy, TypeMonitor:
		return "Operations"
	case TypeNote, TypeQuestion, TypeIdea, TypeReview:
		return "General"
	default:
		return "Other"
	}
}

// GetTypeDescription returns a human-readable description
func GetTypeDescription(itemType ItemType) string {
	descriptions := map[ItemType]string{
		TypeRefactor:     "Code needs restructuring or reorganization",
		TypeOptimize:     "Performance improvements or optimizations needed",
		TypeCleanup:      "Remove dead code, unused imports, or cleanup mess",
		TypeDeprecated:   "Mark old code for removal or replacement",
		TypeBug:          "Known bug that needs fixing",
		TypeFixme:        "Broken code that needs immediate attention",
		TypeTodo:         "General task or feature to implement",
		TypeFeature:      "New functionality to add",
		TypeEnhance:      "Improve or extend existing functionality",
		TypeDoc:          "Documentation needed or outdated",
		TypeTest:         "Unit tests, integration tests, or test cases needed",
		TypeExample:      "Code examples or usage samples needed",
		TypeSecurity:     "Security vulnerability or improvement needed",
		TypeCompliance:   "Regulatory, legal, or standards compliance issue",
		TypeDebt:         "Technical debt that should be addressed",
		TypeArchitecture: "Architectural changes or improvements needed",
		TypeConfig:       "Configuration changes or setup needed",
		TypeDeploy:       "Deployment, CI/CD, or release related tasks",
		TypeMonitor:      "Monitoring, logging, or observability improvements",
		TypeNote:         "Important information, reminder, or explanation",
		TypeQuestion:     "Something that needs clarification or decision",
		TypeIdea:         "Future improvement ideas or suggestions",
		TypeReview:       "Code that needs review or approval",
	}

	if desc, exists := descriptions[itemType]; exists {
		return desc
	}
	return "Unknown item type"
}

// NewItem creates a new Item with initial status
func NewItem(id int, itemType ItemType, title, description, file string, line int, user string) *Item {
	now := time.Now()

	item := &Item{
		ID:          id,
		Type:        itemType,
		Title:       title,
		Description: description,
		File:        file,
		Line:        line,
		Status:      StatusTodo,
		Priority:    GetPriorityFromType(itemType),
		CreatedAt:   now,
		UpdatedAt:   now,
		CurrentUser: user,
		History: []StatusHistory{{
			Status:    StatusTodo,
			Timestamp: now,
			User:      user,
		}},
	}

	return item
}

// For backward compatibility, you can add these methods to maintain existing API:

// Legacy getter methods (for backward compatibility)
func (i *Item) GetIsDone() bool {
	return i.IsDone()
}

func (i *Item) GetIsInProgress() bool {
	return i.IsInProgress()
}

// MarshalJSON provides custom JSON serialization for backward compatibility
func (i *Item) MarshalJSON() ([]byte, error) {
	type Alias Item

	doneAt := i.GetDoneAt()
	inProgressAt := i.GetInProgressAt()

	return json.Marshal(&struct {
		*Alias
		IsDone       bool       `json:"is_done"`
		DoneAt       *time.Time `json:"done_at"`
		DoneBy       string     `json:"done_by"`
		IsInProgress bool       `json:"is_in_progress"`
		InProgressAt *time.Time `json:"in_progress_at"`
		InProgressBy string     `json:"in_progress_by"`
		FullTitle    string     `json:"full_title"`
	}{
		Alias:        (*Alias)(i),
		IsDone:       i.IsDone(),
		DoneAt:       doneAt,
		DoneBy:       i.GetDoneBy(),
		IsInProgress: i.IsInProgress(),
		InProgressAt: inProgressAt,
		InProgressBy: i.GetInProgressBy(),
		FullTitle:    i.GetFullTitle(),
	})
}
