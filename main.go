package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"time"

	"rmine_push_data/action"
	"rmine_push_data/common"
	"rmine_push_data/database"
)

const (
	pollInterval = 10 * time.Second
)

var (
	bucketName string
	region     string
	dsn        string
	endpoint   string
	accessKey  string
	secretKey  string
)

func init() {
	region = common.ConfInfo["nhn.region"]
	bucketName = common.ConfInfo["nhn.storage.bucket.name"]
	dsn = common.ConfInfo["database.url"]
	endpoint = common.ConfInfo["nhn.storage.endpoint.url"]
	accessKey = common.ConfInfo["nhn.storage.accessKey"]
	secretKey = common.ConfInfo["nhn.storage.secretKey"]
}

// logEveryHour logs a message every hour at the start of the hour
func logEveryHour() {
	for {
		now := time.Now()
		// Calculate the duration until the next hour
		next := now.Truncate(time.Hour).Add(time.Hour)
		duration := next.Sub(now)

		// Sleep until the start of the next hour
		time.Sleep(duration)

		// Log the message
		log.Printf("Hourly log: Current time is %s", next.Format(time.RFC1123))
	}
}

func main() {

	// Start the hourly logging function in a new goroutine
	go logEveryHour()

	fmt.Println("Start to RMine Push Data Service....")

	// Ensure the keys are not empty
	if accessKey == "" || secretKey == "" {
		log.Fatalf("AccessKey or SecretKey is empty")
	}

	// Connect to MySQL database
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Load the last checked time from file
	lastChecked, err := common.LoadLastCheckedTime()
	if err != nil {
		lastChecked.IsZero()
		fmt.Printf("failed to load last checked time: %v", err)
	}
	if lastChecked.IsZero() {
		// If there is no last checked time, start from one week ago
		lastChecked = time.Now().Add(-7 * 24 * time.Hour)
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true)}, // Use path-style addressing for compatibility with custom endpoints
	)
	if err != nil {
		log.Fatalf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	for {
		// Fetch and process issues
		issues, err := database.FetchNewIssues(db, lastChecked)
		if err != nil {
			log.Printf("failed to fetch new issues: %v", err)
			continue
		}

		if issues != nil {
			err = action.ProcessIssues(s3Client, bucketName, issues)
			if err != nil {
				log.Printf("failed to process and upload issues: %v", err)
			}
			log.Printf("Finished processing %d issues", len(issues))
		}

		// Fetch and process messages
		messages, err := database.FetchMessages(db, lastChecked)
		if err != nil {
			log.Printf("failed to fetch messages: %v", err)
			continue
		}

		if messages != nil {
			err = action.ProcessMessages(s3Client, bucketName, messages)
			if err != nil {
				log.Printf("failed to process and upload messages: %v", err)
			}
			log.Printf("Finished processing %d messages", len(messages))
		}

		// Fetch and process users
		users, err := database.FetchUsers(db, lastChecked)
		if err != nil {
			log.Printf("failed to fetch users: %v", err)
			continue
		}

		if users != nil {
			err = action.ProcessUsers(s3Client, bucketName, users)
			if err != nil {
				log.Printf("failed to process and upload users: %v", err)
			}
			log.Printf("Finished processing %d users", len(users))
		}

		rawIssue, err := database.FetchIssues(db, lastChecked)
		if err != nil {
			log.Printf("failed to fetch Raw issues: %v", err)
			continue
		}

		if rawIssue != nil {
			err = action.ProcessRawIssues(s3Client, bucketName, rawIssue)
			if err != nil {
				log.Printf("failed to process and upload raw issues: %v", err)
			}
			log.Printf("Finished processing %d raw-issues", len(rawIssue))
		}

		// Update lastChecked time
		lastChecked = time.Now()
		err = common.SaveLastCheckedTime(lastChecked)
		if err != nil {
			log.Printf("failed to save last checked time: %v", err)
		}

		// Sleep for the poll interval
		time.Sleep(pollInterval)
	}
}
