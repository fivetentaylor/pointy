package timeline

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/graph/loaders"
	"github.com/teamreviso/code/pkg/graph/model"
	"github.com/teamreviso/code/pkg/models"
)

func HydrateMentions(ctx context.Context, msg *model.TLMessageV1, mentionedUserIds []string) (*model.TLMessageV1, error) {
	if mentionedUserIds == nil {
		return msg, nil
	}

	_, err := env.UserClaim(ctx)
	if err != nil {
		log.Errorf("error getting current user: %s", err)
		return nil, fmt.Errorf("please login")
	}

	mentionedUsers, err := loaders.GetUsers(ctx, mentionedUserIds)
	if err != nil {
		log.Errorf("error loading mentioned users: %s", err)
		return nil, fmt.Errorf("error loading mentioned users: %s", err)
	}

	if len(mentionedUsers) > 0 {
		// Extract all mentions from the content
		mentions := ExtractMentions(ctx, msg.Content, nil)

		// Create a map of user IDs to users for quick lookup
		userMap := make(map[string]*models.User)
		for _, user := range mentionedUsers {
			userMap[user.ID] = user
		}

		// Update each mention with the correct user information
		for _, mention := range mentions {
			if user, ok := userMap[mention.ID]; ok {
				msg.Content = updateSingleBase64Mention(msg.Content, user, mention.Base64Match)
			}
		}

		// Handle old format mentions if necessary
		if !Base64Regex.MatchString(msg.Content) {
			for _, user := range mentionedUsers {
				// OLD deprecated mention format
				// Replace mentions without dynamic name
				oldMention := fmt.Sprintf("@:user:%s@", user.ID)
				newMention := fmt.Sprintf("@:user:%s:%s@", user.ID, user.DisplayName)
				msg.Content = strings.ReplaceAll(msg.Content, oldMention, newMention)

				// Replace mentions with dynamic name
				oldMentionPattern := fmt.Sprintf("@:user:%s:[^@]*@", user.ID)
				msg.Content = regexp.MustCompile(oldMentionPattern).ReplaceAllString(msg.Content, newMention)
			}
		}
	}

	return msg, nil
}

type UpdateBase64MentionsOptions struct {
	SpecificMatch string
}

func UnfurlBase64Mentions(content string, users []*models.User) string {
	newContent := content
	for _, user := range users {
		matches := Base64Regex.FindAllStringSubmatch(newContent, -1)
		for _, match := range matches {
			userName := user.Name
			newContent = strings.ReplaceAll(newContent, match[0], userName)
		}
	}
	return newContent
}

func UpdateBase64Mentions(content string, user *models.User, options UpdateBase64MentionsOptions) string {
	if options.SpecificMatch == "" {
		// If no specific match is provided, update all matches
		matches := Base64Regex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			content = updateSingleBase64Mention(content, user, match[0])
		}
	} else {
		// Update only the specific match
		content = updateSingleBase64Mention(content, user, options.SpecificMatch)
	}
	return content
}

func updateSingleBase64Mention(content string, user *models.User, base64Match string) string {
	newMention := fmt.Sprintf(":user:%s:%s", user.ID, user.DisplayName)
	encodedStr := base64.StdEncoding.EncodeToString([]byte(newMention))

	replacement := strings.ReplaceAll(content, base64Match, fmt.Sprintf("@@%s@@", encodedStr))
	return replacement
}

var Base64Regex = regexp.MustCompile(`@@([A-Za-z0-9+/=]+)@@`)
var mentionRegex = regexp.MustCompile(`^:user:([a-zA-Z0-9-]+?):(.+)$`)

// return a slice of structs with ID and Name
type Mention struct {
	ID          string
	Name        string
	Base64Match string
}

/*
	 pull any at-mentions from the msg.Content. at-mentions are base64 encoded in the form of @@:user:ID:name@@
		example: @@:user:77d01aca-e6a6-47c8-9df6-61cfe342b739:justin@@ how are you @@:user:87d01aca-e6a6-47c8-9df6-61cfe342b739:joe@@, you are @@:user:77d01aca-e6a6-47c8-9df6-61cfe342b739:justin@@
	 we only want the ID part, and we don't want duplicates
*/
func ExtractMentions(ctx context.Context, content string, existingIDs []string) []Mention {
	mentions := []Mention{}
	mentionMap := make(map[string]bool)
	newUserMap := make(map[string]bool)
	for _, id := range existingIDs {
		mentionMap[id] = true
	}

	matches := Base64Regex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			continue
		}

		decodedStr := string(decodedBytes)
		mentionMatch := mentionRegex.FindStringSubmatch(decodedStr)
		base64Match, mentionID, mentionEmail := match[0], mentionMatch[1], mentionMatch[2]

		if len(mentionMatch) >= 2 && !mentionMap[mentionID] && !newUserMap[mentionEmail] {
			if mentionID == "new-user-id" {
				newUserMap[mentionEmail] = true
			} else {
				mentionMap[mentionID] = true
			}
			mentions = append(mentions, Mention{ID: mentionID, Name: mentionEmail, Base64Match: base64Match})
		}
	}

	return mentions
}

func HydrateFlaggedVersion(ctx context.Context, update *model.TLUpdateV1, timelineUpdate *models.TimelineDocumentUpdateV1) (*model.TLUpdateV1, error) {
	if timelineUpdate.FlaggedVersionId == "" {
		return update, nil
	}

	docVersionTbl := env.Query(ctx).DocumentVersion
	userTbl := env.Query(ctx).User

	docVersion, err := docVersionTbl.Where(docVersionTbl.ID.Eq(timelineUpdate.FlaggedVersionId)).First()
	if err != nil {
		return nil, fmt.Errorf("error getting flagged version: %w", err)
	}

	user, err := userTbl.Where(userTbl.ID.Eq(docVersion.CreatedBy)).First()
	if err != nil {
		return nil, fmt.Errorf("error getting flagged by user: %w", err)
	}

	update.FlaggedVersionID = &docVersion.ID
	update.FlaggedVersionName = &docVersion.Name
	update.FlaggedVersionCreatedAt = &docVersion.CreatedAt
	update.FlaggedByUser = user

	return update, nil
}
