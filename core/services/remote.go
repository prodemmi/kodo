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
	noteService *NoteService
}

func NewRemoteManager(logger *zap.Logger, settings *SettingsService, noteService *NoteService) *RemoteService {
	return &RemoteService{
		logger:      logger,
		settings:    settings,
		noteService: noteService,
	}
}

type SyncResult struct {
	NotesCreated  int
	NotesUpdated  int
	IssuesCreated int
	IssuesUpdated int
	IssuesClosed  int
	Conflicts     []SyncConflict
	Errors        []error
}

type SyncConflict struct {
	Type        string
	NoteID      *int
	IssueNumber *int
	Description string
	LocalData   interface{}
	RemoteData  interface{}
}

func (r *RemoteService) SyncIssuesWithNotes() (*SyncResult, error) {
	result := &SyncResult{}

	settings := r.settings.LoadSettings()
	githubToken := settings.GithubAuth.Token

	if !settings.CodeScanSettings.SyncEnabled {
		return nil, fmt.Errorf("sync not enabled")
	}

	if githubToken == "" {
		return nil, fmt.Errorf("GitHub token not configured")
	}

	owner, repo, err := getRepoOwnerAndName()
	if err != nil {
		return nil, fmt.Errorf("failed to get repo info: %v", err)
	}

	r.logger.Info("Starting GitHub sync",
		zap.String("owner", owner),
		zap.String("repo", repo))

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	localNotes, err := r.noteService.GetNotes("", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to load local notes: %v", err)
	}

	allIssues, err := r.getAllGitHubIssues(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to load GitHub issues: %v", err)
	}

	r.logger.Info("Loaded data for sync",
		zap.Int("local_notes", len(localNotes)),
		zap.Int("github_issues", len(allIssues)))

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

func (r *RemoteService) getAllGitHubIssues(ctx context.Context, client *github.Client, owner, repo string) ([]*github.Issue, error) {
	var allIssues []*github.Issue

	openIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, err
	}
	allIssues = append(allIssues, openIssues...)

	closedIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "closed",
		Since:       time.Now().AddDate(0, -3, 0),
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, err
	}
	allIssues = append(allIssues, closedIssues...)

	return allIssues, nil
}

func (r *RemoteService) performBidirectionalSync(ctx context.Context, client *github.Client, owner, repo string, localNotes []entities.Note, githubIssues []*github.Issue, result *SyncResult) error {

	notesByGitHubID := make(map[int]*entities.Note)
	notesByTitle := make(map[string]*entities.Note)
	issuesByNumber := make(map[int]*github.Issue)
	issuesByNoteID := make(map[string]*github.Issue)
	processedIssues := make(map[int]bool)

	for i := range localNotes {
		note := &localNotes[i]

		if note.GitHubIssueNumber != nil {
			notesByGitHubID[*note.GitHubIssueNumber] = note
		}

		notesByTitle[note.Title] = note
	}

	for _, issue := range githubIssues {
		issuesByNumber[issue.GetNumber()] = issue

		noteID := r.extractNoteIDFromIssue(issue)
		if noteID != "" {
			issuesByNoteID[noteID] = issue
		}
	}

	r.logger.Info("Sync mappings",
		zap.Int("notes_by_github_id", len(notesByGitHubID)),
		zap.Int("notes_by_title", len(notesByTitle)),
		zap.Int("issues_by_number", len(issuesByNumber)),
		zap.Int("issues_by_note_id", len(issuesByNoteID)))

	err := r.syncExistingPairs(ctx, client, owner, repo, notesByGitHubID, issuesByNumber, processedIssues, result)
	if err != nil {
		return err
	}

	err = r.createIssuesForNewNotes(ctx, client, owner, repo, localNotes, notesByGitHubID, issuesByNoteID, processedIssues, result)
	if err != nil {
		return err
	}

	err = r.createNotesForNewIssues(ctx, githubIssues, notesByGitHubID, issuesByNoteID, processedIssues, result)
	if err != nil {
		return err
	}

	return nil
}

