package auth

import (
	"context"
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/service/images"
	"github.com/teamreviso/code/pkg/utils"
)

var (
	ErrEmailAlreadyInUse = errors.New("email is already in use")
)

type userIdentity struct {
	ID           string
	Email        string
	Name         string
	Provider     string
	Picture      *string
	PasswordHash string
}

func findUserIdentity(ctx context.Context, email, provider string) (*userIdentity, error) {
	userTbl := env.Query(ctx).User

	user, err := userTbl.
		Where(userTbl.Email.Eq(email)).
		Where(userTbl.Provider.Eq(provider)).
		First()

	if err != nil {
		return nil, err
	}

	return &userIdentity{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Provider:     user.Provider,
		PasswordHash: user.PasswordHash,
	}, nil
}

func findUserIdentityByOneTimeAccessToken(ctx context.Context, token string) (*userIdentity, error) {
	userTbl := env.Query(ctx).User

	accessTokenTbl := env.Query(ctx).OneTimeAccessToken

	accessToken, err := accessTokenTbl.
		Where(accessTokenTbl.Token.Eq(token)).
		Where(accessTokenTbl.ExpiresAt.Gt(time.Now())).
		First()
	if err != nil {
		log.Errorf("error finding access token: %s", err.Error())
		return nil, err
	}

	user, err := userTbl.
		Where(userTbl.ID.Eq(accessToken.UserID)).
		First()
	if err != nil {
		log.Errorf("error finding user from access token %s: %s", accessToken.UserID, err.Error())
		return nil, err
	}

	accessToken.IsUsed = true
	err = accessTokenTbl.Save(accessToken)
	if err != nil {
		return nil, err
	}

	return &userIdentity{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Provider:     user.Provider,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (ui userIdentity) Identifier() string {
	normalizedEmail := utils.NormalizeEmail(ui.Email)
	return fmt.Sprintf("%s:%s", ui.Provider, normalizedEmail)
}

func (ui userIdentity) token(ctx context.Context) (string, error) {
	user, err := ui.UpdateOrCreateUser(ctx)
	if err != nil {
		return "", err
	}
	return JWT(ctx).GenerateUserToken(user)
}

func (ui userIdentity) FindOrCreateUser(ctx context.Context) (*models.User, error) {
	user, err := ui.User(ctx)
	if err != nil {
		if err.Error() == "record not found" {
			return ui.CreateUser(ctx)
		}
		return nil, fmt.Errorf("error finding or creating user: %w", err)
	}

	return user, nil
}

func (ui userIdentity) UpdateOrCreateUser(ctx context.Context) (*models.User, error) {
	user, err := ui.User(ctx)
	if err != nil {
		if err.Error() == "record not found" {
			return ui.CreateUser(ctx)
		}
		return nil, fmt.Errorf("error finding or creating user: %w", err)
	}

	// don't replace picture if we already have one
	if user.Picture != nil && ui.Picture == nil {
		ui.Picture = user.Picture
	}

	err = ui.UpdateUser(ctx, user)
	return user, err
}

func (ui userIdentity) CreateUser(ctx context.Context) (*models.User, error) {
	q := env.Query(ctx)
	log := env.Log(ctx)

	log.Infof("creating user %s", ui.Identifier())

	normalizedEmail := utils.NormalizeEmail(ui.Email)

	if ok := utils.IsValidEmail(normalizedEmail); !ok {
		return nil, fmt.Errorf("invalid email")
	}

	var (
		avatar image.Image
		err    error
	)
	if ui.Picture != nil && strings.HasPrefix(*ui.Picture, "http") {
		log.Infof("downloading avatar from %s", *ui.Picture)
		avatar, err = images.DownloadImage(*ui.Picture)
		if err != nil {
			log.Errorf("error downloading avatar: %s", err)
		}
	}

	var user *models.User
	err = q.Transaction(func(tx *query.Query) error {
		user = &models.User{
			Name:         ui.Name,
			DisplayName:  ui.Name,
			Email:        normalizedEmail,
			Provider:     ui.Provider,
			Picture:      ui.Picture,
			PasswordHash: ui.PasswordHash,
		}

		log.Infof("creating user %+v", user)

		err := tx.User.Create(user)
		if err != nil {
			log.Errorf("error creating user: %s", err)
			return fmt.Errorf("error creating user: %w", err)
		}

		return nil
	})
	if err != nil {
		if strings.Index(err.Error(), "duplicate key value violates unique constraint") != -1 {
			return nil, ErrEmailAlreadyInUse
		}
		return nil, err
	}

	if avatar != nil {
		err = images.UpdateUserAvatar(ctx, user, avatar)
		if err != nil {
			log.Errorf("error updating user avatar: %s", err)
			return user, nil
		}
	}

	log.Infof("created user %s", ui.Identifier())

	sendSignupEventToLoops(ctx, normalizedEmail)

	return user, nil
}

func (ui userIdentity) UpdateUser(ctx context.Context, user *models.User) error {
	q := env.Query(ctx)

	normalizedEmail := utils.NormalizeEmail(ui.Email)

	if ok := utils.IsValidEmail(normalizedEmail); !ok {
		return fmt.Errorf("invalid email")
	}

	err := q.Transaction(func(tx *query.Query) error {
		user.Name = ui.Name
		user.Email = normalizedEmail
		user.Provider = ui.Provider
		user.Picture = ui.Picture
		user.PasswordHash = ui.PasswordHash

		err := tx.User.Save(user)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		return nil
	})
	if err != nil {
		if strings.Index(err.Error(), "duplicate key value violates unique constraint") != -1 {
			return errors.New("email is already in use")
		}
		return err
	}

	return nil
}

func (ui userIdentity) User(ctx context.Context) (*models.User, error) {
	userTbl := env.Query(ctx).User

	normalizedEmail := utils.NormalizeEmail(ui.Email)

	user, err := userTbl.
		Where(userTbl.Email.Eq(normalizedEmail)).
		First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func sendSignupEventToLoops(ctx context.Context, email string) error {
	if os.Getenv("ENV") == "development" {
		return nil
	}

	loopsClient, err := client.NewLoopsClientFromEnv()
	if err != nil {
		log.Info("Failed to create loops client", "error", err)
		return err
	}

	err = loopsClient.SendEvent(ctx, "signup", client.LoopsContactProperties{Email: email}, nil)
	if err != nil {
		log.Info("Failed to send event to loops", "error", err)
		return err
	}
	return nil
}
