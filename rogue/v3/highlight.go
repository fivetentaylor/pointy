package v3

import "fmt"

func (r *Rogue) HighlightSpan(beforeID, afterID ID) (int, int, error) {
	startID, err := r.VisRightOf(beforeID)
	if err != nil {
		return 0, 0, fmt.Errorf("VisRightOf(%v): %w", beforeID, err)
	}

	endID := afterID
	endChar, err := r.GetCharByID(afterID)
	if err != nil || endChar != '\n' {
		endID, err = r.VisLeftOf(afterID)
		if err != nil {
			return 0, 0, fmt.Errorf("VisLeftOf(%v): %w", afterID, err)
		}
	}

	startIx, _, err := r.Rope.GetIndex(startID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	endIx, _, err := r.Rope.GetIndex(endID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	return startIx, endIx - startIx + 1, nil
}

func (r *Rogue) SpanContainsID(startID, endID, id ID) bool {
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return false
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return false
	}

	_, idIx, err := r.Rope.GetIndex(id)
	if err != nil {
		return false
	}

	return startIx <= idIx && idIx <= endIx
}
