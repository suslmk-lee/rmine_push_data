package database

import (
	"database/sql"
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
        select j.id, i.id as job_id, i.tracker_id, i.project_id, i.subject, i.description, i.due_date , i.status_id, i.assigned_to_id, 
		i.created_on , i.updated_on , i.start_date , i.done_ratio, i.estimated_hours, i.priority_id, i.author_id, j.user_id as commentor_id,
		i.root_id, j.notes, jd.property, jd.prop_key, jd.old_value, jd.value 
		  from bitnami_redmine.issues i
		  left outer join bitnami_redmine.journals j
			 on i.id = j.journalized_id 
		  left outer join bitnami_redmine.journal_details jd
			 on j.id = jd.journal_id 
		   where j.id is not null  
           AND j.created_on > ?
        ORDER BY j.created_on DESC`

	rows, err := db.Query(query, formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []model.Issue
	for rows.Next() {
		var issue model.Issue
		var estimatedHours sql.NullFloat64
		var assignedToId sql.NullInt32
		var dueDate sql.NullTime
		var Property, PropKey, oldValue, value sql.NullString
		if err := rows.Scan(
			&issue.ID, &issue.JobID, &issue.TrackerID, &issue.ProjectID, &issue.Subject, &issue.Description, &dueDate, &issue.StatusID, &assignedToId,
			&issue.CreatedOn, &issue.UpdatedOn, &issue.StartDate, &issue.DoneRatio, &estimatedHours, &issue.PriorityID, &issue.AuthorID, &issue.CommentorID,
			&issue.RootID, &issue.Notes, &Property, &PropKey, &oldValue, &value,
		); err != nil {
			return nil, err
		}
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
		if assignedToId.Valid {
			issue.AssignedToID = assignedToId.Int32
		}
		if Property.Valid {
			issue.Property = Property.String
		}
		if PropKey.Valid {
			issue.PropKey = PropKey.String
		}
		if oldValue.Valid {
			issue.OldValue = oldValue.String
		}
		if value.Valid {
			issue.Value = value.String
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
