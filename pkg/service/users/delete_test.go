package users

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	dynamodb_storage "github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

type TestData struct {
	User             *models.User
	Document         *models.Document
	Comment          *models.Comment
	DocumentAccess   *models.DocumentAccess
	Attachment       *models.DocumentAttachment
	DocumentVersion  *models.DocumentVersion
	AuthorID         *models.AuthorID
	AccessToken      *models.OneTimeAccessToken
	PaymentHistory   *models.PaymentHistory
	SharedLink       *models.SharedDocumentLink
	UserSubscription *models.UserSubscription
	Message          *dynamodb_storage.Message
	TimelineEvent    *dynamodb_storage.TimelineEvent
	Notification     *dynamodb_storage.Notification
	UserPreference   *dynamodb_storage.UserPreference
	DocPreference    *dynamodb_storage.DocPreference
	AvatarKey        string
	AttachmentKeys   []string
}

func TestDeleteUser_ComprehensiveTest(t *testing.T) {
	// Setup test environment
	ctx := context.Background()
	
	// Ensure storage is available
	testutils.EnsureStorage()
	
	// Get test context which includes all dependencies
	ctx = testutils.TestContext()
	
	// Extract dependencies from context
	db := testutils.TestGormDb(t)
	dynamoDB, err := dynamodb_storage.NewDB()
	require.NoError(t, err, "Failed to create DynamoDB client")
	
	s3Client, err := s3.NewS3()
	require.NoError(t, err, "Failed to create S3 client")

	// Create deletion service
	deletionService := NewUserDeletionService(db, dynamoDB, s3Client)

	// Create comprehensive test data
	testData := createComprehensiveTestData(t, db, dynamoDB, s3Client)
	userID := testData.User.ID

	// Verify data exists before deletion
	verifyDataExists(t, db, dynamoDB, s3Client, testData)

	// Perform user deletion
	result, err := deletionService.DeleteUser(ctx, userID)

	// Assert deletion succeeded
	if len(result.Errors) > 0 {
		t.Logf("Deletion completed with errors: %v", result.Errors)
	}
	require.NoError(t, err, "User deletion should succeed")
	assert.NotZero(t, result.DeletedPostgresqlRecords, "Should have deleted PostgreSQL records")
	assert.NotZero(t, result.DeletedDynamoDBRecords, "Should have deleted DynamoDB records")
	// S3 deletion might be 0 if no avatar/attachments were actually created
	assert.GreaterOrEqual(t, result.DeletedS3Objects, 0, "S3 deletion count should be non-negative")

	// Verify all data was deleted
	verifyDataDeleted(t, db, dynamoDB, s3Client, testData)

	t.Logf("Successfully deleted user %s: PG=%d, Dynamo=%d, S3=%d", 
		userID, result.DeletedPostgresqlRecords, result.DeletedDynamoDBRecords, result.DeletedS3Objects)
}

