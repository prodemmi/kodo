package services

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/prodemmi/kodo/core/entities"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type RemoteService struct {
	logger      *zap.Logger
	settings    *SettingsService
	noteStorage *NoteStorage
}

func NewRemoteManager(logger *zap.Logger, settings *SettingsService, noteStorage *NoteStorage) *RemoteService {
	return &RemoteService{
		logger:      logger,
		settings:    settings,
		noteStorage: noteStorage,
	}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	NotesCreated  int
	NotesUpdated  int
	IssuesCreated int
	IssuesUpdated int
	IssuesClosed  int
	Conflicts     []SyncConflict
	Errors        []error
}

// SyncConflict represents a conflict during sync
type SyncConflict struct {
	Type        string
	NoteID      *int
	IssueNumber *int
	Description string
	LocalData   interface{}
	RemoteData  interface{}
}

// Enhanced sync method with bidirectional sync and conflict resolution
func (r *RemoteService) SyncIssuesWithNotes() (*SyncResult, error) {
	result := &SyncResult{}

	// Load settings and validate
	settings := r.settings.LoadSettings()
	githubToken := settings.GithubAuth.Token

	if !settings.CodeScanSettings.SyncEnabled {
		return nil, fmt.Errorf("sync not enabled")
	}

	if githubToken == "" {
		return nil, fmt.Errorf("GitHub token not configured")
	}

	// Get repository info
	owner, repo, err := getRepoOwnerAndName()
	if err != nil {
		return nil, fmt.Errorf("failed to get repo info: %v", err)
	}

	r.logger.Info("Starting GitHub sync",
		zap.String("owner", owner),
		zap.String("repo", repo))

	// Initialize GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get all local notes
	localNotes, err := r.noteStorage.GetNotes("", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to load local notes: %v", err)
	}

	// Get all GitHub issues (including closed ones)
	allIssues, err := r.getAllGitHubIssues(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to load GitHub issues: %v", err)
	}

	r.logger.Info("Loaded data for sync",
		zap.Int("local_notes", len(localNotes)),
		zap.Int("github_issues", len(allIssues)))

	// Perform bidirectional sync
	err = r.performBidirectionalSync(ctx, client, owner, repo, localNotes, allIssues, result)
	if err != nil {
		return result, fmt.Errorf("sync failed: %v", err)
	}

	r.logger.Info("Sync completed",
		zap.Int("notes_created", result.NotesCreated),
		zap.Int("notes_updated", result.NotesUpdated),
		zap.Int("issues_created", result.IssuesCreated),
		zap.Int("issues_updated", result.IssuesUpdated),
		zap.Int("issues_closed", result.IssuesClosed),
		zap.Int("conflicts", len(result.Conflicts)))

	return result, nil
}

// getAllGitHubIssues retrieves all issues from GitHub (open and closed)
func (r *RemoteService) getAllGitHubIssues(ctx context.Context, client *github.Client, owner, repo string) ([]*github.Issue, error) {
	var allIssues []*github.Issue

	// Get open issues
	openIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, err
	}
	allIssues = append(allIssues, openIssues...)

	// Get closed issues (limit to recent ones to avoid too many)
	closedIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "closed",
		Since:       time.Now().AddDate(0, -3, 0), // Last 3 months
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, err
	}
	allIssues = append(allIssues, closedIssues...)

	return allIssues, nil
}

