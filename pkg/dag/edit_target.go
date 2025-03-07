package dag

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type EditTargetAction string

const (
	EditTargetActionReplace EditTargetAction = "replace"
	EditTargetActionAppend  EditTargetAction = "append"
	EditTargetActionPrepend EditTargetAction = "prepend"
)

type EditTarget struct {
	ID       string
	BeforeID v3.ID
	AfterID  v3.ID
	Markdown string
	Action   EditTargetAction
}

type innerChunk struct {
	start v3.ID
	end   v3.ID
	size  int
}

// MaxEditDistance calculates the maximum edit distance for levenshtein edit distance
func MaxEditDistance(length int, k float64, b int) int {
	return int(math.Floor(k*float64(length))) + b
}

func ChunkDocument(document *v3.Rogue, maxTokens int) ([]EditTarget, error) {
	chars := document.GetUint16()

	var chunks []innerChunk
	var inGroup bool
	var groupStart int

	addChunk := func(i int) error {
		startId, err := document.Rope.GetVisID(groupStart)
		if err != nil {
			return fmt.Errorf("error getting vis id: %s", err)
		}
		endId, err := document.Rope.GetVisID(i)
		if err != nil {
			return fmt.Errorf("error getting vis id: %s", err)
		}
		mkdown, err := document.GetMarkdown(startId, endId)
		if err != nil {
			return fmt.Errorf("error getting text between: %s", err)
		}
		size := estimateTokens(mkdown)

		chunks = append(chunks, innerChunk{
			start: startId,
			end:   endId,
			size:  size,
		})

		return nil
	}

	for i, c := range chars {
		if utils.IsBreak(c) {
			if inGroup {
				err := addChunk(i)
				if err != nil {
					return nil, fmt.Errorf("error adding chunk: %s", err)
				}
				inGroup = false
			}
		} else {
			if !inGroup {
				groupStart = i
				inGroup = true
			}
		}
	}

	if inGroup {
		err := addChunk(len(chars))
		if err != nil {
			return nil, fmt.Errorf("error adding chunk: %s", err)
		}
	}

	editTargets := []EditTarget{}

	var result []innerChunk
	if len(chunks) == 0 {
		return editTargets, nil
	}

	// Initialize the first group
	currentGroup := innerChunk{
		start: chunks[0].start,
		end:   chunks[0].end,
		size:  chunks[0].size,
	}

	for i := 1; i < len(chunks); i++ {
		chunk := chunks[i]
		// Check if adding the next chunk exceeds the size limit
		if currentGroup.size+chunk.size <= maxTokens {
			// Add to current group
			currentGroup.end = chunk.end
			currentGroup.size += chunk.size
		} else {
			// Finalize the current group and start a new one
			result = append(result, currentGroup)
			currentGroup = innerChunk{
				start: chunk.start,
				end:   chunk.end,
				size:  chunk.size,
			}
		}
	}

	result = append(result, currentGroup)

	for _, chunk := range result {
		mkd, err := document.GetMarkdown(chunk.start, chunk.end)
		if err != nil {
			return nil, fmt.Errorf("error getting markdown: %s", err)
		}
		editTargets = append(editTargets, EditTarget{
			ID:       uuid.New().String(),
			BeforeID: chunk.start,
			AfterID:  chunk.end,
			Markdown: mkd,
		})
	}

	return editTargets, nil
}
