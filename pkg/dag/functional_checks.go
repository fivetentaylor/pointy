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
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/rogue"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type FunctionalCheckFile struct {
	ID        string `json:"id"`
	DagName   string `json:"dagName"`
	CheckName string `json:"checkName"`
	CreatedAt int64  `json:"createdAt"`

	DocumentId string `json:"documentId"`

	BeforeAddress    string `json:"beforeAddress"`
	SerializedRogue  string `json:"serializedRogue"`
	SerializedThread string `json:"serializedThread"`
}

func ListFuncationalChecks(ctx context.Context, dagName string) ([]FunctionalCheckFile, error) {
	log := env.Log(ctx)
	log.Info("listing checks", "dagName", dagName)

	var checks []FunctionalCheckFile
	err := filepath.Walk(filepath.Join("checks", "dags", dagName), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "check.json" {
			bts, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var check FunctionalCheckFile
			err = json.Unmarshal(bts, &check)
			if err != nil {
				return err
			}

			checks = append(checks, check)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(checks, func(i, j int) bool {
		return checks[i].CreatedAt > checks[j].CreatedAt
	})

	return checks, nil
}

func GetFuncationalCheck(ctx context.Context, dagName, id string) (*FunctionalCheckFile, error) {
	log := env.Log(ctx)
	log.Info("getting check", "dagName", dagName, "id", id)

	path := filepath.Join("checks", "dags", dagName, id, "check.json")
	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var check FunctionalCheckFile
	err = json.Unmarshal(bts, &check)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling check: %s", err)
	}

	return &check, nil
}

func (f FunctionalCheckFile) Prepare(ctx context.Context, userId string) (*dynamo.Message, *dynamo.Message, *dynamo.Thread, error) {
	log := env.Log(ctx)

	err := f.resetDocument(ctx, userId)
	if err != nil {
		log.Error("[funcational check] error ensuring document exists and resetting", "error", err)
		return nil, nil, nil, fmt.Errorf("error ensuring document exists and resetting it: %s", err)
	}

	input, output, thread, err := f.createNewMessageThread(ctx, userId)
	if err != nil {
		log.Error("[funcational check] error creating new message thread", "error", err)
		return nil, nil, nil, fmt.Errorf("error creating new message thread: %s", err)
	}

	return input, output, thread, nil
}

func (f FunctionalCheckFile) Run(ctx context.Context, userId string, d *Dag) (*FunctionalCheckFileResult, error) {
	log := env.Log(ctx)

	input, output, thread, err := f.Prepare(ctx, userId)
	if err != nil {
		log.Error("[funcational check] error preparing", "error", err)
		return nil, err
	}

	log.Info("[funcational check] created new message thread", "inputMessageId", input.MessageID, "outputMessageId", output.MessageID, "threadId", thread.ThreadID)

	result := &FunctionalCheckFileResult{
		DagName:    f.DagName,
		CheckID:    f.ID,
		ThreadId:   thread.ThreadID,
		DocumentId: f.DocumentId,
	}

	err = result.CaptureBeforeState(ctx)
	if err != nil {
		log.Error("[funcational check] error capturing before state", "error", err)
		return nil, err
	}

	d.ParentId = f.DocumentId
	err = d.Run(ctx, map[string]any{
		"docId":    f.DocumentId,
		"threadId": thread.ThreadID,
		"authorId": "check",
		"userId":   userId,

		"inputMessageId":  input.MessageID,
		"outputMessageId": output.MessageID,
	})
	if err != nil {
		log.Error("[funcational check] failed to respond to thread", "error", err)
		return nil, stackerr.Wrap(err)
	}

	result.DagId = d.Uuid

	err = result.CaptureResult(ctx)
	if err != nil {
		log.Error("[funcational check] error capturing result", "error", err)
		return nil, err
	}

	err = result.Save()
	if err != nil {
		log.Error("[funcational check] error saving result", "error", err)
		return result, err
	}

	err = result.Evaluate(ctx)
	if err != nil {
		log.Error("[funcational check] error evaluating state", "error", err)
	}

	err = result.Save()
	if err != nil {
		log.Error("[funcational check] error saving result after evaluation", "error", err)
		return result, err
	}

	return result, nil
}

func CreateFunctionalCheckFile(
	ctx context.Context, docId, threadID, dagName, checkName string) (*FunctionalCheckFile, error) {
	doc, err := rogue.CurrentDocument(ctx, docId)
	if err != nil {
		return nil, err
	}

	rb, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	dydb := env.Dynamo(ctx)

	messages, err := dydb.GetMessagesForThread(threadID)
	if err != nil {
		return nil, err
	}

	mb, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	fc := &FunctionalCheckFile{
		DocumentId: docId,
		CheckName:  checkName,
		DagName:    dagName,

		SerializedRogue:  string(rb),
		SerializedThread: string(mb),
	}

	// extract before address
	lastMessage := messages[len(messages)-1]
	if lastMessage.MessageMetadata != nil {
		fc.BeforeAddress = lastMessage.MessageMetadata.ContentAddressBefore
	}

	return fc, nil
}

func (f FunctionalCheckFile) createNewMessageThread(ctx context.Context, userId string) (*dynamo.Message, *dynamo.Message, *dynamo.Thread, error) {
	dydb := env.Dynamo(ctx)

	thread := &dynamo.Thread{
		DocID:  f.DocumentId,
		UserID: userId,
		Title:  fmt.Sprintf("Funcational check: %s", f.DagName),
	}
	err := dydb.CreateThread(thread)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating thread: %s", err)
	}

	var messages []*dynamo.Message
	err = json.Unmarshal([]byte(f.SerializedThread), &messages)
	if err != nil {
		return nil, nil, nil, err
	}

	outputIndex := len(messages) - 1
	inputIndex := outputIndex - 1

	var input, output *dynamo.Message
	for i, m := range messages {
		m.MessageID = uuid.NewString()
		m.ContainerID = fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, thread.ThreadID)

		if m.UserID != constants.RevisoUserID {
			m.UserID = userId
		}

		if i == inputIndex {
			input = m
		}

		if i == outputIndex {
			output = m

			// Clear out content and attachments
			m.LifecycleStage = dynamo.MessageLifecycleStagePending
			m.Attachments.Attachments = []*models.Attachment{}
			m.Content = ""
			m.MessageMetadata = &models.MessageMetadata{}
		}

		err := dydb.CreateMessage(m)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return input, output, thread, nil
}

