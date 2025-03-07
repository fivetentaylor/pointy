package graph_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

func (suite *GraphTestSuite) TestUploadImage() {
	ctx := testutils.TestContext()
	user := testutils.CreateAdmin(suite.T(), ctx)
	ctx = testutils.WithUserClaimForUser(user)(ctx)

	// Create a test document
	docID := uuid.NewString()
	testutils.CreateTestDocument(suite.T(), ctx, docID, "")
	testutils.AddOwnerToDocument(suite.T(), ctx, docID, user.ID)

	// Open an actual image file
	file, err := os.Open("taylor.png") // Make sure this file exists in your project
	require.NoError(suite.T(), err)
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	require.NoError(suite.T(), err)

	// Create Upload struct with actual file data
	actualFile := graphql.Upload{
		File:        file,
		Filename:    fileInfo.Name(),
		Size:        fileInfo.Size(),
		ContentType: "image/png",
	}

	// Prepare uploadImage mutation
	variables := map[string]interface{}{
		"file":  actualFile,
		"docId": docID,
	}

	// Run uploadImage mutation
	result, err := suite.RunUploadQuery(ctx, uploadImageMutation, variables)
	require.NoError(suite.T(), err)

	// Assert the result
	uploadedImage, ok := result["uploadImage"].(map[string]interface{})
	require.True(suite.T(), ok)
	require.NotEmpty(suite.T(), uploadedImage["id"])
	require.Equal(suite.T(), docID, uploadedImage["docId"])

	status := "LOADING"
	for status == "LOADING" {
		result, err := suite.RunQuery(ctx, makeQuery(getImageQuery, map[string]interface{}{"docId": docID, "imageId": uploadedImage["id"]}))
		fmt.Printf("result: %v\n", result)
		require.NoError(suite.T(), err)
		getImage, ok := result["getImage"].(map[string]interface{})
		require.True(suite.T(), ok)
		status = getImage["status"].(string)
	}

	// Verify the image was uploaded by listing document images
	listImagesResult, err := suite.RunQuery(ctx, makeQuery(listDocumentImagesQuery, map[string]interface{}{"docId": docID}))
	require.NoError(suite.T(), err)

	images, ok := listImagesResult["listDocumentImages"].([]interface{})
	require.True(suite.T(), ok)
	require.Equal(suite.T(), 1, len(images))

	uploadedImageID := uploadedImage["id"].(string)
	firstImage := images[0].(map[string]interface{})
	require.Equal(suite.T(), uploadedImageID, firstImage["id"])
	require.Equal(suite.T(), docID, firstImage["docId"])

	// Optionally, verify the content of the uploaded image
	getImageResult, err := suite.RunQuery(ctx, makeQuery(getImageSignedUrlQuery, map[string]interface{}{"docId": docID, "imageId": uploadedImageID}))
	require.NoError(suite.T(), err)

	signedURL, ok := getImageResult["getImageSignedUrl"].(map[string]interface{})["url"].(string)
	require.True(suite.T(), ok)
	require.NotEmpty(suite.T(), signedURL)

	// Download the image using the signed URL and compare with original
	resp, err := http.Get(signedURL)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// If you want to see the link open up
	// openURLInBrowser(signedURL)
}

func (suite *GraphTestSuite) TestUploadGif() {
	ctx := testutils.TestContext()
	user := testutils.CreateAdmin(suite.T(), ctx)
	ctx = testutils.WithUserClaimForUser(user)(ctx)

	// Create a test document
	docID := uuid.NewString()
	testutils.CreateTestDocument(suite.T(), ctx, docID, "")
	testutils.AddOwnerToDocument(suite.T(), ctx, docID, user.ID)

	// Open an actual image file
	file, err := os.Open("piper_cat.gif") // Make sure this file exists in your project
	require.NoError(suite.T(), err)
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	require.NoError(suite.T(), err)

	// Create Upload struct with actual file data
	actualFile := graphql.Upload{
		File:        file,
		Filename:    fileInfo.Name(),
		Size:        fileInfo.Size(),
		ContentType: "image/gif",
	}

	// Prepare uploadImage mutation
	variables := map[string]interface{}{
		"file":  actualFile,
		"docId": docID,
	}

	// Run uploadImage mutation
	result, err := suite.RunUploadQuery(ctx, uploadImageMutation, variables)
	require.NoError(suite.T(), err)

	// Assert the result
	uploadedImage, ok := result["uploadImage"].(map[string]interface{})
	require.True(suite.T(), ok)
	require.NotEmpty(suite.T(), uploadedImage["id"])
	require.Equal(suite.T(), docID, uploadedImage["docId"])

	status := "LOADING"
	for status == "LOADING" {
		result, err := suite.RunQuery(ctx, makeQuery(getImageQuery, map[string]interface{}{"docId": docID, "imageId": uploadedImage["id"]}))
		fmt.Printf("result: %v\n", result)
		require.NoError(suite.T(), err)
		getImage, ok := result["getImage"].(map[string]interface{})
		require.True(suite.T(), ok)
		status = getImage["status"].(string)
	}

	// Verify the image was uploaded by listing document images
	listImagesResult, err := suite.RunQuery(ctx, makeQuery(listDocumentImagesQuery, map[string]interface{}{"docId": docID}))
	require.NoError(suite.T(), err)

	images, ok := listImagesResult["listDocumentImages"].([]interface{})
	require.True(suite.T(), ok)
	require.Equal(suite.T(), 1, len(images))

	uploadedImageID := uploadedImage["id"].(string)
	firstImage := images[0].(map[string]interface{})
	require.Equal(suite.T(), uploadedImageID, firstImage["id"])
	require.Equal(suite.T(), docID, firstImage["docId"])

	// Optionally, verify the content of the uploaded image
	getImageResult, err := suite.RunQuery(ctx, makeQuery(getImageSignedUrlQuery, map[string]interface{}{"docId": docID, "imageId": uploadedImageID}))
	require.NoError(suite.T(), err)

	signedURL, ok := getImageResult["getImageSignedUrl"].(map[string]interface{})["url"].(string)
	require.True(suite.T(), ok)
	require.NotEmpty(suite.T(), signedURL)

	// Download the image using the signed URL and compare with original
	resp, err := http.Get(signedURL)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// If you want to see the link open up
	// openURLInBrowser(signedURL)
}

// GraphQL query/mutation strings
const (
	uploadImageMutation = `
		mutation($file: Upload!, $docId: ID!) {
			uploadImage(file: $file, docId: $docId) {
				id
				docId
			}
		}
	`

	getImageQuery = `
    query($docId: ID!, $imageId: ID!) {
      getImage(docId: $docId, imageId: $imageId) {
        id
        docId
        status
      }
    }
  `

	listDocumentImagesQuery = `
		query($docId: ID!) {
			listDocumentImages(docId: $docId) {
				id
				docId
			}
		}
	`

	getImageSignedUrlQuery = `
		query($docId: ID!, $imageId: ID!) {
			getImageSignedUrl(docId: $docId, imageId: $imageId) {
				url
				expiresAt
			}
		}
	`
)

func openURLInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
