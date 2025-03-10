package attachments_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"

	"github.com/fivetentaylor/pointy/pkg/service/attachments"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

func Test_ExtractedText(t *testing.T) {
	testutils.EnsureStorage()

	ctx := testutils.TestContext()
	user := testutils.CreateUser(t, ctx)
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "")

	testCases := []struct {
		desc                string
		filepath            string
		contentType         string
		expectTextToInclude string
		expectError         string
	}{
		{
			desc:                "extracted text from DOCX",
			filepath:            "./testdata/example.docx",
			contentType:         "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			expectTextToInclude: `Login with Google 🖤🤪 👍✅`,
		},
		{
			desc:        "extracted text from PDF",
			filepath:    "./testdata/attention.pdf",
			contentType: "application/pdf",
			expectTextToInclude: `The dominant sequence transduction models are based on complex recurrent or
convolutional neural networks that include an encoder and a decoder.`,
		},
		// 		{
		// 			desc:        "extracted text from pages doc",
		// 			filepath:    "./testdata/pages.pages",
		// 			contentType: "application/x-iwork-pages-sffpages",
		// 			expectTextToInclude: `The dominant sequence transduction models are based on complex recurrent or
		// convolutional neural networks that include an encoder and a decoder.`,
		// 		},
		{
			desc:                "extracted text from RTF",
			filepath:            "./testdata/reinventing.rtf",
			contentType:         "text/rtf",
			expectTextToInclude: `This is the story of why we reinvented the text editor to integrate into LLMs, and why we spent a year doing it.`,
		},
		{
			desc:        "extracted text from empty markdown",
			filepath:    "./testdata/empty.md",
			contentType: "text/markdown",
			expectError: "no text extracted from attachment",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			file, size := testutils.File(t, tC.filepath)

			attachment, err := attachments.Upload(ctx, graphql.Upload{
				File:        file,
				ContentType: tC.contentType,
				Filename:    "example.docx",
				Size:        size,
			}, docID, user.ID)

			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(attachment)
			text, err := attachments.ExtractedText(ctx, attachment)
			if tC.expectError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tC.expectError) {
					t.Errorf("expected error to contain %s, got %s", tC.expectError, err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			// text should contain the expected text
			if !strings.Contains(text, tC.expectTextToInclude) {
				outname := fmt.Sprintf("%s-extracted.txt", strings.ReplaceAll(tC.desc, " ", "-"))
				t.Errorf("expected text to contain %s, see %s", tC.expectTextToInclude, outname)

				f, err := os.OpenFile(
					outname,
					os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				if _, err := f.WriteString(text); err != nil {
					t.Fatal(err)
				}
			}

		})
	}
}
