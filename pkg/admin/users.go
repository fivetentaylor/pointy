package admin

import (
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"

	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/users"
	dynamodb_storage "github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
)

const usersPerPage = 20

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	searchEmail := r.URL.Query().Get("search")
	
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	offset := (page - 1) * usersPerPage
	
	q := query.Use(env.RawDB(ctx))
	
	// Build query for count
	var totalUsers int64
	var err error
	if searchEmail != "" {
		totalUsers, err = q.User.Where(q.User.Email.Like("%" + searchEmail + "%")).Count()
	} else {
		totalUsers, err = q.User.Count()
	}
	if err != nil {
		log.Errorf("Failed to count users: %v", err)
		http.Error(w, "Failed to count users", http.StatusInternalServerError)
		return
	}
	
	// Get users for current page
	var usersList []*models.User
	if searchEmail != "" {
		usersList, err = q.User.Where(q.User.Email.Like("%" + searchEmail + "%")).
			Order(q.User.CreatedAt.Desc()).
			Offset(offset).
			Limit(usersPerPage).
			Find()
	} else {
		usersList, err = q.User.Order(q.User.CreatedAt.Desc()).
			Offset(offset).
			Limit(usersPerPage).
			Find()
	}
	if err != nil {
		log.Errorf("Failed to get users: %v", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	
	// Convert to template format
	adminUsers := make([]templates.AdminUser, len(usersList))
	for i, user := range usersList {
		adminUsers[i] = templates.AdminUser{
			ID:               user.ID,
			Name:             user.Name,
			Email:            user.Email,
			DisplayName:      user.DisplayName,
			Provider:         user.Provider,
			Admin:            user.Admin,
			Educator:         user.Educator,
			StripeCustomerID: user.StripeCustomerID,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
		}
	}
	
	totalPages := int((totalUsers + int64(usersPerPage) - 1) / int64(usersPerPage))
	
	paginationInfo := templates.PaginationInfo{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  int(totalUsers),
		ItemsPerPage: usersPerPage,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}
	
	templates.Users(adminUsers, paginationInfo, searchEmail).Render(ctx, w)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userID")
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	// Create deletion service
	dynamoDB, err := dynamodb_storage.NewDB()
	if err != nil {
		log.Errorf("Failed to create DynamoDB client: %v", err)
		http.Error(w, "Failed to initialize deletion service", http.StatusInternalServerError)
		return
	}
	
	s3Client, err := s3.NewS3()
	if err != nil {
		log.Errorf("Failed to create S3 client: %v", err)
		http.Error(w, "Failed to initialize deletion service", http.StatusInternalServerError)
		return
	}
	
	deletionService := users.NewUserDeletionService(env.RawDB(ctx), dynamoDB, s3Client)
	
	// Perform user deletion
	result, err := deletionService.DeleteUser(ctx, userID)
	if err != nil {
		log.Errorf("User deletion failed for %s: %v", userID, err)
		w.Header().Set("HX-Trigger", "user-delete-error")
		templates.DeleteUserError(userID, err.Error()).Render(ctx, w)
		return
	}
	
	log.Infof("Successfully deleted user %s: PG=%d, Dynamo=%d, S3=%d", 
		userID, result.DeletedPostgresqlRecords, result.DeletedDynamoDBRecords, result.DeletedS3Objects)
	
	// Return success response for HTMX
	w.Header().Set("HX-Trigger", "user-deleted")
	templates.DeleteUserSuccess(userID, result).Render(ctx, w)
}