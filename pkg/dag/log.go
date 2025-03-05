package dag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
)

type LogNode struct {
	Next Node
}

func (n LogNode) Run(ctx context.Context) (Node, error) {
	state := GetState(ctx)

	log := env.Log(ctx)
	log.Info("state", "state", state)

	return n.Next, nil
}

func saveLogFile(ctx context.Context, filename string, content any) {
	log := env.Log(ctx)
	now := time.Now()
	invert := 999999999999999999 - (now.UnixNano() / 10)

	d := GetDag(ctx)
	runPath := fmt.Sprintf("%s/%s/%d/%s_%s", d.ParentId, d.Name, invert, d.Uuid, now.Format(time.RFC3339))
	s3 := env.S3(ctx)

	node := GetCurrentNode(ctx)
	if node != nil {
		filename = fmt.Sprintf("%s-%s", GetNodeType(node), filename)
	}

	key := fmt.Sprintf("%s/%s/%s", constants.DagsDir, runPath, filename)

	if str, ok := content.(string); ok {
		err := s3.PutObject(s3.Bucket, key, "text/plain", []byte(str))
		if err != nil {
			log.Error("error saving log", "error", err, "key", key)
		}
		return
	}

	// otherwise, save it as json
	bts, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Error("error marshalling content", "error", err, "content", content)
	}

	err = s3.PutObject(s3.Bucket, key, "application/json", bts)
	if err != nil {
		log.Error("error saving log", "error", err, "key", key)
	}

	return
}

type HTTPLogger struct {
	id      string
	counter int
	*http.Client
}

func (l *HTTPLogger) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	defer l.incr()

	start := time.Now()

	// Log the request
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	reqLog := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
		"body":   string(reqBody),
	}
	go saveLogFile(ctx, fmt.Sprintf("request-%d.json", start.Unix()), reqLog)

	// Execute the request
	resp, err := l.Client.Do(req)
	if err != nil {
		go saveLogFile(ctx, fmt.Sprintf("error-%d.json", start.Unix()), map[string]interface{}{"error": err})
		return nil, err
	}

	// Log the response
	var respBody bytes.Buffer
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, &respBody))

	respLog := map[string]interface{}{
		"status":      resp.Status,
		"status_code": resp.StatusCode,
		"header":      resp.Header,
		"body":        respBody.String(),
		"duration_ms": time.Since(start).Milliseconds(),
	}
	go saveLogFile(ctx, fmt.Sprintf("response-%d.json", start.Unix()), respLog)

	return resp, nil
}

func (l *HTTPLogger) incr() {
	l.counter++
}
