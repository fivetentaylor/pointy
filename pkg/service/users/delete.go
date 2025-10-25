package users

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/query"
	dynamodb_storage "github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
	"gorm.io/gorm"
)

type UserDeletionService struct {
	DB       *gorm.DB
	DynamoDB *dynamodb_storage.DB
	S3       *s3.S3
}

type UserDeletionResult struct {
	UserID                    string
	DeletedPostgresqlRecords  int
	DeletedDynamoDBRecords    int
	DeletedS3Objects          int
	Errors                    []string
	CompletedAt               time.Time
}

func NewUserDeletionService(db *gorm.DB, dynamodb *dynamodb_storage.DB, s3Client *s3.S3) *UserDeletionService {
	return &UserDeletionService{
		DB:       db,
		DynamoDB: dynamodb,
		S3:       s3Client,
	}
}

// DeleteUser performs a comprehensive deletion of all user data across PostgreSQL, DynamoDB, and S3
func (uds *UserDeletionService) DeleteUser(ctx context.Context, userID string) (*UserDeletionResult, error) {
	log.Infof("Starting comprehensive user deletion for userID: %s", userID)

	result := &UserDeletionResult{
		UserID:  userID,
		Errors:  []string{},
	}

	// Step 1: Delete PostgreSQL records (in proper order due to foreign key constraints)
	pgCount, err := uds.deletePostgreSQLData(ctx, userID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("PostgreSQL deletion error: %v", err))
		log.Errorf("PostgreSQL deletion failed: %v", err)
	} else {
		result.DeletedPostgresqlRecords = pgCount
		log.Infof("Deleted %d PostgreSQL records for user %s", pgCount, userID)
	}

	// Step 2: Delete DynamoDB records
	dynamoCount, err := uds.deleteDynamoDBData(ctx, userID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("DynamoDB deletion error: %v", err))
		log.Errorf("DynamoDB deletion failed: %v", err)
	} else {
		result.DeletedDynamoDBRecords = dynamoCount
		log.Infof("Deleted %d DynamoDB records for user %s", dynamoCount, userID)
	}

	// Step 3: Delete S3 objects
	s3Count, err := uds.deleteS3Data(ctx, userID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("S3 deletion error: %v", err))
		log.Errorf("S3 deletion failed: %v", err)
	} else {
		result.DeletedS3Objects = s3Count
		log.Infof("Deleted %d S3 objects for user %s", s3Count, userID)
	}

	result.CompletedAt = time.Now()

	if len(result.Errors) > 0 {
		log.Warnf("User deletion completed with %d errors for user %s", len(result.Errors), userID)
		return result, fmt.Errorf("user deletion completed with errors: %v", result.Errors)
	}

	log.Infof("Successfully completed user deletion for userID: %s", userID)
	return result, nil
}