func (f FunctionalCheckFile) resetDocument(ctx context.Context, userId string) error {
	log := env.Log(ctx)
	q := env.Query(ctx)
	docTbl := q.Document

	// Find doc
	var doc *models.Document
	docs, err := docTbl.Where(docTbl.ID.Eq(f.DocumentId)).Find()
	if err != nil {
		return stackerr.Errorf("[check %s] error finding document for check: %s", f.ID, err)
	}
	if len(docs) == 0 {
		doc = &models.Document{
			ID:           f.DocumentId,
			Title:        f.CheckName,
			IsPublic:     true,
			RootParentID: f.DocumentId,
			// TODO figure out rogue version
		}

		err := docTbl.Create(doc)
		if err != nil {
			return fmt.Errorf("[check %s] error creating document for check (user %s, doc %s) inserting %#v: %s", f.ID, userId, f.DocumentId, doc, err)
		}
	}

	docAccessTbl := q.DocumentAccess
	accesses, err := docAccessTbl.Where(docAccessTbl.UserID.Eq(userId), docAccessTbl.DocumentID.Eq(f.DocumentId)).Find()
	if err != nil {
		return stackerr.Errorf("[check %s] error finding access for check: %s", f.ID, err)
	}
	if len(accesses) == 0 {
		access := &models.DocumentAccess{
			UserID:      userId,
			DocumentID:  f.DocumentId,
			AccessLevel: "owner",
		}
		err := docAccessTbl.Create(access)
		if err != nil {
			return fmt.Errorf("[check %s] error creating access for check (user %s, doc %s) inserting %#v: %s", f.ID, userId, f.DocumentId, access, err)
		}
	}

	err = rogue.SnapshotDoc(ctx, f.DocumentId)

	seq, err := rogue.GetLastS3Seq(ctx, f.DocumentId)
	log.Info("last s3 seq", "seq", seq)
	if err != nil {
		return stackerr.Errorf("[check %s] error getting last s3 seq: %s", f.ID, err)
	}
	newSeq := seq + 1

	var rdoc *v3.Rogue
	log.Info("json", "json", f.SerializedRogue)
	err = json.Unmarshal([]byte(f.SerializedRogue), &rdoc)
	if err != nil {
		return stackerr.Errorf("[check %s] error unmarshaling doc: %s", f.ID, err)
	}

	log.Info("[check %s] resetting document", "id", f.DocumentId, "seq", newSeq, "address", f.BeforeAddress)

	var addr v3.ContentAddress
	if f.BeforeAddress != "" {
		err = json.Unmarshal([]byte(f.BeforeAddress), &addr)
		if err != nil {
			return stackerr.Errorf("[check %s] error unmarshaling address: %s", f.ID, err)
		}
	}

	rdoc, err = rdoc.GetOldRogue(&addr)
	if err != nil {
		return stackerr.Errorf("[check %s] error getting address rogue: %s", f.ID, err)
	}

	// Save rogue doc
	err = rogue.SaveDocToS3(ctx, f.DocumentId, newSeq, rdoc)
	if err != nil {
		return stackerr.Errorf("[check %s] error saving doc: %s", f.ID, err)
	}

	return nil
}

func (f FunctionalCheckFile) Save() error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	if f.CreatedAt == 0 {
		f.CreatedAt = time.Now().Unix()
	}

	bts, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	file := fmt.Sprintf("./checks/dags/%s/%s/check.json", f.DagName, f.ID)
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
