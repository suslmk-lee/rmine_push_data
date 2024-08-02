package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"rmine_push_data/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

// CreateCloudEvent creates a CloudEvent for the given data
func CreateCloudEvent(source, eventType string, data interface{}) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetSource(source)
	event.SetType(eventType)
	event.SetTime(time.Now())

	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return event, err
	}
	return event, nil
}

// UploadToS3 uploads the given data to S3 with the specified key
func UploadToS3(s3Client *s3.S3, bucketName string, data []byte, key string) error {
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err := s3Client.PutObject(input)
	return err
}

// ProcessIssues processes a list of issues by creating CloudEvents and uploading them to S3
func ProcessIssues(s3Client *s3.S3, bucketName string, issues []model.Issue) error {
	for _, issue := range issues {
		event, err := CreateCloudEvent("redmine/issues", "com.example.issue", issue)
		if err != nil {
			log.Printf("failed to create CloudEvent: %v", err)
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal CloudEvent: %v", err)
			continue
		}

		key := fmt.Sprintf("rmine_push_data/issues/%d.json", issue.ID)
		err = UploadToS3(s3Client, bucketName, data, key)
		if err != nil {
			log.Printf("failed to upload data to S3: %v", err)
			continue
		}
	}
	return nil
}

// ProcessMessages processes a list of messages by creating CloudEvents and uploading them to S3
func ProcessMessages(s3Client *s3.S3, bucketName string, messages []model.Message) error {
	for _, message := range messages {
		event, err := CreateCloudEvent("redmine/messages", "com.example.message", message)
		if err != nil {
			log.Printf("failed to create CloudEvent: %v", err)
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal CloudEvent: %v", err)
			continue
		}

		key := fmt.Sprintf("rmine_push_data/messages/%d.json", message.ID)
		err = UploadToS3(s3Client, bucketName, data, key)
		if err != nil {
			log.Printf("failed to upload data to S3: %v", err)
			continue
		}
	}
	return nil
}

// ProcessJournalDetails processes a list of journal details by creating CloudEvents and uploading them to S3
func ProcessJournalDetails(s3Client *s3.S3, bucketName string, journalDetails []model.JournalDetail) error {
	for _, journalDetail := range journalDetails {
		event, err := CreateCloudEvent("redmine/journal_details", "com.example.journal_detail", journalDetail)
		if err != nil {
			log.Printf("failed to create CloudEvent: %v", err)
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal CloudEvent: %v", err)
			continue
		}

		key := fmt.Sprintf("rmine_push_data/journal_details/%d.json", journalDetail.ID)
		err = UploadToS3(s3Client, bucketName, data, key)
		if err != nil {
			log.Printf("failed to upload data to S3: %v", err)
			continue
		}
	}
	return nil
}

// ProcessUsers processes a list of users by creating CloudEvents and uploading them to S3
func ProcessUsers(s3Client *s3.S3, bucketName string, users []model.User) error {
	for _, user := range users {
		event, err := CreateCloudEvent("redmine/users", "com.example.user", user)
		if err != nil {
			log.Printf("failed to create CloudEvent: %v", err)
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal CloudEvent: %v", err)
			continue
		}

		key := fmt.Sprintf("rmine_push_data/users/%d.json", user.ID)
		err = UploadToS3(s3Client, bucketName, data, key)
		if err != nil {
			log.Printf("failed to upload data to S3: %v", err)
			continue
		}
	}
	return nil
}