// deletePostgreSQLData removes all user data from PostgreSQL tables
func (uds *UserDeletionService) deletePostgreSQLData(ctx context.Context, userID string) (int, error) {
	q := query.Use(uds.DB)
	totalDeleted := 0

	// Delete in order to handle foreign key constraints
	tables := []struct {
		name   string
		delete func() (int64, error)
	}{
		// Delete dependent records first
		{"comments", func() (int64, error) {
			result, err := q.Comment.Unscoped().Where(q.Comment.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"document_attachments", func() (int64, error) {
			result, err := q.DocumentAttachment.Where(q.DocumentAttachment.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"document_versions (created_by)", func() (int64, error) {
			result, err := q.DocumentVersion.Where(q.DocumentVersion.CreatedBy.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"document_versions (updated_by)", func() (int64, error) {
			result, err := q.DocumentVersion.Where(q.DocumentVersion.UpdatedBy.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"document_access", func() (int64, error) {
			result, err := q.DocumentAccess.Where(q.DocumentAccess.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"shared_document_links", func() (int64, error) {
			result, err := q.SharedDocumentLink.Where(q.SharedDocumentLink.InviterID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"author_ids", func() (int64, error) {
			result, err := q.AuthorID.Where(q.AuthorID.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"one_time_access_tokens", func() (int64, error) {
			result, err := q.OneTimeAccessToken.Where(q.OneTimeAccessToken.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"payment_history", func() (int64, error) {
			result, err := q.PaymentHistory.Where(q.PaymentHistory.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		{"user_subscriptions", func() (int64, error) {
			result, err := q.UserSubscription.Where(q.UserSubscription.UserID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
		// Finally delete the main user record
		{"users", func() (int64, error) {
			result, err := q.User.Where(q.User.ID.Eq(userID)).Delete()
			return result.RowsAffected, err
		}},
	}

	for _, table := range tables {
		deleted, err := table.delete()
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to delete from %s: %w", table.name, err)
		}
		totalDeleted += int(deleted)
		if deleted > 0 {
			log.Debugf("Deleted %d records from %s for user %s", deleted, table.name, userID)
		}
	}

	return totalDeleted, nil
}

// deleteDynamoDBData removes all user data from DynamoDB
func (uds *UserDeletionService) deleteDynamoDBData(ctx context.Context, userID string) (int, error) {
	totalDeleted := 0

	// Delete messages where user is the author
	msgCount, err := uds.deleteUserMessages(userID)
	if err != nil {
		return totalDeleted, fmt.Errorf("failed to delete user messages: %w", err)
	}
	totalDeleted += msgCount

	// Delete timeline events for the user
	timelineCount, err := uds.deleteUserTimelineEvents(userID)
	if err != nil {
		return totalDeleted, fmt.Errorf("failed to delete user timeline events: %w", err)
	}
	totalDeleted += timelineCount

	// Delete notifications for the user
	notifCount, err := uds.deleteUserNotifications(userID)
	if err != nil {
		return totalDeleted, fmt.Errorf("failed to delete user notifications: %w", err)
	}
	totalDeleted += notifCount

	// Delete user preferences
	prefCount, err := uds.deleteUserPreferences(userID)
	if err != nil {
		return totalDeleted, fmt.Errorf("failed to delete user preferences: %w", err)
	}
	totalDeleted += prefCount

	return totalDeleted, nil
}

// deleteUserMessages removes all messages created by the user
func (uds *UserDeletionService) deleteUserMessages(userID string) (int, error) {
	// This is a complex operation as we need to scan for messages by userID
	// In a production environment, you might want to add a GSI on userID for efficiency
	
	// For now, we'll need to scan the table and filter by userID
	// This is inefficient but necessary without a GSI on userID
	
	log.Warnf("Scanning DynamoDB table for user messages - this may be slow without GSI on userID")
	
	input := &dynamodb.ScanInput{
		TableName:        &uds.DynamoDB.TableName,
		FilterExpression: aws.String("userID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {S: aws.String(userID)},
		},
	}

	deleted := 0
	err := uds.DynamoDB.Client.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range page.Items {
			// Extract PK and SK for deletion
			if pk, exists := item["PK"]; exists && pk.S != nil {
				if sk, exists := item["SK"]; exists && sk.S != nil {
					deleteInput := &dynamodb.DeleteItemInput{
						TableName: &uds.DynamoDB.TableName,
						Key: map[string]*dynamodb.AttributeValue{
							"PK": pk,
							"SK": sk,
						},
					}
					
					_, err := uds.DynamoDB.Client.DeleteItem(deleteInput)
					if err != nil {
						log.Errorf("Failed to delete DynamoDB item %s/%s: %v", *pk.S, *sk.S, err)
						return false
					}
					deleted++
				}
			}
		}
		return true
	})

	return deleted, err
}

// deleteUserTimelineEvents removes all timeline events for the user
func (uds *UserDeletionService) deleteUserTimelineEvents(userID string) (int, error) {
	// Similar to messages, this requires a scan unless there's a GSI on userID
	input := &dynamodb.ScanInput{
		TableName:        &uds.DynamoDB.TableName,
		FilterExpression: aws.String("userID = :userID AND begins_with(SK, :timePrefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID":    {S: aws.String(userID)},
			":timePrefix": {S: aws.String("time#")},
		},
	}

	deleted := 0
	err := uds.DynamoDB.Client.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range page.Items {
			if pk, exists := item["PK"]; exists && pk.S != nil {
				if sk, exists := item["SK"]; exists && sk.S != nil {
					deleteInput := &dynamodb.DeleteItemInput{
						TableName: &uds.DynamoDB.TableName,
						Key: map[string]*dynamodb.AttributeValue{
							"PK": pk,
							"SK": sk,
						},
					}
					
					_, err := uds.DynamoDB.Client.DeleteItem(deleteInput)
					if err != nil {
						log.Errorf("Failed to delete timeline event %s/%s: %v", *pk.S, *sk.S, err)
						return false
					}
					deleted++
				}
			}
		}
		return true
	})

	return deleted, err
}

// deleteUserNotifications removes all notifications for the user
func (uds *UserDeletionService) deleteUserNotifications(userID string) (int, error) {
	// Notifications use PK pattern: notif#userID
	pk := fmt.Sprintf("notif#%s", userID)
	
	input := &dynamodb.QueryInput{
		TableName:              &uds.DynamoDB.TableName,
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(pk)},
		},
	}

	deleted := 0
	err := uds.DynamoDB.Client.QueryPages(input, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		for _, item := range page.Items {
			if sk, exists := item["SK"]; exists && sk.S != nil {
				deleteInput := &dynamodb.DeleteItemInput{
					TableName: &uds.DynamoDB.TableName,
					Key: map[string]*dynamodb.AttributeValue{
						"PK": {S: aws.String(pk)},
						"SK": sk,
					},
				}
				
				_, err := uds.DynamoDB.Client.DeleteItem(deleteInput)
				if err != nil {
					log.Errorf("Failed to delete notification %s/%s: %v", pk, *sk.S, err)
					return false
				}
				deleted++
			}
		}
		return true
	})

	return deleted, err
}

// deleteUserPreferences removes user preferences and document preferences
func (uds *UserDeletionService) deleteUserPreferences(userID string) (int, error) {
	deleted := 0

	// Delete user preferences (PK: userPref#userID)
	userPrefPK := fmt.Sprintf("userPref#%s", userID)
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: &uds.DynamoDB.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(userPrefPK)},
			"SK": {S: aws.String(userID)},
		},
		ReturnValues: aws.String("ALL_OLD"), // This will return the deleted item if it existed
	}
	
	result, err := uds.DynamoDB.Client.DeleteItem(deleteInput)
	if err != nil {
		log.Errorf("Failed to delete user preferences: %v", err)
	} else if result.Attributes != nil {
		// Only count as deleted if there were actually attributes returned (item existed)
		deleted++
	}

	// Delete document notification preferences (PK: docNotifPref#userID)
	docPrefPK := fmt.Sprintf("docNotifPref#%s", userID)
	
	queryInput := &dynamodb.QueryInput{
		TableName:              &uds.DynamoDB.TableName,
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(docPrefPK)},
		},
	}

	err = uds.DynamoDB.Client.QueryPages(queryInput, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		for _, item := range page.Items {
			if sk, exists := item["SK"]; exists && sk.S != nil {
				deleteInput := &dynamodb.DeleteItemInput{
					TableName: &uds.DynamoDB.TableName,
					Key: map[string]*dynamodb.AttributeValue{
						"PK": {S: aws.String(docPrefPK)},
						"SK": sk,
					},
				}
				
				_, err := uds.DynamoDB.Client.DeleteItem(deleteInput)
				if err != nil {
					log.Errorf("Failed to delete document preference %s/%s: %v", docPrefPK, *sk.S, err)
					return false
				}
				deleted++
			}
		}
		return true
	})

	return deleted, err
}

// deleteS3Data removes all S3 objects associated with the user
func (uds *UserDeletionService) deleteS3Data(ctx context.Context, userID string) (int, error) {
	totalDeleted := 0

	// Delete user avatar
	avatarKey := fmt.Sprintf(constants.UserAvatarKeyFormat, userID)
	exists, err := uds.S3.Exists(uds.S3.ImagesBucket, avatarKey)
	if err != nil {
		log.Errorf("Failed to check if avatar exists: %v", err)
	} else if exists {
		err = uds.S3.DeleteAll(uds.S3.ImagesBucket, avatarKey)
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to delete user avatar: %w", err)
		}
		totalDeleted++
		log.Debugf("Deleted user avatar: %s", avatarKey)
	}

	// Delete user's document attachments
	// We need to find all attachments uploaded by this user from PostgreSQL first
	q := query.Use(uds.DB)
	attachments, err := q.DocumentAttachment.Where(q.DocumentAttachment.UserID.Eq(userID)).Find()
	if err != nil {
		return totalDeleted, fmt.Errorf("failed to query user attachments: %w", err)
	}

	for _, attachment := range attachments {
		// Delete original attachment
		originalKey := fmt.Sprintf(constants.DocumentAttachmentOriginalKey, attachment.S3ID)
		err = uds.S3.DeleteAll(uds.S3.Bucket, originalKey)
		if err != nil {
			log.Errorf("Failed to delete attachment original %s: %v", originalKey, err)
		} else {
			totalDeleted++
		}

		// Delete attachment file
		fileKey := fmt.Sprintf(constants.DocumentAttachmentFileKey, attachment.S3ID, attachment.Filename)
		err = uds.S3.DeleteAll(uds.S3.Bucket, fileKey)
		if err != nil {
			log.Errorf("Failed to delete attachment file %s: %v", fileKey, err)
		} else {
			totalDeleted++
		}
	}

	// Note: Document images are tied to documents, not users directly
	// They should be cleaned up when documents are deleted
	// If you need to delete user-specific document images, you would need
	// additional logic to find which images were created by this user

	return totalDeleted, nil
}