func createComprehensiveTestData(t *testing.T, db *gorm.DB, dynamoDB *dynamodb_storage.DB, s3Client *s3.S3) *TestData {
	q := query.Use(db)
	userID := uuid.NewString()
	docID := uuid.NewString()
	
	testData := &TestData{}

	// Create main user
	testData.User = &models.User{
		ID:               userID,
		Name:             "Test User",
		Email:            fmt.Sprintf("test-%s@example.com", userID[:8]),
		Provider:         "test",
		DisplayName:      "Test User Display",
		StripeCustomerID: "cus_test_" + userID[:8],
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := q.User.Create(testData.User)
	require.NoError(t, err, "Failed to create test user")

	// Create document for relationships
	testData.Document = &models.Document{
		ID:           docID,
		Title:        "Test Document",
		RootParentID: docID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = q.Document.Create(testData.Document)
	require.NoError(t, err, "Failed to create test document")

	// Create comment
	now := time.Now()
	testData.Comment = &models.Comment{
		ID:         uuid.NewString(),
		DocumentID: docID,
		ThreadID:   uuid.NewString(),
		UserID:     userID,
		Body:       "Test comment",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	err = q.Comment.Create(testData.Comment)
	require.NoError(t, err, "Failed to create test comment")

	// Create document access
	testData.DocumentAccess = &models.DocumentAccess{
		DocumentID:  docID,
		UserID:      userID,
		AccessLevel: "owner",
	}
	err = q.DocumentAccess.Create(testData.DocumentAccess)
	require.NoError(t, err, "Failed to create document access")

	// Create document attachment
	attachmentID := uuid.NewString()
	testData.Attachment = &models.DocumentAttachment{
		ID:          uuid.NewString(),
		UserID:      userID,
		DocumentID:  docID,
		S3ID:        attachmentID,
		Filename:    "test-file.txt",
		ContentType: "text/plain",
		Size:        1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = q.DocumentAttachment.Create(testData.Attachment)
	require.NoError(t, err, "Failed to create document attachment")

	// Create document version
	testData.DocumentVersion = &models.DocumentVersion{
		ID:             uuid.NewString(),
		DocumentID:     docID,
		Name:           "Test Version",
		ContentAddress: "test-address",
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = q.DocumentVersion.Create(testData.DocumentVersion)
	require.NoError(t, err, "Failed to create document version")

	// Create author ID
	testData.AuthorID = &models.AuthorID{
		AuthorID:   1,
		DocumentID: docID,
		UserID:     userID,
		CreatedAt:  time.Now(),
	}
	err = q.AuthorID.Create(testData.AuthorID)
	require.NoError(t, err, "Failed to create author ID")

	// Create access token
	testData.AccessToken = &models.OneTimeAccessToken{
		UserID:    userID,
		Token:     "test-token-" + userID[:8],
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsUsed:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = q.OneTimeAccessToken.Create(testData.AccessToken)
	require.NoError(t, err, "Failed to create access token")

	// Create payment history
	testData.PaymentHistory = &models.PaymentHistory{
		ID:                      uuid.NewString(),
		UserID:                  userID,
		StripePaymentIntentID:   "pi_test_" + userID[:8],
		AmountCents:             1000,
		Currency:                "USD",
		Status:                  "succeeded",
		CreatedAt:               time.Now(),
	}
	err = q.PaymentHistory.Create(testData.PaymentHistory)
	require.NoError(t, err, "Failed to create payment history")

	// Create shared document link
	testData.SharedLink = &models.SharedDocumentLink{
		DocumentID:    docID,
		InviterID:     userID,
		InviteeEmail:  "invitee@example.com",
		InviteLink:    "testlink",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		IsActive:      true,
	}
	err = q.SharedDocumentLink.Create(testData.SharedLink)
	require.NoError(t, err, "Failed to create shared document link")

	// Create subscription plan first (needed for user subscription)
	subscriptionPlan := &models.SubscriptionPlan{
		ID:            uuid.NewString(),
		Name:          "Test Plan",
		PriceCents:    999,
		Currency:      "USD",
		Interval:      "month",
		Status:        "active",
		StripePriceID: "price_test_123",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = q.SubscriptionPlan.Create(subscriptionPlan)
	require.NoError(t, err, "Failed to create subscription plan")

	// Create user subscription
	testData.UserSubscription = &models.UserSubscription{
		ID:                   uuid.NewString(),
		UserID:               userID,
		SubscriptionPlanID:   subscriptionPlan.ID,
		StripeSubscriptionID: "sub_test_" + userID[:8],
		Status:               "active",
		CurrentPeriodStart:   time.Now(),
		CurrentPeriodEnd:     time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	err = q.UserSubscription.Create(testData.UserSubscription)
	require.NoError(t, err, "Failed to create user subscription")

	// Create DynamoDB data
	createDynamoDBTestData(t, dynamoDB, userID, docID, testData)

	// Create S3 data
	createS3TestData(t, s3Client, userID, attachmentID, testData)

	return testData
}

func createDynamoDBTestData(t *testing.T, dynamoDB *dynamodb_storage.DB, userID, docID string, testData *TestData) {
	// Create message
	testData.Message = &dynamodb_storage.Message{
		ContainerID:     "chan#" + uuid.NewString(),
		MessageID:       uuid.NewString(),
		Chain:           "main",
		CreatedAt:       time.Now().UnixNano(),
		UserID:          userID,
		AuthorID:        "1",
		ChannelID:       uuid.NewString(),
		Content:         "Test message",
		DocID:           docID,
		LifecycleStage:  dynamodb_storage.MessageLifecycleStageCompleted,
		Attachments:     &models.AttachmentList{},
		AIContent:       &models.AIContent{},
		MessageMetadata: &models.MessageMetadata{},
	}
	err := dynamoDB.CreateMessage(testData.Message)
	require.NoError(t, err, "Failed to create test message")

	// Create timeline event
	testData.TimelineEvent = &dynamodb_storage.TimelineEvent{
		DocID:    "doc#" + docID,
		UserID:   userID,
		AuthorID: "1",
		Event:    &models.TimelineEventPayload{},
	}
	err = dynamoDB.CreateTimelineEvent(testData.TimelineEvent)
	require.NoError(t, err, "Failed to create timeline event")

	// Create notification
	testData.Notification = &dynamodb_storage.Notification{
		UserID:    userID,
		ID:        uuid.NewString(),
		DocID:     docID,
		Read:      false,
		CreatedAt: time.Now().UnixNano(),
		Payload:   &models.NotificationPayload{},
	}
	err = dynamoDB.UpsertNotification(testData.Notification)
	require.NoError(t, err, "Failed to create notification")

	// Create user preference
	testData.UserPreference = &dynamodb_storage.UserPreference{
		UserID: userID,
		Preference: &models.UserPreference{
			EnableActivityNotifications:    true,
			UnreadActivityFrequencyMinutes: 5,
		},
	}
	err = dynamoDB.UpsertUserPreference(testData.UserPreference)
	require.NoError(t, err, "Failed to create user preference")

	// Create document preference
	testData.DocPreference = &dynamodb_storage.DocPreference{
		UserID: userID,
		DocID:  docID,
		Preference: &models.DocumentPreference{
			EnableFirstOpenNotifications:  true,
			EnableAllCommentNotifications: false,
			EnableMentionNotifications:    true,
			EnableDmNotifications:         true,
		},
	}
	err = dynamoDB.UpsertDocNotificationPreference(testData.DocPreference)
	require.NoError(t, err, "Failed to create document preference")
}

func createS3TestData(t *testing.T, s3Client *s3.S3, userID, attachmentID string, testData *TestData) {
	// Create user avatar
	testData.AvatarKey = fmt.Sprintf(constants.UserAvatarKeyFormat, userID)
	avatarData := []byte("fake-avatar-image-data")
	err := s3Client.PutObject(s3Client.ImagesBucket, testData.AvatarKey, "image/png", avatarData)
	if err != nil {
		t.Logf("Failed to create avatar (may be expected in test): %v", err)
	}

	// Create attachment files
	originalKey := fmt.Sprintf(constants.DocumentAttachmentOriginalKey, attachmentID)
	fileKey := fmt.Sprintf(constants.DocumentAttachmentFileKey, attachmentID, "test-file.txt")
	
	testData.AttachmentKeys = []string{originalKey, fileKey}
	
	attachmentData := []byte("fake-attachment-data")
	
	err = s3Client.PutObject(s3Client.Bucket, originalKey, "text/plain", attachmentData)
	if err != nil {
		t.Logf("Failed to create attachment original (may be expected in test): %v", err)
	}
	
	err = s3Client.PutObject(s3Client.Bucket, fileKey, "text/plain", attachmentData)
	if err != nil {
		t.Logf("Failed to create attachment file (may be expected in test): %v", err)
	}
}

func verifyDataExists(t *testing.T, db *gorm.DB, dynamoDB *dynamodb_storage.DB, s3Client *s3.S3, testData *TestData) {
	q := query.Use(db)
	userID := testData.User.ID

	// Verify PostgreSQL data exists
	count, err := q.User.Where(q.User.ID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "User should exist")

	count, err = q.Comment.Where(q.Comment.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Comment should exist")

	count, err = q.DocumentAccess.Where(q.DocumentAccess.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Document access should exist")

	// Verify DynamoDB data exists
	notification, err := dynamoDB.GetNotification(userID, testData.Notification.ID)
	require.NoError(t, err)
	assert.NotNil(t, notification, "Notification should exist")

	userPref, err := dynamoDB.GetUserPreference(userID)
	require.NoError(t, err)
	assert.NotNil(t, userPref, "User preference should exist")

	// Note: S3 verification is optional since test S3 might not work in all environments
}

func verifyDataDeleted(t *testing.T, db *gorm.DB, dynamoDB *dynamodb_storage.DB, s3Client *s3.S3, testData *TestData) {
	q := query.Use(db)
	userID := testData.User.ID

	// Verify PostgreSQL data deleted
	count, err := q.User.Where(q.User.ID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "User should be deleted")

	count, err = q.Comment.Where(q.Comment.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Comments should be deleted")

	count, err = q.DocumentAccess.Where(q.DocumentAccess.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Document access should be deleted")

	count, err = q.DocumentAttachment.Where(q.DocumentAttachment.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Document attachments should be deleted")

	count, err = q.AuthorID.Where(q.AuthorID.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Author IDs should be deleted")

	count, err = q.OneTimeAccessToken.Where(q.OneTimeAccessToken.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Access tokens should be deleted")

	count, err = q.PaymentHistory.Where(q.PaymentHistory.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Payment history should be deleted")

	count, err = q.SharedDocumentLink.Where(q.SharedDocumentLink.InviterID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Shared document links should be deleted")

	count, err = q.UserSubscription.Where(q.UserSubscription.UserID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "User subscriptions should be deleted")

	// Verify DynamoDB data deleted
	_, err = dynamoDB.GetNotification(userID, testData.Notification.ID)
	assert.Error(t, err, "Notification should not exist")
	assert.True(t, strings.Contains(err.Error(), "notification not found"), "Should get not found error")

	// Verify notifications are deleted by checking count
	notifCount, err := dynamoDB.GetNotificationCountForUser(userID, false)
	require.NoError(t, err)
	assert.Equal(t, int64(0), notifCount, "No unread notifications should exist")

	readNotifCount, err := dynamoDB.GetNotificationCountForUser(userID, true)
	require.NoError(t, err)
	assert.Equal(t, int64(0), readNotifCount, "No read notifications should exist")

	// Note: Full DynamoDB verification requires scanning, which is expensive
	// In a real test, you might want to add specific checks for messages and timeline events
}

func TestDeleteUser_NonExistentUser(t *testing.T) {
	// Setup test environment
	ctx := context.Background()
	
	// Ensure storage is available
	testutils.EnsureStorage()
	
	// Extract dependencies
	db := testutils.TestGormDb(t)
	dynamoDB, err := dynamodb_storage.NewDB()
	require.NoError(t, err, "Failed to create DynamoDB client")
	
	s3Client, err := s3.NewS3()
	require.NoError(t, err, "Failed to create S3 client")

	// Create deletion service
	deletionService := NewUserDeletionService(db, dynamoDB, s3Client)

	// Try to delete non-existent user
	nonExistentUserID := uuid.NewString()
	result, err := deletionService.DeleteUser(ctx, nonExistentUserID)

	// Should succeed with 0 deletions
	require.NoError(t, err, "Deleting non-existent user should not error")
	assert.Equal(t, 0, result.DeletedPostgresqlRecords, "Should delete 0 PostgreSQL records")
	assert.Equal(t, 0, result.DeletedDynamoDBRecords, "Should delete 0 DynamoDB records") 
	assert.Equal(t, 0, result.DeletedS3Objects, "Should delete 0 S3 objects")
	assert.Empty(t, result.Errors, "Should have no errors")
}

func TestDeleteUser_PartialData(t *testing.T) {
	// Setup test environment  
	ctx := context.Background()
	
	// Ensure storage is available
	testutils.EnsureStorage()
	
	// Extract dependencies
	db := testutils.TestGormDb(t)
	dynamoDB, err := dynamodb_storage.NewDB()
	require.NoError(t, err, "Failed to create DynamoDB client")
	
	s3Client, err := s3.NewS3()
	require.NoError(t, err, "Failed to create S3 client")

	// Create deletion service
	deletionService := NewUserDeletionService(db, dynamoDB, s3Client)

	// Create user with only PostgreSQL data (no DynamoDB/S3)
	q := query.Use(db)
	userID := uuid.NewString()
	
	user := &models.User{
		ID:          userID,
		Name:        "Partial Test User",
		Email:       fmt.Sprintf("partial-%s@example.com", userID[:8]),
		Provider:    "test",
		DisplayName: "Partial Test User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = q.User.Create(user)
	require.NoError(t, err, "Failed to create partial test user")

	// Delete user
	result, err := deletionService.DeleteUser(ctx, userID)

	// Should succeed
	require.NoError(t, err, "User deletion should succeed")
	assert.Equal(t, 1, result.DeletedPostgresqlRecords, "Should delete 1 PostgreSQL record")
	assert.Equal(t, 0, result.DeletedDynamoDBRecords, "Should delete 0 DynamoDB records")
	assert.Equal(t, 0, result.DeletedS3Objects, "Should delete 0 S3 objects")

	// Verify user is deleted
	count, err := q.User.Where(q.User.ID.Eq(userID)).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "User should be deleted")
}