// performBidirectionalSync handles the main sync logic
func (r *RemoteService) performBidirectionalSync(ctx context.Context, client *github.Client, owner, repo string, localNotes []entities.Note, githubIssues []*github.Issue, result *SyncResult) error {
	// Create maps for efficient lookups
	notesByGitHubID := make(map[int]*entities.Note)
	notesByTitle := make(map[string]*entities.Note)
	issuesByNumber := make(map[int]*github.Issue)
	issuesByNoteID := make(map[string]*github.Issue)
	processedIssues := make(map[int]bool)

	// Index local notes
	for i := range localNotes {
		note := &localNotes[i]

		// Index by GitHub issue number if it exists
		if note.GitHubIssueNumber != nil {
			notesByGitHubID[*note.GitHubIssueNumber] = note
		}

		// Index by title for fallback matching
		notesByTitle[note.Title] = note
	}

	// Index GitHub issues
	for _, issue := range githubIssues {
		issuesByNumber[issue.GetNumber()] = issue

		// Check if issue has a note ID in its body
		noteID := r.extractNoteIDFromIssue(issue)
		if noteID != "" {
			issuesByNoteID[noteID] = issue
		}
	}

	// Log mappings for debugging
	r.logger.Info("Sync mappings",
		zap.Int("notes_by_github_id", len(notesByGitHubID)),
		zap.Int("notes_by_title", len(notesByTitle)),
		zap.Int("issues_by_number", len(issuesByNumber)),
		zap.Int("issues_by_note_id", len(issuesByNoteID)))

	// Phase 1: Update existing pairs (note <-> issue)
	err := r.syncExistingPairs(ctx, client, owner, repo, notesByGitHubID, issuesByNumber, processedIssues, result)
	if err != nil {
		return err
	}

	// Phase 2: Create GitHub issues for new notes
	err = r.createIssuesForNewNotes(ctx, client, owner, repo, localNotes, notesByGitHubID, issuesByNoteID, processedIssues, result)
	if err != nil {
		return err
	}

	// Phase 3: Create notes for new GitHub issues (that don't have associated notes)
	err = r.createNotesForNewIssues(ctx, githubIssues, notesByGitHubID, issuesByNoteID, processedIssues, result)
	if err != nil {
		return err
	}

	return nil
}

// syncExistingPairs syncs notes and issues that are already paired
func (r *RemoteService) syncExistingPairs(ctx context.Context, client *github.Client, owner, repo string, notesByGitHubID map[int]*entities.Note, issuesByNumber map[int]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for issueNumber, note := range notesByGitHubID {
		issue, exists := issuesByNumber[issueNumber]
		if !exists {
			// Issue was deleted on GitHub
			r.logger.Warn("GitHub issue no longer exists", zap.Int("issue_number", issueNumber), zap.Int("note_id", note.ID))

			// Remove the GitHub reference from the note
			note.GitHubIssueNumber = nil
			note.GitHubIssueURL = nil
			_, err := r.noteStorage.UpdateNoteWithHistory(note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to update note %d: %v", note.ID, err))
			} else {
				result.NotesUpdated++
			}
			continue
		}

		// Mark this issue as processed
		processedIssues[issueNumber] = true

		// Check if sync is needed
		needsSync, direction := r.needsSyncCheck(note, issue)
		if !needsSync {
			continue
		}

		r.logger.Info("Sync check",
			zap.Int("note_id", note.ID),
			zap.Int("issue_number", issueNumber),
			zap.String("direction", direction),
			zap.Time("note_updated", note.UpdatedAt),
			zap.Time("issue_updated", issue.GetUpdatedAt().Time))

		switch direction {
		case "note_to_issue":
			err := r.updateIssueFromNote(ctx, client, owner, repo, note, issue, result)
			if err != nil {
				result.Errors = append(result.Errors, err)
			}
		case "issue_to_note":
			err := r.updateNoteFromIssue(issue, note, result)
			if err != nil {
				result.Errors = append(result.Errors, err)
			}
		case "conflict":
			// Prioritize based on which was updated more recently
			if note.UpdatedAt.After(issue.GetUpdatedAt().Time) {
				r.logger.Info("Resolving conflict: prioritizing local note",
					zap.Int("note_id", note.ID),
					zap.Int("issue_number", issueNumber))
				err := r.updateIssueFromNote(ctx, client, owner, repo, note, issue, result)
				if err != nil {
					result.Errors = append(result.Errors, err)
				}
			} else {
				r.logger.Info("Resolving conflict: prioritizing GitHub issue",
					zap.Int("note_id", note.ID),
					zap.Int("issue_number", issueNumber))
				err := r.updateNoteFromIssue(issue, note, result)
				if err != nil {
					result.Errors = append(result.Errors, err)
				}
			}
		}
	}

	return nil
}

