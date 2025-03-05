package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/suite"
	"github.com/teamreviso/code/pkg/config"
	"github.com/teamreviso/code/pkg/graph"
	"github.com/teamreviso/code/pkg/server"
	"github.com/teamreviso/code/pkg/server/auth"
	"github.com/teamreviso/code/pkg/testutils"
)

type GraphTestSuite struct {
	suite.Suite

	handler      http.Handler
	server       *server.Server
	jwt          *auth.JWTEngine
	addr         string
	cancelWorker context.CancelFunc
}

func (suite *GraphTestSuite) SetupTest() {
	testutils.EnsureStorage()

	suite.cancelWorker = testutils.RunWorker(suite.T())
	gormdb := testutils.TestGormDb(suite.T())

	jwtsecret := "testtesttesttest"

	cfg := config.Server{
		Addr:            ":0",
		JWTSecret:       jwtsecret,
		OpenAIKey:       "sk-idontwork",
		AllowedOrigins:  []string{"https://*", "http://*"},
		WorkerRedisAddr: os.Getenv("REDIS_URL"),
		GoogleOauth: &config.GoogleOauth{
			ClientID:     "fake-client-id",
			ClientSecret: "fake-client-secret",
			RedirectURL:  "fake-redirect-url",
		},
	}

	s, err := server.NewServer(cfg, gormdb)
	if err != nil {
		suite.T().Fatal(fmt.Errorf("error creating server: %w", err))
	}

	// Start the server in a goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			suite.T().Fatalf("Failed to start server: %v", err)
		}
	}()

	suite.jwt = auth.NewJWT(jwtsecret)

	time.Sleep(500 * time.Millisecond)

	suite.server = s
	suite.addr = s.Listener.Addr().String() // get the addr of the server
	suite.addr = strings.Replace(suite.addr, "[::]", "localhost", 1)
	suite.handler = graph.NewHandler()
}

func TestGraphTestSuite(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}

func (suite *GraphTestSuite) RunQuery(ctx context.Context, query string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v\n\n%v", status, rr.Body.String())
	}

	response := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	var ok bool

	// no error in response
	if err, ok := response["errors"]; ok {
		return nil, fmt.Errorf("unexpected error in response: %v", err)
	}

	// no data in response
	var dataRaw interface{}
	if dataRaw, ok = response["data"]; !ok {
		return nil, fmt.Errorf("unexpected no data in response: %v", dataRaw)
	}

	var data map[string]interface{}
	if data, ok = dataRaw.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("unexpected type of data in response: %v", data)
	}

	return data, nil
}

func RunQuery[T any](ctx context.Context, handler http.Handler, query string) (*T, error) {
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v\n\n%v", status, rr.Body.String())
	}

	var response struct {
		Data   T              `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return &response.Data, nil
}

type GraphQLError struct {
	Message string `json:"message"`
	// Add other fields as needed
}

func makeQuery(operationString string, variables map[string]interface{}) string {
	request := map[string]interface{}{
		"query": operationString, // This key is used for both queries and mutations
	}

	if variables != nil && len(variables) > 0 {
		request["variables"] = variables
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func mapGet(m map[string]interface{}, keys ...string) (interface{}, bool) {
	var current interface{} = m

	for _, key := range keys {
		v := reflect.ValueOf(current)

		if v.Kind() != reflect.Map {
			return nil, false
		}

		mapValue := v.MapIndex(reflect.ValueOf(key))
		if !mapValue.IsValid() {
			return nil, false
		}

		current = mapValue.Interface()
	}

	return current, true
}

func (suite *GraphTestSuite) RunUploadQuery(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	if isFileUpload(variables) {
		req, err = createMultipartRequest(query, variables)
	} else {
		hydratedQuery := makeQuery(query, variables)
		req, err = http.NewRequest("POST", "/", bytes.NewBufferString(hydratedQuery))
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	suite.handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v\n\n%v", status, rr.Body.String())
	}

	response := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	if err, ok := response["errors"]; ok {
		return nil, fmt.Errorf("unexpected error in response: %v", err)
	}

	dataRaw, ok := response["data"]
	if !ok {
		return nil, fmt.Errorf("unexpected no data in response: %v", dataRaw)
	}

	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type of data in response: %v", data)
	}

	return data, nil
}

func isFileUpload(variables map[string]interface{}) bool {
	for _, v := range variables {
		if _, ok := v.(graphql.Upload); ok {
			return true
		}
	}

	return false
}

func createMultipartRequest(query string, variables map[string]interface{}) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Prepare the operations
	fileMap := make(map[string][]string)
	var fileUpload graphql.Upload
	for k, v := range variables {
		if upload, ok := v.(graphql.Upload); ok {
			fileUpload = upload
			fileMap["0"] = []string{fmt.Sprintf("variables.%s", k)}
			variables[k] = nil
			break // We're only handling one file for now
		}
	}

	operations := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	operationsJSON, err := json.Marshal(operations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operations: %v", err)
	}
	err = writer.WriteField("operations", string(operationsJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to write operations field: %v", err)
	}

	// Write the map
	mapJSON, err := json.Marshal(fileMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map: %v", err)
	}
	err = writer.WriteField("map", string(mapJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to write map field: %v", err)
	}

	// Add file part
	if fileUpload.File != nil {
		part, err := writer.CreateFormFile("0", fileUpload.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %v", err)
		}
		_, err = io.Copy(part, fileUpload.File)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file content: %v", err)
		}
		// Reset the file reader for potential future use
		if seeker, ok := fileUpload.File.(io.Seeker); ok {
			_, err = seeker.Seek(0, io.SeekStart)
			if err != nil {
				return nil, fmt.Errorf("failed to reset file reader: %v", err)
			}
		}
	} else {
		return nil, fmt.Errorf("no file found in variables")
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
