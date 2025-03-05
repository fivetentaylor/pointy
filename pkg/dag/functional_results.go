package dag

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/rogue"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	v3 "github.com/teamreviso/code/rogue/v3"
)

type FunctionalCheckFileResult struct {
	ID        string `json:"id"`
	CheckID   string `json:"checkId"`
	DagName   string `json:"dagName"`
	CreatedAt int64  `json:"createdAt"`

	DocumentId string `json:"documentId"`
	ThreadId   string `json:"threadId"`
	DagId      string `json:"dagId"`

	SerializedRogueBefore  string `json:"serializedRogueBefore"`
	SerializedThreadBefore string `json:"serializedThreadBefore"`
	SerializedRogueResult  string `json:"serializedRogueResult"`
	SerializedThreadResult string `json:"serializedThreadResult"`

	Prompt      string `json:"prompt"`
	RawResponse string `json:"rawResponse"`

	Justification string  `json:"justification"`
	Score         float64 `json:"score"`
	Assessment    string  `json:"assessment"`
}

func ListFuncationalCheckResults(ctx context.Context, checkID string) ([]FunctionalCheckFileResult, error) {
	log := env.Log(ctx)
	log.Info("listing check results", "checkID", checkID)

	path := filepath.Join("results", "dags", checkID)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []FunctionalCheckFileResult
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "result.json" {
			log.Info("found result", "path", path)

			bts, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var result FunctionalCheckFileResult
			err = json.Unmarshal(bts, &result)
			if err != nil {
				return err
			}

			results = append(results, result)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("check not found: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt > results[j].CreatedAt
	})

	return results, nil
}

func GetFuncationalCheckResult(ctx context.Context, checkID, id string) (*FunctionalCheckFileResult, error) {
	log := env.Log(ctx)
	log.Info("getting result", "id", id, "checkID", checkID)

	path := filepath.Join("results", "dags", checkID, id, "result.json")
	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result FunctionalCheckFileResult
	err = json.Unmarshal(bts, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling check: %s", err)
	}

	return &result, nil
}

func (r *FunctionalCheckFileResult) CaptureBeforeState(ctx context.Context) error {
	doc, err := rogue.CurrentDocument(ctx, r.DocumentId)
	if err != nil {
		return err
	}

	rb, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	dydb := env.Dynamo(ctx)

	messages, err := dydb.GetMessagesForThread(r.ThreadId)
	if err != nil {
		return err
	}

	mb, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	r.SerializedRogueBefore = string(rb)
	r.SerializedThreadBefore = string(mb)

	return nil
}

func (r *FunctionalCheckFileResult) CaptureResult(ctx context.Context) error {
	doc, err := rogue.CurrentDocument(ctx, r.DocumentId)
	if err != nil {
		return err
	}

	rb, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	dydb := env.Dynamo(ctx)

	messages, err := dydb.GetMessagesForThread(r.ThreadId)
	if err != nil {
		return err
	}

	mb, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	r.SerializedRogueResult = string(rb)
	r.SerializedThreadResult = string(mb)

	return nil
}

func (r *FunctionalCheckFileResult) BeforeDoc() (*v3.Rogue, error) {
	var result *v3.Rogue
	err := json.Unmarshal([]byte(r.SerializedRogueBefore), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling before doc: %s", err)
	}

	return result, nil
}

func (r *FunctionalCheckFileResult) ResultDoc() (*v3.Rogue, error) {
	var result *v3.Rogue
	err := json.Unmarshal([]byte(r.SerializedRogueResult), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling result doc: %s", err)
	}

	return result, nil
}

func (r *FunctionalCheckFileResult) InputMessage() (*dynamo.Message, error) {
	var messages []dynamo.Message

	err := json.Unmarshal([]byte(r.SerializedThreadResult), &messages)
	if err != nil {
		return nil, err
	}

	return &messages[len(messages)-2], nil
}

func (r *FunctionalCheckFileResult) OutputMessage() (*dynamo.Message, error) {
	var messages []dynamo.Message

	err := json.Unmarshal([]byte(r.SerializedThreadResult), &messages)
	if err != nil {
		return nil, err
	}

	return &messages[len(messages)-1], nil
}

func (r *FunctionalCheckFileResult) Diff() (string, error) {
	if r.SerializedRogueBefore == "" || r.SerializedRogueResult == "" {
		return "", nil
	}

	var before *v3.Rogue
	err := json.Unmarshal([]byte(r.SerializedRogueBefore), &before)
	if err != nil {
		return "", err
	}

	var result *v3.Rogue
	err = json.Unmarshal([]byte(r.SerializedRogueResult), &result)
	if err != nil {
		return "", err
	}

	beforeMkdown, err := before.GetFullMarkdown()
	if err != nil {
		return "", err
	}

	resultMkdown, err := result.GetFullMarkdown()
	if err != nil {
		return "", err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(beforeMkdown, resultMkdown, false)

	return dmp.DiffPrettyText(diffs), nil
}

func (r *FunctionalCheckFileResult) Save() error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.CreatedAt == 0 {
		r.CreatedAt = time.Now().Unix()
	}

	bts, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	file := fmt.Sprintf("./results/dags/%s/%s/result.json", r.CheckID, r.ID)
	err = os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bts, 0644)
	if err != nil {
		return err
	}

	return nil
}
