package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/fivetentaylor/pointy/pkg/graph"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

type MessagingTestSuite struct {
	suite.Suite
	handler http.Handler
}

func (suite *MessagingTestSuite) SetupTest() {
	testutils.EnsureStorage()

	suite.handler = graph.NewHandler()
}

func TestMessagingTestSuite(t *testing.T) {
	suite.Run(t, new(MessagingTestSuite))
}

func (suite *MessagingTestSuite) RunQuery(ctx context.Context, query string) (map[string]interface{}, error) {
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

	log.Println(rr.Body.String())

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
