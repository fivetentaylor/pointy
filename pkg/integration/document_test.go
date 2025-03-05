// package integration_test
//
// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"testing"
// 	"time"
//
// 	"github.com/google/uuid"
// 	"github.com/playwright-community/playwright-go"
// 	"github.com/stretchr/testify/suite"
//
// 	"github.com/teamreviso/code/pkg/config"
// 	"github.com/teamreviso/code/pkg/server"
// 	"github.com/teamreviso/code/pkg/server/auth"
// 	"github.com/teamreviso/code/pkg/testutils"
// )
//
// type DocumentSuite struct {
// 	suite.Suite
//
// 	server       *server.Server
// 	jwt          *auth.JWTEngine
// 	addr         string
// 	cancelWorker context.CancelFunc
// }
//
// func (suite *DocumentSuite) SetupTest() {
// 	err := playwright.Install()
// 	if err != nil {
// 		suite.T().Fatal(err)
// 	}
//
// 	testutils.EnsureStorage()
// 	suite.cancelWorker = testutils.RunWorker(suite.T())
// 	gormdb := testutils.TestGormDb(suite.T())
//
// 	jwtsecret := "testtesttesttest"
//
// 	cfg := config.Server{
// 		Addr:            ":0",
// 		JWTSecret:       jwtsecret,
// 		OpenAIKey:       "sk-idontwork",
// 		AllowedOrigins:  []string{"https://*", "http://*"},
// 		WorkerRedisAddr: os.Getenv("REDIS_URL"),
// 		GoogleOauth: &config.GoogleOauth{
// 			ClientID:     "fake-client-id",
// 			ClientSecret: "fake-client-secret",
// 			RedirectURL:  "fake-redirect-url",
// 		},
// 	}
//
// 	s, err := server.NewServer(cfg, gormdb)
// 	if err != nil {
// 		suite.T().Fatal(fmt.Errorf("error creating server: %w", err))
// 	}
//
// 	// Start the server in a goroutine
// 	go func() {
// 		if err := s.ListenAndServe(); err != nil {
// 			suite.T().Fatalf("Failed to start server: %v", err)
// 		}
// 	}()
//
// 	suite.jwt = auth.NewJWT(jwtsecret)
//
// 	time.Sleep(500 * time.Millisecond)
//
// 	suite.server = s
// 	suite.addr = s.Listener.Addr().String() // get the addr of the server
// 	suite.addr = strings.Replace(suite.addr, "[::]", "localhost", 1)
// }
//
// func (suite *DocumentSuite) TearDownTest() {
// 	suite.cancelWorker()
// }
//
// func (suite *DocumentSuite) TestSimpleDocumentInteraction() {
// 	t := suite.T()
// 	ctx := testutils.TestContext()
// 	user := testutils.CreateAdmin(t, ctx)
// 	docID := uuid.NewString()
// 	testutils.CreateTestDocument(t, ctx, docID, "")
// 	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)
//
// 	pw, err := playwright.Run()
// 	suite.Require().NoError(err)
//
// 	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
// 		// Headless: playwright.Bool(false),
// 	})
// 	suite.Require().NoError(err)
//
// 	page, err := browser.NewPage()
// 	suite.Require().NoError(err)
//
// 	page.On("console", func(msg playwright.ConsoleMessage) {
// 		t.Logf("CONSOLE: %s: %s", msg.Type(), msg.Text())
// 	})
//
// 	page.On("pageerror", func(err error) {
// 		t.Logf("PAGE ERROR: %s", err.Error())
// 	})
//
// 	page.On("requestfailed", func(request playwright.Request) {
// 		t.Logf("REQUEST FAILED: %s", request.URL())
// 	})
//
// 	resp, err := page.Goto(fmt.Sprintf("http://%s/test/documents/%s", suite.addr, docID), playwright.PageGotoOptions{
// 		WaitUntil: playwright.WaitUntilStateLoad,
// 	})
// 	suite.Require().NoError(err)
//
// 	fmt.Printf("resp: %+v\n", resp)
//
// 	var locator playwright.Locator
// 	for i := 0; i < 30; i++ {
// 		locator = page.Locator("css=[contenteditable]")
// 		count, _ := locator.Count()
// 		t.Logf("count: %d, locator: %#v", count, locator)
// 		if count > 0 {
// 			break
// 		}
// 		time.Sleep(100 * time.Millisecond)
//
// 	}
//
// 	time.Sleep(500 * time.Millisecond)
// 	_, err = page.WaitForFunction("() => { return window.re.loaded; }", nil)
// 	suite.Require().NoError(err, "expected re.loaded to be true")
//
// 	err = locator.Focus()
// 	suite.Require().NoError(err)
//
// 	t.Logf("Starting typing...")
// 	value, err := locator.InnerHTML()
// 	suite.Require().NoError(err)
//
// 	for _, c := range "hello, world!" {
// 		t.Logf("key: %s", string(c))
// 		page.Keyboard().Down(string(c))
// 		time.Sleep(200 * time.Millisecond)
// 		value, err := locator.InnerHTML()
// 		if err != nil {
// 			suite.T().Fatal(err)
// 		}
// 		t.Logf("value: %s", value)
// 		page.Keyboard().Up(string(c))
// 	}
//
// 	expected := `<p data-rid="q_1"><span data-rid="1_3">hello, world!</span></p>`
//
// 	value, err = locator.InnerHTML()
// 	suite.Require().NoError(err)
// 	suite.Equal(expected, value, "expected innerHTML %q but got %q", expected, value)
//
// 	for i := 0; i < 30; i++ {
// 		html := testutils.GetDocumentHtml(t, ctx, docID)
// 		if html == expected {
// 			break
// 		}
// 		time.Sleep(100 * time.Millisecond)
// 	}
//
// 	html := testutils.GetDocumentHtml(t, ctx, docID)
// 	suite.Equal(expected, html, "expected document html %q but got %q", expected, html)
//
// 	suite.Require().NoError(browser.Close())
// 	suite.Require().NoError(pw.Stop())
// }
//
// func TestDocumentSuite(t *testing.T) {
// 	t.Parallel()
// 	suite.Run(t, new(DocumentSuite))
// }
