package dag

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	v3 "github.com/teamreviso/code/rogue/v3"
)

type FunctionalCheckExampleFile struct {
	ID        string `json:"id"`
	CheckID   string `json:"checkId"`
	DagName   string `json:"dagName"`
	CreatedAt int64  `json:"createdAt"`

	DocumentId string `json:"documentId"`

	SerializedRogueBefore  string `json:"serializedRogueBefore"`
	SerializedThreadBefore string `json:"serializedThreadBefore"`
	SerializedRogueResult  string `json:"serializedRogueResult"`
	SerializedThreadResult string `json:"serializedThreadResult"`

	Approved bool `json:"approved"`
}

func ListFuncationalCheckExamples(ctx context.Context, dagName, checkID string) ([]FunctionalCheckExampleFile, error) {
	log := env.Log(ctx)
	log.Info("listing check results", "checkID", checkID)

	var examples []FunctionalCheckExampleFile
	path := filepath.Join("checks", "dags", dagName, checkID, "examples")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return examples, nil
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			log.Info("found example", "path", path)

			bts, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var example FunctionalCheckExampleFile
			err = json.Unmarshal(bts, &example)
			if err != nil {
				return err
			}

			examples = append(examples, example)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(examples, func(i, j int) bool {
		return examples[i].CreatedAt > examples[j].CreatedAt
	})

	return examples, nil
}

func GetFuncationalCheckExample(ctx context.Context, dagName, checkID, id string) (*FunctionalCheckExampleFile, error) {
	log := env.Log(ctx)
	log.Info("getting result", "id", id, "checkID", checkID)

	path := filepath.Join("checks", "dags", dagName, checkID, "examples", fmt.Sprintf("%s.json", id))
	bts, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var result FunctionalCheckExampleFile
	err = json.Unmarshal(bts, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling check: %s", err)
	}

	return &result, nil
}

func CreateFuncationalCheckExampleFile(
	ctx context.Context, result *FunctionalCheckFileResult,
) (*FunctionalCheckExampleFile, error) {
	fce := &FunctionalCheckExampleFile{
		ID:      result.ID,
		CheckID: result.CheckID,
		DagName: result.DagName,

		DocumentId: result.DocumentId,
		CreatedAt:  result.CreatedAt,

		SerializedRogueBefore:  result.SerializedRogueBefore,
		SerializedThreadBefore: result.SerializedThreadBefore,
		SerializedRogueResult:  result.SerializedRogueResult,
		SerializedThreadResult: result.SerializedThreadResult,
	}

	return fce, nil
}

func (r *FunctionalCheckExampleFile) OutputMessage() (*dynamo.Message, error) {
	var messages []dynamo.Message

	err := json.Unmarshal([]byte(r.SerializedThreadResult), &messages)
	if err != nil {
		return nil, err
	}

	return &messages[len(messages)-1], nil
}

func (r *FunctionalCheckExampleFile) BeforeDoc() (*v3.Rogue, error) {
	var result *v3.Rogue
	err := json.Unmarshal([]byte(r.SerializedRogueBefore), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling before doc: %s", err)
	}

	return result, nil
}

func (r *FunctionalCheckExampleFile) ResultDoc() (*v3.Rogue, error) {
	var result *v3.Rogue
	err := json.Unmarshal([]byte(r.SerializedRogueResult), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling result doc: %s", err)
	}

	return result, nil
}

func (r *FunctionalCheckExampleFile) Diff() (string, error) {
	if r.SerializedRogueBefore == "" || r.SerializedRogueResult == "" {
		return "", nil
	}

	before, err := r.BeforeDoc()
	if err != nil {
		return "", err
	}

	result, err := r.ResultDoc()
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

func (f FunctionalCheckExampleFile) Save() error {
	if f.ID == "" {
		return fmt.Errorf("id must be set to the result id")
	}
	if f.CreatedAt == 0 {
		f.CreatedAt = time.Now().Unix()
	}

	bts, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	file := fmt.Sprintf("./checks/dags/%s/%s/examples/%s.json", f.DagName, f.CheckID, f.ID)
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