func (r *RemoteService) syncExistingPairs(ctx context.Context, client *github.Client, owner, repo string, notesByGitHubID map[int]*entities.Note, issuesByNumber map[int]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for issueNumber, note := range notesByGitHubID {
		issue, exists := issuesByNumber[issueNumber]
		if !exists {

			r.logger.Warn("GitHub issue no longer exists", zap.Int("issue_number", issueNumber), zap.Int("note_id", note.ID))

			note.GitHubIssueNumber = nil
			note.GitHubIssueURL = nil
			_, err := r.noteService.UpdateNoteWithHistory(note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to update note %d: %v", note.ID, err))
			} else {
				result.NotesUpdated++
			}
			continue
		}

		processedIssues[issueNumber] = true

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

func (r *RemoteService) createIssuesForNewNotes(ctx context.Context, client *github.Client, owner, repo string, localNotes []entities.Note, notesByGitHubID map[int]*entities.Note, issuesByNoteID map[string]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for _, note := range localNotes {

		if note.GitHubIssueNumber != nil {
			continue
		}

		noteIDStr := fmt.Sprintf("note-%d", note.ID)
		if existingIssue, exists := issuesByNoteID[noteIDStr]; exists {

			err := r.linkNoteToIssue(&note, existingIssue)
			if err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}

			processedIssues[existingIssue.GetNumber()] = true
			result.NotesUpdated++
			continue
		}

		err := r.createGitHubIssueFromNote(ctx, client, owner, repo, &note, result)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	return nil
}

func (r *RemoteService) createNotesForNewIssues(ctx context.Context, githubIssues []*github.Issue, notesByGitHubID map[int]*entities.Note, issuesByNoteID map[string]*github.Issue, processedIssues map[int]bool, result *SyncResult) error {
	for _, issue := range githubIssues {
		issueNumber := issue.GetNumber()

		if processedIssues[issueNumber] {
			continue
		}

		if issue.PullRequestLinks != nil {
			continue
		}

		if !r.hasNoteManagedLabel(issue) {

			err := r.createNoteFromGitHubIssue(issue, result)
			if err != nil {
				result.Errors = append(result.Errors, err)
			}
		}
	}

	return nil
}

func (r *RemoteService) hasNoteManagedLabel(issue *github.Issue) bool {
	for _, label := range issue.Labels {
		if label.GetName() == "note-managed" {
			return true
		}
	}
	return false
}

func (r *RemoteService) createNoteFromGitHubIssue(issue *github.Issue, result *SyncResult) error {

	content := r.extractNoteContentFromIssue(issue)

	tags, category := r.convertLabelsToNoteTags(issue.Labels)

	author := "unknown"

	if issue.User != nil && issue.User.Login != nil {
		author = fmt.Sprintf("%s %s", issue.User.GetLogin(), issue.User.GetEmail())
	}

	issueNumber := issue.GetNumber()
	now := time.Now().UTC()
	status := issue.GetState()

	createdNote, err := r.noteService.CreateNoteWithHistory(
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

	note := r.noteService.getNoteByID(createdNote.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve created note: %v", err)
	}

	note.GitHubIssueNumber = &issueNumber
	note.GitHubIssueURL = issue.HTMLURL
	note.GitHubLastSync = &now
	note.Status = &status

	_, err = r.noteService.UpdateNoteWithHistory(
		note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update note with GitHub info: %v", err)
	}

	result.NotesCreated++
	return nil
}

func (r *RemoteService) needsSyncCheck(note *entities.Note, issue *github.Issue) (bool, string) {
	issueUpdated := issue.GetUpdatedAt().UTC()
	noteUpdated := note.UpdatedAt.UTC()

	var lastSyncTime time.Time
	if note.GitHubLastSync != nil {
		lastSyncTime = note.GitHubLastSync.UTC()
	}

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

	if noteUpdated.After(issueUpdated) {
		return true, "note_to_issue"
	}
	if issueUpdated.After(noteUpdated) {
		return true, "issue_to_note"
	}

	return true, "conflict"
}

func (r *RemoteService) updateIssueFromNote(ctx context.Context, client *github.Client, owner, repo string, note *entities.Note, issue *github.Issue, result *SyncResult) error {

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

	r.logger.Info("Updating issue", zap.Any("issue_request", issueRequest))

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

	now := time.Now().UTC()
	note.GitHubLastSync = &now
	note.GitHubIssueURL = updatedIssue.HTMLURL

	_, err = r.noteService.UpdateNoteWithHistory(note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID)
	if err != nil {
		return fmt.Errorf("failed to update note sync info: %v", err)
	}

	result.IssuesUpdated++
	return nil
}

func (r *RemoteService) updateNoteFromIssue(issue *github.Issue, note *entities.Note, result *SyncResult) error {

	content := r.extractNoteContentFromIssue(issue)

	tags, category := r.convertLabelsToNoteTags(issue.Labels)

	status := issue.GetState()

	r.logger.Info("Updating note",
		zap.Int("note_id", note.ID),
		zap.String("title", issue.GetTitle()),
		zap.String("content", content),
		zap.Strings("tags", tags),
		zap.String("category", category))

	now := time.Now().UTC()
	note.GitHubLastSync = &now
	note.Status = &status
	note.Tags = tags
	note.Category = category

	_, err := r.noteService.UpdateNoteWithHistory(note.ID, issue.GetTitle(), content, tags, category, note.Pinned, note.FolderID)
	if err != nil {
		return fmt.Errorf("failed to update note from issue: %v", err)
	}

	result.NotesUpdated++
	return nil
}

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

	err = r.linkNoteToIssue(note, createdIssue)
	if err != nil {
		return fmt.Errorf("failed to link note to created issue: %v", err)
	}

	result.IssuesCreated++
	return nil
}

func (r *RemoteService) linkNoteToIssue(note *entities.Note, issue *github.Issue) error {
	now := time.Now().UTC()
	issueNumber := issue.GetNumber()

	note.GitHubIssueNumber = &issueNumber
	note.GitHubIssueURL = issue.HTMLURL
	note.GitHubLastSync = &now

	status := issue.GetState()
	note.Status = &status

	_, err := r.noteService.UpdateNoteWithHistory(
		note.ID, note.Title, note.Content, note.Tags, note.Category, note.Pinned, note.FolderID,
	)
	return err
}

func (r *RemoteService) formatIssueBodyFromNote(note *entities.Note) string {
	var body strings.Builder

	body.WriteString(note.Content)
	body.WriteString("\n\n")

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

func (r *RemoteService) extractNoteContentFromIssue(issue *github.Issue) string {
	body := issue.GetBody()

	metadataStart := strings.Index(body, "---\n<!-- Note Metadata")
	if metadataStart != -1 {
		return strings.TrimSpace(body[:metadataStart])
	}

	return body
}

func (r *RemoteService) extractNoteIDFromIssue(issue *github.Issue) string {
	body := issue.GetBody()

	re := regexp.MustCompile(`Note-ID:\s*(note-\d+)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func (r *RemoteService) convertNoteTagsToLabels(tags []string, category string) []string {
	var labels []string

	if category != "" && category != "general" {
		labels = append(labels, fmt.Sprintf("category:%s", category))
	}

	labels = append(labels, "note-managed")

	for _, tag := range tags {
		if tag != "" {
			labels = append(labels, fmt.Sprintf("tag:%s", tag))
		}
	}

	return labels
}

func (r *RemoteService) convertLabelsToNoteTags(labels []*github.Label) ([]string, string) {
	var tags []string
	category := "general"

	for _, label := range labels {
		labelName := label.GetName()

		if strings.HasPrefix(labelName, "category:") {
			category = strings.TrimPrefix(labelName, "category:")
			continue
		}

		if strings.HasPrefix(labelName, "tag:") {
			tag := strings.TrimPrefix(labelName, "tag:")
			tags = append(tags, tag)
			continue
		}

		if labelName == "note-managed" {
			continue
		}

		tags = append(tags, labelName)
	}

	return tags, category
}

func getRepoOwnerAndName() (string, string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", err
	}

	url := strings.TrimSpace(string(out))

	re := regexp.MustCompile(`[:/]([^/]+)/([^/]+?)(\.git)?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("cannot parse repo owner and name from URL: %s", url)
	}

	owner := matches[1]
	repo := matches[2]
	return owner, repo, nil
}
