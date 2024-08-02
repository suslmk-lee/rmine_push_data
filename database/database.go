package database

import (
	"database/sql"
	"fmt"
	"rmine_push_data/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB connects to the MySQL database
func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FetchMessages(db *sql.DB, lastChecked time.Time) ([]model.Message, error) {
	query := `
			select m.id, m.board_id, m.parent_id, m.subject, replace(m.content, '"', ''), 
			       m.author_id, m.last_reply_id, m.created_on, m.updated_on, m.locked, m.sticky
  from bitnami_redmine.messages m, bitnami_redmine.users a, bitnami_redmine.email_addresses ea 
 where m.author_id = a.id 
   and a.id = ea.user_id
   and ea.is_default = 1
   and m.updated_on > ?
   and m.parent_id is null 
order by m.updated_on desc`

	formattedTime := lastChecked.Format("2006-01-02 15:04:05")

	rows, err := db.Query(query, formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(
			&message.ID, &message.BoardID, &message.ParentID, &message.Subject, &message.Content,
			&message.AuthorID, &message.LastReplyID, &message.CreatedOn, &message.UpdatedOn, &message.Locked, &message.Sticky,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func FetchJournalDetail(db *sql.DB, lastChecked time.Time) ([]model.JournalDetail, error) {
	query := `
			select j.id, j.journal_id, j.property, j.prop_key, j.old_value, j.value from bitnami_redmine.journal_details j`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journalDetails []model.JournalDetail
	for rows.Next() {
		var journalDetail model.JournalDetail
		if err := rows.Scan(
			&journalDetail.ID, &journalDetail.JournalID, &journalDetail.Property, &journalDetail.PropKey, &journalDetail.OldValue, &journalDetail.Value,
		); err != nil {
			return nil, err
		}
		journalDetails = append(journalDetails, journalDetail)
	}
	return journalDetails, nil
}

func FetchUsers(db *sql.DB, lastChecked time.Time) ([]model.User, error) {
	query := `
			select u.id, u.login, u.hashed_password, u.firstname, u.lastname, u.admin, u.status, u.last_login_on, u.language,  
       u.auth_source_id , u.created_on , u.updated_on , u.type, u.mail_notification , u.salt , u.must_change_passwd , u.passwd_changed_on
  from bitnami_redmine.users u
   where u.updated_on > ?
     and u.login != ''`

	formattedTime := lastChecked.Format("2006-01-02 15:04:05")

	rows, err := db.Query(query, formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID, &user.Login, &user.HashedPassword, &user.FirstName, &user.LastName, &user.Admin, &user.Status, &user.LastLoginOn,
			&user.Language, &user.AuthSourceID, &user.CreatedOn, &user.UpdatedOn, &user.Type, &user.MailNotification, &user.Salt,
			&user.MustChangePasswd, &user.PasswdChangedOn,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func FetchIssues(db *sql.DB, lastChecked time.Time) ([]model.IssueDetail, error) {
	formattedTime := lastChecked.Format("2006-01-02 15:04:05")
	query := `
			select i.id, i.tracker_id, i.project_id, i.subject, i.description, i.due_date , i.status_id, i.assigned_to_id, 
				i.created_on , i.updated_on , i.start_date , i.done_ratio, i.priority_id, i.author_id, 
				i.project_id , i.root_id 
			  from bitnami_redmine.issues i
			 where i.updated_on > ?
			order by i.updated_on desc`

	rows, err := db.Query(query, formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []model.IssueDetail
	for rows.Next() {
		var issue model.IssueDetail
		var dueDate sql.NullTime
		if err := rows.Scan(
			&issue.ID, &issue.TrackerId, &issue.ProjectId, &issue.Subject, &issue.Description, &issue.DueDate, &issue.StatusId, &issue.AssignedToId,
			&issue.CreatedOn, &issue.UpdatedOn, &issue.StartDate, &issue.DoneRatio, &issue.PriorityId, &issue.AuthorId,
			&issue.ProjectId, &issue.RootId,
		); err != nil {
			return nil, err
		}
		if dueDate.Valid {
			issue.DueDate = dueDate.Time
		} else {
			issue.DueDate = time.Time{}
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

// FetchNewIssues fetches new issues from the MySQL database
func FetchNewIssues(db *sql.DB, lastChecked time.Time) ([]model.Issue, error) {
	formattedTime := lastChecked.Format("2006-01-02 15:04:05")
	query := `
        SELECT j.id, i.id as 'job_id', is2.name, u.firstname, u.lastname, i.start_date, i.due_date, i.done_ratio, i.estimated_hours,
        (SELECT e.name FROM bitnami_redmine.enumerations e WHERE e.type = 'IssuePriority' AND i.priority_id = e.id) AS priority,
        (SELECT b.firstname FROM bitnami_redmine.users b WHERE i.author_id = b.id) AS author,
        i.subject, i.description,
        (SELECT b.firstname FROM bitnami_redmine.users b WHERE j.user_id = b.id) AS commentor,
        j.notes, j.created_on
        FROM bitnami_redmine.issues i
        JOIN bitnami_redmine.issue_statuses is2 ON i.status_id = is2.id
        JOIN bitnami_redmine.users u ON i.assigned_to_id = u.id
        JOIN bitnami_redmine.journals j ON i.id = j.journalized_id
        WHERE j.created_on > ?
        ORDER BY j.created_on DESC`

	rows, err := db.Query(query, formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []model.Issue
	for rows.Next() {
		var issue model.Issue
		var assigneeFirstName, assigneeLastName, commentorFirstName sql.NullString
		var estimatedHours sql.NullFloat64
		var dueDate sql.NullTime
		if err := rows.Scan(
			&issue.ID, &issue.JobID, &issue.Status, &assigneeFirstName, &assigneeLastName, &issue.StartDate, &dueDate,
			&issue.DoneRatio, &estimatedHours, &issue.Priority, &issue.Author, &issue.Subject,
			&issue.Description, &commentorFirstName, &issue.Notes, &issue.CreatedOn,
		); err != nil {
			return nil, err
		}
		issue.Assignee = fmt.Sprintf("%s %s", assigneeFirstName.String, assigneeLastName.String)
		issue.Commentor = commentorFirstName.String
		if estimatedHours.Valid {
			issue.EstimatedHours = estimatedHours.Float64
		} else {
			issue.EstimatedHours = 0
		}
		if dueDate.Valid {
			issue.DueDate = dueDate.Time
		} else {
			issue.DueDate = time.Time{}
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
