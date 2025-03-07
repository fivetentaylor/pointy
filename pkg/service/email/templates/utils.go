package templates

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fivetentaylor/pointy/pkg/constants"
)

func webHostUrl(ctx context.Context, path string) string {
	baseURL := ctx.Value(constants.WebHostContextKey).(string)
	base, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse base url for app_host %s: %v", base, err))
	}

	rel, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	// Resolve the reference
	fullURL := base.ResolveReference(rel)

	return fullURL.String()
}

func appHostUrl(ctx context.Context, path string) string {
	baseURL := ctx.Value(constants.AppHostContextKey).(string)
	base, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse base url for app_host %s: %v", base, err))
	}

	rel, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	fullURL := base.ResolveReference(rel)
	return fullURL.String()
}
