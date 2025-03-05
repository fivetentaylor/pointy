package v3

import (
	"encoding/json"
	"slices"
	"strings"
	"unicode"
)

const OpsPerStatSegment = 100

var ignoreAuthors = []string{"root", "q"}

type DocStats struct {
	WordCount      int `json:"wordCount"`
	ParagraphCount int `json:"paragraphCount"`
}

func (s *DocStats) AsJS() map[string]interface{} {
	bts, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	out := map[string]interface{}{}
	err = json.Unmarshal(bts, &out)
	if err != nil {
		panic(err)
	}

	return out
}

type OpStats struct {
	Inserts []int `json:"inserts"`
	Deletes []int `json:"deletes"`

	InsertsByPrefix map[string][]int `json:"insertsByPrefix"`
	DeletesByPrefix map[string][]int `json:"deletesByPrefix"`

	CurrentCharsByPrefix map[string]int `json:"currentCharsByPrefix"`

	Segments int `json:"segments"`
}

func (ops *OpStats) AsJS() map[string]interface{} {
	bts, err := json.Marshal(ops)
	if err != nil {
		panic(err)
	}

	out := map[string]interface{}{}
	err = json.Unmarshal(bts, &out)
	if err != nil {
		panic(err)
	}

	return out
}

func (ops *OpStats) AddToInserts(prefix string, segment, inc int) {
	if _, ok := ops.InsertsByPrefix[prefix]; !ok {
		ops.InsertsByPrefix[prefix] = make([]int, ops.Segments)
	}

	ops.Inserts[segment] += inc
	ops.InsertsByPrefix[prefix][segment] += inc
}

func (ops *OpStats) AddToDeletes(prefix string, segment, inc int) {
	if _, ok := ops.DeletesByPrefix[prefix]; !ok {
		ops.DeletesByPrefix[prefix] = make([]int, ops.Segments)
	}

	ops.Deletes[segment] += inc
	ops.DeletesByPrefix[prefix][segment] += inc
}

func (r *Rogue) DocStats() (*DocStats, error) {
	text := r.GetText()
	wordCount := len(strings.Fields(text))

	lines := strings.Split(text, "\n")
	paragraphCount := 0
	inParagraph := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			inParagraph = false
		} else if !inParagraph {
			paragraphCount++
			inParagraph = true
		}
	}

	return &DocStats{
		WordCount:      wordCount,
		ParagraphCount: paragraphCount,
	}, nil
}

func (r *Rogue) OpStats() (*OpStats, error) {
	segments := (r.OpIndex.MaxSeq() / OpsPerStatSegment) + 1
	startID, err := r.Rope.GetTotID(0)
	if err != nil {
		return nil, err
	}

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	if err != nil {
		return nil, err
	}

	tot, err := r.Rope.GetTotBetween(startID, endID)
	if err != nil {
		return nil, err
	}

	out := &OpStats{
		Inserts: make([]int, segments),
		Deletes: make([]int, segments),

		InsertsByPrefix:      make(map[string][]int),
		DeletesByPrefix:      make(map[string][]int),
		CurrentCharsByPrefix: make(map[string]int),

		Segments: segments,
	}

	for author, v := range r.OpIndex.AuthorOps {
		if slices.Contains(ignoreAuthors, author) {
			continue
		}

		var prefix string
		firstAuthorRune := []rune(author)[0]
		if !unicode.IsDigit(firstAuthorRune) && !unicode.IsLetter(firstAuthorRune) {
			prefix = string(firstAuthorRune)
		}

		var processOp func(i Op) error
		processOp = func(i Op) error {
			switch op := i.(type) {
			case InsertOp:
				segment := (op.ID.Seq / OpsPerStatSegment)
				length := len(op.Text)
				out.AddToInserts(prefix, segment, length)
			case DeleteOp:
				segment := (op.ID.Seq / OpsPerStatSegment)
				length := op.SpanLength
				out.AddToDeletes(prefix, segment, length)
			case MultiOp:
				for _, op := range op.Mops {
					err := processOp(op)
					if err != nil {
						return err
					}
				}
			}

			return nil
		}

		err := v.Dft(processOp)

		if err != nil {
			return nil, err
		}
	}

	currentCharsByAuthor := make(map[string]int)
	for idx, id := range tot.IDs {
		if !tot.IsDeleted[idx] {
			currentCharsByAuthor[id.Author]++
		}
	}

	for author, count := range currentCharsByAuthor {
		if slices.Contains(ignoreAuthors, author) {
			continue
		}

		var prefix string
		firstAuthorRune := []rune(author)[0]
		if !unicode.IsDigit(firstAuthorRune) && !unicode.IsLetter(firstAuthorRune) {
			prefix = string(firstAuthorRune)
		}

		out.CurrentCharsByPrefix[prefix] += count
	}

	return out, nil

}
