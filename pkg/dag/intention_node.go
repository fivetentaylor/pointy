package dag

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

const intentionPrompt = `You are an AI assistant designed to analyze user queries in a writing chat bot context. Your task is to determine the user's intent based on their input. Carefully examine the user's message and categorize it according to the following criteria:

Query Type:
a) General question (not requiring full document context)
b) Document analysis request (requiring full document context, but no updates)
c) Document update request
If it's a document update request, determine the scope:
a) Whole document update
b) Subsection update

Provide your analysis in the following format:
Explanation: [Brief explanation of your reasoning]
Query Type: [General Question / Document Analysis Request / Document Update Request]
Document Scope: [Full Document / Subsection / Not Applicable]
Update Required: [Yes / No]
Confidence: [High / Medium / Low]
Example user inputs and how to analyze them:

User: "What's the best way to start a persuasive essay?"
Analysis:
Explanation: The user is asking for general advice about essay writing, not referring to a specific document or requesting any analysis or updates.
Query Type: General Question
Document Scope: Not Applicable
Update Required: No
Confidence: High

User: "What do you think of my essay?"
Analysis:
Explanation: The user is asking for an opinion on their essay, which requires reading and analyzing the full document, but not making any updates to it.
Query Type: Document Analysis Request
Document Scope: Full Document
Update Required: No
Confidence: High

User: "Can you write me an example todo list?"
Analysis:
Explanation: The user explicitly asks to write, indicating a whole document update.
Query Type: Document Update Request
Document Scope: Full Document
Update Required: Yes
Confidence: High

User: "Can you revise the entire document to improve the flow?"
Analysis:
Explanation: The user explicitly asks for a revision of the entire document, indicating a whole document update.
Query Type: Document Update Request
Document Scope: Full Document
Update Required: Yes
Confidence: High

User: "Please update the conclusion paragraph to better summarize the main points."
Analysis:
Explanation: The user requests an update to a specific part of the document (the conclusion paragraph), indicating a subsection update.
Query Type: Document Update Request
Document Scope: Subsection
Update Required: Yes
Confidence: High

User: "Is my argument in the third paragraph consistent with my thesis?"
Analysis:
Explanation: The user is asking about the consistency of a specific part of their essay, which requires reading the full document for context, but doesn't require any updates.
Query Type: Document Analysis Request
Document Scope: Full Document
Update Required: No
Confidence: High

User: "What should I title this essay?"
Analysis:
Explanation: The user is asking for advice on titling their essay, which require reading the full document.
Query Type: Document Analysis Request
Document Scope: Full Document
Update Required: No
Confidence: High

Always analyze the user's input carefully, considering both explicit requests and implied intentions. If you're unsure about the intent, use a lower confidence rating and explain your reasoning. Pay special attention to queries that require reading the full document but don't necessarily require updates.

User: `

type IntentionNode struct {
	BaseLLMNode

	DocumentAnalysisNode Node
	DocumentUpdateNode   Node
	DefaultNode          Node
}

type IntentionNodeInputs struct {
	PromptName string `key:"promptName"` // selectedPromptName
	DocId      string `key:"docId"`
	ThreadId   string `key:"threadId"`
	AuthorId   string `key:"authorId"`

	InputMessageId  string `key:"inputMessageId"`
	OutputMessageId string `key:"outputMessageId"`
}

func (n *IntentionNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)
	log.Info("👾 running intention node")

	input := &IntentionNodeInputs{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	inMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.InputMessageId)
	if err != nil {
		return nil, fmt.Errorf("error getting input message: %s", err)
	}

	outMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.OutputMessageId)
	if err != nil {
		return nil, fmt.Errorf("error getting output message: %s", err)
	}

	outMsg.LifecycleStage = dynamo.MessageLifecycleStagePending
	outMsg.LifecycleReason = "Understanding User Intent"
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	resp, err := n.GenerateFromSinglePrompt(ctx, intentionPrompt+inMsg.Content)
	if err != nil {
		return nil, fmt.Errorf("error generating response: %s", err)
	}

	log.Info("[intention] response", "response", resp)

	data, err := parseIntentionText(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %s", err)
	}

	outMsg.LifecycleReason = "Analysis Complete"
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	SetStateKey(ctx, "intention", data)

	if data.UpdateRequired == "Yes" {
		return n.DocumentUpdateNode, nil
	}

	if data.QueryType == "Document Analysis Request" {
		return n.DocumentAnalysisNode, nil
	}

	return n.DefaultNode, nil
}

func lastCompletedMessage(messages []*dynamo.Message) *dynamo.Message {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.LifecycleStage == dynamo.MessageLifecycleStageCompleted {
			return msg
		}
	}
	return nil
}

type IntentionData struct {
	QueryType      string
	DocumentScope  string
	UpdateRequired string
	Confidence     string
}

func parseIntentionText(text string) (IntentionData, error) {
	var data IntentionData
	scanner := bufio.NewScanner(strings.NewReader(text))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Query Type:") {
			data.QueryType = strings.TrimSpace(strings.TrimPrefix(line, "Query Type:"))
		} else if strings.HasPrefix(line, "Document Scope:") {
			data.DocumentScope = strings.TrimSpace(strings.TrimPrefix(line, "Document Scope:"))
		} else if strings.HasPrefix(line, "Confidence:") {
			data.Confidence = strings.TrimSpace(strings.TrimPrefix(line, "Confidence:"))
		} else if strings.HasPrefix(line, "Update Required:") {
			data.UpdateRequired = strings.TrimSpace(strings.TrimPrefix(line, "Update Required:"))
		}
	}

	if err := scanner.Err(); err != nil {
		return IntentionData{}, err
	}

	return data, nil
}