// createIssuesForNewNotes creates GitHub issues for notes that don't have associated issues
func (r *RemoteService) createIssuesForNewNotes(ctx context.Context, client *github.Client, owner, repo string, localNotes []entities.Note, notesByGitHubID map[int]*entities.Note, issuesByNoteID map[string]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for _, note := range localNotes {
		// Skip if note already has a GitHub issue
		if note.GitHubIssueNumber != nil {
			continue
		}

		// Check if there's already an issue for this note (by note ID in issue body)
		noteIDStr := fmt.Sprintf("note-%d", note.ID)
		if existingIssue, exists := issuesByNoteID[noteIDStr]; exists {
			// Link the existing issue to the note
			err := r.linkNoteToIssue(&note, existingIssue)
			if err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			// Mark this issue as processed
			processedIssues[existingIssue.GetNumber()] = true
			result.NotesUpdated++
			continue
		}

		// Create new GitHub issue
		err := r.createGitHubIssueFromNote(ctx, client, owner, repo, &note, result)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	return nil
}

// createNotesForNewIssues creates notes for GitHub issues that don't have associated notes
func (r *RemoteService) createNotesForNewIssues(ctx context.Context, githubIssues []*github.Issue, notesByGitHubID map[int]*entities.Note, issuesByNoteID map[string]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for _, issue := range githubIssues {
		issueNumber := issue.GetNumber()

		// Skip if already processed or if note already exists for this issue
		if processedIssues[issueNumber] {
			continue
		}

		// Skip if this is a pull request
		if issue.PullRequestLinks != nil {
			continue
		}

		// Check if issue has note-managed label (created by our system)
		if !r.hasNoteManagedLabel(issue) {
			// This is a manual GitHub issue, create a note for it
			err := r.createNoteFromGitHubIssue(issue, result)
			if err != nil {
				result.Errors = append(result.Errors, err)
			}
		}
	}

	return nil
}

// hasNoteManagedLabel checks if an issue has the note-managed label
func (r *RemoteService) hasNoteManagedLabel(issue *github.Issue) bool {
	for _, label := range issue.Labels {
		if label.GetName() == "note-managed" {
			return true
		}
	}
	return false
}

// createNoteFromGitHubIssue creates a new note from a GitHub issue
func (r *RemoteService) createNoteFromGitHubIssue(issue *github.Issue, result *SyncResult) error {
	// Extract content from issue body (remove any metadata)
	content := r.extractNoteContentFromIssue(issue)

	// Convert labels to tags and category
	tags, category := r.convertLabelsToNoteTags(issue.Labels)

	author := "unknown"
	// Get issue author
	if issue.User != nil && issue.User.Login != nil {
		author = fmt.Sprintf("%s %s", issue.User.GetLogin(), issue.User.GetEmail())
	}

	// Create the note
	issueNumber := issue.GetNumber()
	now := time.Now().UTC()
	status := issue.GetState()

	createdNote, err := r.noteStorage.CreateNoteWithHistory(
		issue.GetTitle(),
		content,
		tags,
		category,
		nil,
		&author,
	)
	if err != nil {
		return fmt.Errorf("failed to create note from issue #%d: %v", issueNumber, err)
	}

	// Update the created note with GitHub information
	note := r.noteStorage.getNoteByID(createdNote.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve created note: %v", err)
	}

	note.GitHubIssueNumber = &issueNumber
	note.GitHubIssueURL = issue.HTMLURL
	note.GitHubLastSync = &now
	note.Status = &status

	_, err = r.noteStorage.UpdateNoteWithHistory(
		note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update note with GitHub info: %v", err)
	}

	result.NotesCreated++
	return nil
}

// needsSyncCheck determines if a note-issue pair needs syncing and in which direction
func (r *RemoteService) needsSyncCheck(note *entities.Note, issue *github.Issue) (bool, string) {
	// Ensure UTC for consistent comparison
	issueUpdated := issue.GetUpdatedAt().Time.UTC()
	noteUpdated := note.UpdatedAt.UTC()

	// If GitHub sync time is tracked, use it for comparison
	var lastSyncTime time.Time
	if note.GitHubLastSync != nil {
		lastSyncTime = note.GitHubLastSync.UTC()
	}

	// Check if either was modified since last sync (with 10-second buffer)
	noteModifiedSinceSync := lastSyncTime.IsZero() || noteUpdated.After(lastSyncTime.Add(10*time.Second))
	issueModifiedSinceSync := lastSyncTime.IsZero() || issueUpdated.After(lastSyncTime.Add(10*time.Second))

	if !noteModifiedSinceSync && !issueModifiedSinceSync {
		return false, ""
	}

	if noteModifiedSinceSync && !issueModifiedSinceSync {
		return true, "note_to_issue"
	}

	if !noteModifiedSinceSync && issueModifiedSinceSync {
		return true, "issue_to_note"
	}

	// Both modified, decide based on which was updated more recently
	if noteUpdated.After(issueUpdated) {
		return true, "note_to_issue"
	}
	if issueUpdated.After(noteUpdated) {
		return true, "issue_to_note"
	}

	// Timestamps are equal within buffer, treat as conflict
	return true, "conflict"
}

// updateIssueFromNote updates a GitHub issue with note data
func (r *RemoteService) updateIssueFromNote(ctx context.Context, client *github.Client, owner, repo string, note *entities.Note, issue *github.Issue, result *SyncResult) error {
	// Prepare issue update
	body := r.formatIssueBodyFromNote(note)
	labels := r.convertNoteTagsToLabels(note.Tags, note.Category)

	// Determine issue state
	state := "open"
	if note.Status != nil && *note.Status == "closed" {
		state = "closed"
	}

	issueRequest := &github.IssueRequest{
		Title:  &note.Title,
		Body:   &body,
		Labels: &labels,
		State:  &state,
	}

	// Log the update request
	r.logger.Info("Updating issue", zap.Any("issue_request", issueRequest))

	// Update the issue with retry logic
	var updatedIssue *github.Issue
	var err error
	for retries := 3; retries > 0; retries-- {
		updatedIssue, _, err = client.Issues.Edit(ctx, owner, repo, issue.GetNumber(), issueRequest)
		if err == nil {
			break
		}
		r.logger.Warn("Retrying issue update", zap.Int("retries_left", retries-1), zap.Error(err))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to update issue #%d: %v", issue.GetNumber(), err)
	}

	// Update note's sync information
	now := time.Now().UTC()
	note.GitHubLastSync = &now
	note.GitHubIssueURL = updatedIssue.HTMLURL

	_, err = r.noteStorage.UpdateNoteWithHistory(note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID)
	if err != nil {
		return fmt.Errorf("failed to update note sync info: %v", err)
	}

	result.IssuesUpdated++
	return nil
}

// updateNoteFromIssue updates a note with GitHub issue data
func (r *RemoteService) updateNoteFromIssue(issue *github.Issue, note *entities.Note, result *SyncResult) error {
	// Extract note content from issue body (remove metadata)
	content := r.extractNoteContentFromIssue(issue)

	// Convert labels to tags and category
	tags, category := r.convertLabelsToNoteTags(issue.Labels)

	// Determine status
	status := issue.GetState()

	// Log update details
	r.logger.Info("Updating note",
		zap.Int("note_id", note.ID),
		zap.String("title", issue.GetTitle()),
		zap.String("content", content),
		zap.Strings("tags", tags),
		zap.String("category", category))

	// Update note
	now := time.Now().UTC()
	note.GitHubLastSync = &now
	note.Status = &status
	note.Tags = tags
	note.Category = category

	_, err := r.noteStorage.UpdateNoteWithHistory(note.ID, issue.GetTitle(), content, tags, category, note.Pinned, note.FolderID)
	if err != nil {
		return fmt.Errorf("failed to update note from issue: %v", err)
	}

	result.NotesUpdated++
	return nil
}

// createGitHubIssueFromNote creates a new GitHub issue from a note
func (r *RemoteService) createGitHubIssueFromNote(ctx context.Context, client *github.Client, owner, repo string, note *entities.Note, result *SyncResult) error {
	body := r.formatIssueBodyFromNote(note)
	labels := r.convertNoteTagsToLabels(note.Tags, note.Category)

	state := "open"
	if note.Status != nil && *note.Status == "closed" {
		state = "closed"
	}

	issueRequest := &github.IssueRequest{
		Title:  &note.Title,
		Body:   &body,
		Labels: &labels,
		State:  &state,
	}

	// Create the issue with retry logic
	var createdIssue *github.Issue
	var err error
	for retries := 3; retries > 0; retries-- {
		createdIssue, _, err = client.Issues.Create(ctx, owner, repo, issueRequest)
		if err == nil {
			break
		}
		r.logger.Warn("Retrying issue creation", zap.Int("retries_left", retries-1), zap.Error(err))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to create issue for note %d: %v", note.ID, err)
	}

	// Update note with issue information
	err = r.linkNoteToIssue(note, createdIssue)
	if err != nil {
		return fmt.Errorf("failed to link note to created issue: %v", err)
	}

	result.IssuesCreated++
	return nil
}

// linkNoteToIssue establishes the bidirectional link between a note and GitHub issue
func (r *RemoteService) linkNoteToIssue(note *entities.Note, issue *github.Issue) error {
	now := time.Now().UTC()
	issueNumber := issue.GetNumber()

	note.GitHubIssueNumber = &issueNumber
	note.GitHubIssueURL = issue.HTMLURL
	note.GitHubLastSync = &now

	// Determine status from issue state
	status := issue.GetState()
	note.Status = &status

	_, err := r.noteStorage.UpdateNoteWithHistory(
		note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID,
	)
	return err
}

// formatIssueBodyFromNote formats a note into a GitHub issue body with metadata
func (r *RemoteService) formatIssueBodyFromNote(note *entities.Note) string {
	var body strings.Builder

	// Add note content
	body.WriteString(note.Content)
	body.WriteString("\n\n")

	// Add metadata section
	body.WriteString("---\n")
	body.WriteString("<!-- Note Metadata - DO NOT EDIT MANUALLY -->\n")
	body.WriteString(fmt.Sprintf("Note-ID: note-%d\n", note.ID))
	body.WriteString(fmt.Sprintf("Author: %s\n", note.Author))
	body.WriteString(fmt.Sprintf("Created: %s\n", note.CreatedAt.Format(time.RFC3339)))
	body.WriteString(fmt.Sprintf("Updated: %s\n", note.UpdatedAt.Format(time.RFC3339)))
	body.WriteString(fmt.Sprintf("Category: %s\n", note.Category))

	if len(note.Tags) > 0 {
		body.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(note.Tags, ",")))
	}

	if note.GitBranch != nil {
		body.WriteString(fmt.Sprintf("Git-Branch: %s\n", *note.GitBranch))
	}
	if note.GitCommit != nil {
		body.WriteString(fmt.Sprintf("Git-Commit: %s\n", *note.GitCommit))
	}

	body.WriteString("<!-- End Note Metadata -->\n")

	return body.String()
}

