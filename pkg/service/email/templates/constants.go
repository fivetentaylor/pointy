package templates

import (
	"context"
	"fmt"
	"regexp"

	"github.com/fivetentaylor/pointy/pkg/models"
)

const AttachmentMaxContentLength = 200

var atMentionRegex = regexp.MustCompile(`@:user:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}:([^@]+)@`)
var atMentionRegexV2 = regexp.MustCompile(`@@([A-Za-z0-9+/=]+)@@`)

func DocUrl(ctx context.Context, doc *models.Document) string {
	return appHostUrl(ctx, fmt.Sprintf("/drafts/%s", doc.ID))
}
