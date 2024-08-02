package model

import (
	"database/sql"
	"time"
)

// Issue represents an issue in the Redmine system
type Issue struct {
	ID             int       `json:"id"`
	JobID          int       `json:"job_id"`
	Status         string    `json:"status"`
	Assignee       string    `json:"assignee"`
	StartDate      time.Time `json:"start_date"`
	DueDate        time.Time `json:"due_date"`
	DoneRatio      int       `json:"done_ratio"`
	EstimatedHours float64   `json:"estimated_hours"`
	Priority       string    `json:"priority"`
	Author         string    `json:"author"`
	Subject        string    `json:"subject"`
	Description    string    `json:"description"`
	Commentor      string    `json:"commentor"`
	Notes          string    `json:"notes"`
	CreatedOn      time.Time `json:"created_on"`
}

type IssueDetail struct {
	ID             int       `json:"id"`
	TrackerId      int       `json:"job_id"`
	ProjectId      string    `json:"status"`
	Subject        string    `json:"subject"`
	Description    string    `json:"description"`
	DueData        time.Time `json:"due_date"`
	CategoryId     int       `json:"category_id"`
	StatusId       int       `json:"status_id"`
	AssignedToId   int       `json:"assigned_to_id"`
	PriorityId     int       `json:"priority_id"`
	FixedVersionId int       `json:"fixed_version_id"`
	AuthorId       int       `json:"author_id"`
	LockVersion    int       `json:"lock_version"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	StartDate      time.Time `json:"start_date"`
	DoneRatio      int       `json:"done_ratio"`
	EstimatedHours float64   `json:"estimated_hours"`
	ParentId       int       `json:"parent_id"`
	RootId         int       `json:"root_id"`
	IsPrivate      int       `json:"is_private"`
	ClosedOn       time.Time `json:"closed_on"`
}

type Message struct {
	ID          int           `json:"id"`
	BoardID     int           `json:"board_id"`
	ParentID    sql.NullInt64 `json:"parent_id"` // Use sql.NullInt64 to handle possible NULL values
	Subject     string        `json:"subject"`
	Content     string        `json:"content"`
	AuthorID    int           `json:"author_id"`
	LastReplyID sql.NullInt64 `json:"last_reply_id"` // Use sql.NullInt64 to handle possible NULL values
	CreatedOn   time.Time     `json:"created_on"`
	UpdatedOn   time.Time     `json:"updated_on"`
	Locked      bool          `json:"locked"`
	Sticky      bool          `json:"sticky"`
}

type JournalDetail struct {
	ID        int    `json:"id"`
	JournalID int    `json:"journal_id"`
	Property  string `json:"property"`
	PropKey   string `json:"prop_key"`
	OldValue  string `json:"old_value"`
	Value     string `json:"value"`
}

type User struct {
	ID               int           `json:"id"`
	Login            string        `json:"login"`
	HashedPassword   string        `json:"hashed_password"`
	FirstName        string        `json:"firstname"`
	LastName         string        `json:"lastname"`
	Admin            bool          `json:"admin"`
	Status           int           `json:"status"`
	LastLoginOn      sql.NullTime  `json:"last_login_on"` // Use sql.NullTime to handle possible NULL values
	Language         string        `json:"language"`
	AuthSourceID     sql.NullInt64 `json:"auth_source_id"` // Use sql.NullInt64 to handle possible NULL values
	CreatedOn        time.Time     `json:"created_on"`
	UpdatedOn        time.Time     `json:"updated_on"`
	Type             string        `json:"type"`
	MailNotification string        `json:"mail_notification"`
	Salt             string        `json:"salt"`
	MustChangePasswd bool          `json:"must_change_passwd"`
	PasswdChangedOn  sql.NullTime  `json:"passwd_changed_on"` // Use sql.NullTime to handle possible NULL values
}
