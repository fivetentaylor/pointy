package dag

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type EvaluationPromptData struct {
	USER_MESSAGE         string
	INITIAL_DOCUMENT     string
	KNOWN_GOOD_RESPONSES []struct {
		MESSAGE          string
		UPDATED_DOCUMENT string
	}
	KNOWN_BAD_RESPONSES []struct {
		MESSAGE          string
		UPDATED_DOCUMENT string
	}
	MESSAGE          string
	UPDATED_DOCUMENT string
}

type Score float64

// UnmarshalXML custom unmarshaler for Score to handle conversion from string to int
func (s *Score) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	err := d.DecodeElement(&content, &start)
	if err != nil {
		return err
	}
	content = strings.TrimSpace(content)
	val, err := strconv.Atoi(content)
	if err != nil {
		return err
	}
	*s = Score(val)
	return nil
}

type Evaluation struct {
	XMLName       xml.Name `xml:"evaluation"`
	Justification string   `xml:"justification"`
	Score         Score    `xml:"score"`
	Assessment    string   `xml:"assessment"`
}

const evaluationPrompt = `You are tasked with evaluating the quality of a new AI-generated response by comparing it to known good responses. Your goal is to determine if the new response is also a good response based on its similarity and appropriateness relative to the known good responses and it's dissimilarity to the known bad responses.

First, review the user message and initial document that were used as inputs for all responses:

<user_message>
{{.USER_MESSAGE}}
</user_message>

<initial_document>
{{.INITIAL_DOCUMENT}}
</initial_document>

Now, examine the known good responses. Each response consists of a message back to the user and an updated document:

<known_good_responses>
{{- range $val := .KNOWN_GOOD_RESPONSES}}
<known_good_response>
<message>
{{$val.MESSAGE}}
</message>
<updated_document>
{{$val.UPDATED_DOCUMENT}}
</updated_document>
</known_good_response>

{{- end}}	
</known_good_responses>

Now, examine the known bad responses. Each response consists of a message back to the user and an updated document:

<known_bad_responses>
{{- range $val := .KNOWN_BAD_RESPONSES}}
<known_bad_response>
<message>
{{$val.MESSAGE}}
</message>
<updated_document>
{{$val.UPDATED_DOCUMENT}}
</updated_document>
</known_bad_response>

{{- end}}	
</known_bad_responses>

Next, review the new response that needs to be evaluated:

<new_response>
<message>
{{.MESSAGE}}
</message>
<updated_document>
{{.UPDATED_DOCUMENT}}
</updated_document>
</new_response>

To analyze and compare the responses, follow these steps:

1. Identify the key elements and main points addressed in the known good responses.
2. Compare the new response to the known good responses in terms of:
   a. Content coverage: Does it address the same key points?
   b. Tone and style: Is it consistent with the tone and style of the known good responses?
   c. Document updates: Are the changes made to the document similar in nature and extent?
   d. Accuracy: Does it provide accurate information based on the initial document and user message?
   e. Completeness: Does it fully address the user's query or request?

3. Note any significant differences or similarities between the new response and the known good responses.

4. Consider whether any differences are acceptable variations or potential issues.

Based on your analysis, provide a final assessment of whether the new response is a good response. Include a detailed justification for your assessment, referencing specific aspects of the responses in your explanation.

Present your evaluation in the following format:

<evaluation>
<justification>
[Provide your detailed justification here, explaining why you believe the new response is or is not a good response. Reference specific aspects of the responses to support your reasoning.]
</justification>

<score>
[Provide a score from 1 to 10, where 1 indicates the new change is completely inconsistent with the known good changes, and 10 indicates it is highly consistent and likely to be a good change.]
</score>

<assessment>
[State your final assessment as either "GOOD" if you believe the new response is a good response, or "NOT GOOD" if you believe it is not a good response.]
</assessment>
</evaluation>
`

var evaluationPromptTemplate = template.Must(template.New("eval").Parse(evaluationPrompt))

func (r *FunctionalCheckFileResult) Evaluate(ctx context.Context) error {
	examples, err := ListFuncationalCheckExamples(ctx, r.DagName, r.CheckID)
	if err != nil {
		return fmt.Errorf("error listing examples for evaluation: %s", err)
	}

	good := []struct {
		MESSAGE          string
		UPDATED_DOCUMENT string
	}{}
	bad := []struct {
		MESSAGE          string
		UPDATED_DOCUMENT string
	}{}

	for _, example := range examples {
		msg, err := example.OutputMessage()
		if err != nil {
			return fmt.Errorf("error getting output message for example %s: %s", example.ID, err)
		}

		doc, err := example.ResultDoc()
		if err != nil {
			return fmt.Errorf("error getting result doc for example %s: %s", example.ID, err)
		}

		mkdown, err := doc.GetFullMarkdown()
		if err != nil {
			return fmt.Errorf("error getting full markdown for example %s: %s", example.ID, err)
		}

		if example.Approved {
			good = append(good, struct {
				MESSAGE          string
				UPDATED_DOCUMENT string
			}{
				MESSAGE:          msg.FullContent(),
				UPDATED_DOCUMENT: mkdown,
			})
		} else {
			bad = append(bad, struct {
				MESSAGE          string
				UPDATED_DOCUMENT string
			}{
				MESSAGE:          msg.FullContent(),
				UPDATED_DOCUMENT: mkdown,
			})
		}

	}

	msg, err := r.InputMessage()
	if err != nil {
		return fmt.Errorf("error getting output message: %s", err)
	}

	doc, err := r.BeforeDoc()
	if err != nil {
		return fmt.Errorf("error getting result doc: %s", err)
	}

	mkdown, err := doc.GetFullMarkdown()
	if err != nil {
		return fmt.Errorf("error getting full markdown: %s", err)
	}

	outmsg, err := r.OutputMessage()
	if err != nil {
		return fmt.Errorf("error getting output message: %s", err)
	}

	updateDoc, err := r.ResultDoc()
	if err != nil {
		return fmt.Errorf("error getting result doc: %s", err)
	}

	outputMkdown, err := updateDoc.GetFullMarkdown()
	if err != nil {
		return fmt.Errorf("error getting full markdown: %s", err)
	}

	data := EvaluationPromptData{
		USER_MESSAGE:         msg.FullContent(),
		INITIAL_DOCUMENT:     mkdown,
		KNOWN_GOOD_RESPONSES: good,
		KNOWN_BAD_RESPONSES:  bad,
		MESSAGE:              outmsg.FullContent(),
		UPDATED_DOCUMENT:     outputMkdown,
	}

	var buf strings.Builder
	err = evaluationPromptTemplate.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	r.Prompt = buf.String()

	llm, err := anthropic.New(
		anthropic.WithModel("claude-3-5-sonnet-20240620"),
		anthropic.WithToken(os.Getenv("ANTHROPIC_API_KEY")),
	)
	if err != nil {
		return fmt.Errorf("error creating llm: %s", err)
	}

	response, err := llms.GenerateFromSinglePrompt(ctx, llm, buf.String())
	if err != nil {
		return fmt.Errorf("error generating response: %s", err)
	}

	r.RawResponse = response
	var eval Evaluation
	err = xml.Unmarshal([]byte(response), &eval)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %s", err)
	}

	r.Justification = eval.Justification
	r.Score = float64(eval.Score)
	r.Assessment = strings.TrimSpace(eval.Assessment)

	return nil
}