// extractNoteContentFromIssue extracts the actual note content from an issue body
func (r *RemoteService) extractNoteContentFromIssue(issue *github.Issue) string {
	body := issue.GetBody()

	// Find the metadata section and remove it
	metadataStart := strings.Index(body, "---\n<!-- Note Metadata")
	if metadataStart != -1 {
		return strings.TrimSpace(body[:metadataStart])
	}

	return body
}

// extractNoteIDFromIssue extracts the note ID from an issue's metadata
func (r *RemoteService) extractNoteIDFromIssue(issue *github.Issue) string {
	body := issue.GetBody()

	// Look for Note-ID: note-X pattern
	re := regexp.MustCompile(`Note-ID:\s*(note-\d+)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// convertNoteTagsToLabels converts note tags and category to GitHub labels
func (r *RemoteService) convertNoteTagsToLabels(tags []string, category string) []string {
	var labels []string

	// Add category as a label with prefix
	if category != "" && category != "general" {
		labels = append(labels, fmt.Sprintf("category:%s", category))
	}

	// Add note-managed label to identify synced issues
	labels = append(labels, "note-managed")

	// Add tags as labels with tag prefix to avoid conflicts
	for _, tag := range tags {
		if tag != "" {
			labels = append(labels, fmt.Sprintf("tag:%s", tag))
		}
	}

	return labels
}

// convertLabelsToNoteTags converts GitHub labels to note tags and category
func (r *RemoteService) convertLabelsToNoteTags(labels []*github.Label) ([]string, string) {
	var tags []string
	category := "general"

	for _, label := range labels {
		labelName := label.GetName()

		// Check if it's a category label
		if strings.HasPrefix(labelName, "category:") {
			category = strings.TrimPrefix(labelName, "category:")
			continue
		}

		// Check if it's a tag label
		if strings.HasPrefix(labelName, "tag:") {
			tag := strings.TrimPrefix(labelName, "tag:")
			tags = append(tags, tag)
			continue
		}

		// Skip management labels
		if labelName == "note-managed" {
			continue
		}

		// Add other labels as tags (for backward compatibility)
		tags = append(tags, labelName)
	}

	return tags, category
}

// getRepoOwnerAndName extracts repository owner and name from git remote
func getRepoOwnerAndName() (string, string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", err
	}

	url := strings.TrimSpace(string(out))

	// Support SSH and HTTPS URLs
	// Example SSH: git@github.com:owner/repo.git
	// Example HTTPS: https://github.com/owner/repo.git
	re := regexp.MustCompile(`[:/]([^/]+)/([^/]+?)(\.git)?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("cannot parse repo owner and name from URL: %s", url)
	}

	owner := matches[1]
	repo := matches[2]
	return owner, repo, nil
}
