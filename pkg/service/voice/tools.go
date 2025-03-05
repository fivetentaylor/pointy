package voice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

type Tool struct {
	Definition ToolDefinition
	Function   func(ctx context.Context, event ResponseOutputItemDone, session *Session) error
}

var UpdateDocumentTool = Tool{
	Definition: ToolDefinition{
		Type:        "function",
		Name:        "update_document",
		Description: "Update the users document",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolProperty{
				"description_of_changes": {
					Type:        "string",
					Description: "A request to another AI agent to update the document. Be very descriptive",
				},
			},
			Required: []string{"description_of_changes"},
		},
	},
	Function: UpdateDocument,
}

func UpdateDocument(
	ctx context.Context,
	event ResponseOutputItemDone,
	session *Session,
) error {
	log := env.SLog(ctx)
	inputKey := fmt.Sprintf("input-%s", event.ResponseID)
	outputKey := fmt.Sprintf("output-%s", event.ResponseID)

	log.Info("🔈 [voice.UpdateDocument] updating document", "event", event, "arguments", event.Item.Arguments)

	inMsgContent, err := GetArgument[string](event, "description_of_changes")
	if err != nil {
		log.Error("[voice.UpdateDocument] error getting input message", "error", err)
		return err
	}
	if inMsgContent == "" {
		log.Error("[voice.UpdateDocument] error input message is empty")
		return fmt.Errorf("input message is empty")
	}

	addrStr, err := session.CurrentMarshalledContentAddress(ctx)
	if err != nil {
		log.Error("[voice.UpdateDocument] error getting content address", "error", err)
		return err
	}

	inMsg, err := session.GetRequestMessage(ctx, inputKey, &dynamo.Message{
		Hidden:  true,
		Content: inMsgContent,
	})
	if err != nil {
		log.Error("[voice.UpdateDocument] error getting output message", "error", err)
		return err
	}
	defer session.RemoveMessage(ctx, inputKey)

	outMsg, err := session.GetResponseMessage(ctx, outputKey, &dynamo.Message{
		LifecycleReason: "Thinking",
		LifecycleStage:  dynamo.MessageLifecycleStagePending,
		MessageMetadata: &models.MessageMetadata{
			ContentAddressBefore: addrStr,
		},
	})
	if err != nil {
		log.Error("[voice.UpdateDocument] error getting output message", "error", err)
		return err
	}

	d := UpdateDocumentDag()
	d.ParentId = session.documentID
	d.OnError = func(ctx context.Context, node dag.Node, err error) {
		log.Error("[voice.UpdateDocument] error running dag", "error", err)
		err = session.SendToRealtime(ctx, ConversationItemCreate{
			Base: Base{
				Type:    "conversation.item.create",
				EventID: GenerateID("evt_", 21),
			},
			Item: ItemDetails{
				ID:     GenerateID("msg", 21),
				Type:   "function_call_output",
				CallID: event.Item.CallID,
				Output: "{\"error\": \"" + err.Error() + "\"}",
			},
		})
		if err != nil {
			log.Error("[voice.UpdateDocument] error sending error event", "error", err)
		}

		err = session.SendToRealtime(ctx, ResponseCreate{
			Response: ResponseDetails{
				Instructions: onRevisionErrorInstructions,
			},
			Base: Base{
				EventID: GenerateID("evt_", 21),
				Type:    "response.create",
			},
		})
		if err != nil {
			log.Error("[voice.UpdateDocument] error sending second response event", "error", err)
		}
	}
	d.OnComplete = func(ctx context.Context, d *dag.Dag) {
		explanation := d.State().Get("revision_explanation")
		conclusion := d.State().Get("revision_conclusion")

		output := map[string]any{
			"success":     "true",
			"explanation": explanation,
			"conclusion":  conclusion,
		}

		bts, err := json.Marshal(output)
		if err != nil {
			log.Error("[voice.UpdateDocument] error encoding output", "error", err)
			return
		}

		err = session.SendToRealtime(ctx, ConversationItemCreate{
			Base: Base{
				Type:    "conversation.item.create",
				EventID: GenerateID("evt_", 21),
			},
			Item: ItemDetails{
				ID:     GenerateID("msg", 21),
				Type:   "function_call_output",
				CallID: event.Item.CallID,
				Output: string(bts),
			},
		})
		if err != nil {
			log.Error("[voice.UpdateDocument] error sending error event", "error", err)
		}

		instructions := fmt.Sprintf(
			onRevisionSuccessInstructions,
			explanation,
			conclusion,
		)

		log.Info("🔈 [voice.UpdateDocument] creating response", "instructions", instructions)

		err = session.SendToRealtime(ctx, ResponseCreate{
			Response: ResponseDetails{
				Instructions: instructions,
				Metadata: map[string]any{
					"outputKey": outputKey,
				},
			},
			Base: Base{
				EventID: GenerateID("evt_", 21),
				Type:    "response.create",
			},
		})
		if err != nil {
			log.Error("[voice.UpdateDocument] error sending response event", "error", err)
		}

		outMsg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
		err = messaging.UpdateMessage(ctx, outMsg)
		if err != nil {
			log.Error("[voice.UpdateDocument] error updating message", "error", err)
		}
	}

	err = d.Run(ctx, map[string]any{
		"docId":    session.documentID,
		"threadId": session.threadID,
		"authorId": session.authorID,
		"userId":   session.userID,

		"inputMessage":  inMsg,
		"outputMessage": outMsg,
	})
	if err != nil {
		log.Error("[voice.UpdateDocument] error running voice update dag", "error", err)
		return err
	}

	return nil
}
