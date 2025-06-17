package webhook

import "time"

// PullRequestEvent represents the generic structure for pull request events.
// It includes fields common to all PR events and specific payload structs.
type PullRequestEvent struct {
	EventKey         string      `json:"eventKey"`
	Date             string      `json:"date"` // Example: "2017-09-19T09:58:11+1000"
	Actor            User        `json:"actor"`
	PullRequest      PullRequest `json:"pullRequest"`
	PreviousFromHash string      `json:"previousFromHash,omitempty"` // Only for pr:from_ref_updated
}

// PullRequest represents the detailed pull request object within the webhook payload.
type PullRequest struct {
	ID           int64         `json:"id"`
	Version      int           `json:"version"`
	Title        string        `json:"title"`
	State        string        `json:"state"` // "OPEN", "DECLINED", "MERGED"
	Open         bool          `json:"open"`
	Closed       bool          `json:"closed"`
	Draft        bool          `json:"draft"`
	CreatedDate  int64         `json:"createdDate"` // Unix timestamp in milliseconds
	UpdatedDate  int64         `json:"updatedDate"` // Unix timestamp in milliseconds
	FromRef      Ref           `json:"fromRef"`
	ToRef        Ref           `json:"toRef"`
	Locked       bool          `json:"locked"`
	Author       Participant   `json:"author"`
	Reviewers    []Participant `json:"reviewers"`
	Participants []Participant `json:"participants"`
}

// Ref represents a Git reference (branch).
type Ref struct {
	ID           string     `json:"id"`
	DisplayID    string     `json:"displayId"`
	LatestCommit string     `json:"latestCommit"`
	Repository   Repository `json:"repository"`
}

// Repository represents a Bitbucket repository.
type Repository struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
}

// Project represents a Bitbucket project.
type Project struct {
	Key    string `json:"key"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
}

// Participant represents a user involved in a pull request (author, reviewer).
type Participant struct {
	User     User   `json:"user"`
	Role     string `json:"role"` // e.g., "AUTHOR", "REVIEWER", "PARTICIPANT"
	Approved bool   `json:"approved"`
	Status   string `json:"status"` // e.g., "UNAPPROVED", "APPROVED"
}

// User represents a Bitbucket user.
type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

// GetCreatedDateAsTime converts the Unix millisecond timestamp to a time.Time object.
func (pr *PullRequest) GetCreatedDateAsTime() time.Time {
	return time.Unix(0, pr.CreatedDate*int64(time.Millisecond))
}

// GetUpdatedDateAsTime converts the Unix millisecond timestamp to a time.Time object.
func (pr *PullRequest) GetUpdatedDateAsTime() time.Time {
	return time.Unix(0, pr.UpdatedDate*int64(time.Millisecond))
}
