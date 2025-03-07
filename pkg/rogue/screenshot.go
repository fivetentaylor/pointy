package rogue

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/server/auth"
)

const screenshotWidth = int64(640)
const screenshotHeight = int64(600)
const screenshotUrlExpireDuration = 30 * time.Minute

func SeqToScreenshotS3Path(docID string, seq int64, theme string) string {
	return path.Join(
		constants.S3Prefix,
		docID,
		"screenshots",
		fmt.Sprintf("%s-%s.png", invertSeq(seq), theme),
	)
}

func GetLastScreenshotsURL(ctx context.Context, docID string) ([]string, error) {
	s3 := env.S3(ctx)
	screenPath := path.Join(constants.S3Prefix, docID, "screenshots")
	screenKeys, err := s3.ListFirstTwo(s3.Bucket, screenPath)
	if err != nil {
		return nil, fmt.Errorf("error getting last screenshot: %s", err)
	}

	darkUrl, err := s3.GetPresignedUrl(
		s3.Bucket,
		screenKeys[0],
		screenshotUrlExpireDuration,
	)
	if err != nil {
		return nil, err
	}

	lightUrl, err := s3.GetPresignedUrl(
		s3.Bucket,
		screenKeys[1],
		screenshotUrlExpireDuration,
	)
	if err != nil {
		return nil, err
	}

	return []string{darkUrl, lightUrl}, nil
}

func (ds *DocStore) ScreenshotDoc(ctx context.Context, docID string) error {
	appHost := os.Getenv("WEB_HOST")
	seq, _, err := ds.GetCurrentDoc(ctx, docID)
	if err != nil {
		return fmt.Errorf("ScreenshotDoc(%q): error getting current doc: %s", docID, err)
	}

	cctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	cctx, cancel = context.WithTimeout(cctx, 25*time.Second)
	defer cancel()

	screenshotURL := path.Join(appHost, "screenshot", docID)
	darkScreenshotURL := path.Join(appHost, "screenshot", docID, "?theme=dark")

	var screenshot []byte
	var darkScreenshot []byte

	log.Info("screenshot url", "url", screenshotURL)

	jwt := auth.NewJWT(os.Getenv("JWT_SECRET"))

	userTbl := env.Query(ctx).User
	user, err := userTbl.Where(userTbl.ID.Eq(constants.RevisoUserID)).First()
	if err != nil {
		return fmt.Errorf("ScreenshotDoc(%q): error getting reviso user: %s", docID, err)
	}

	token, err := jwt.GenerateUserToken(user)
	if err != nil {
		return fmt.Errorf("error generating token: %s", err)
	}

	headers := map[string]interface{}{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	err = chromedp.Run(cctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.EmulateViewport(screenshotWidth, screenshotHeight),
		// attach headers

		chromedp.Navigate(screenshotURL),
		chromedp.WaitVisible("#screenshot-page", chromedp.ByQuery),
		chromedp.FullScreenshot(&screenshot, 100),

		chromedp.Navigate(darkScreenshotURL),
		chromedp.WaitVisible("#screenshot-page", chromedp.ByQuery),
		chromedp.FullScreenshot(&darkScreenshot, 100),
	)
	if err != nil {
		return fmt.Errorf("SnapshotDoc(%q): error taking screenshot: %s", docID, err)
	}

	key := SeqToScreenshotS3Path(docID, seq, "light")
	err = ds.S3.PutCachedObject(ds.S3.Bucket, key, "image/x-png", screenshot)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	key = SeqToScreenshotS3Path(docID, seq, "dark")
	err = ds.S3.PutCachedObject(ds.S3.Bucket, key, "image/x-png", darkScreenshot)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	log.Info("screenshots saved", "url", screenshotURL)

	return nil
}
