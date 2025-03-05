package v3

type Selection struct {
	StartSpanID     ID
	StartSpanOffset int
	EndSpanID       ID
	EndSpanOffset   int
}

func (s Selection) AsJS() map[string]interface{} {
	return map[string]interface{}{
		"startSpanID": s.StartSpanID.AsJS(),
		"startOffset": s.StartSpanOffset,
		"endSpanID":   s.EndSpanID.AsJS(),
		"endOffset":   s.EndSpanOffset,
	}
}

// GetSelection returns a selection between two IDs
func (r Rogue) GetSelection(startID, afterID ID, address *ContentAddress, smartQuote bool) (*Selection, error) {
	startSpanID, startSpanOffset, err := r.EnclosingSpanID(startID, address, smartQuote)
	if err != nil {
		return nil, err
	}

	afterSpanID, afterSpanOffset := startSpanID, startSpanOffset

	if startID != afterID {
		afterSpanID, afterSpanOffset, err = r.EnclosingSpanID(afterID, address, smartQuote)
		if err != nil {
			return nil, err
		}
	}

	return &Selection{
		StartSpanID:     startSpanID,
		StartSpanOffset: startSpanOffset,
		EndSpanID:       afterSpanID,
		EndSpanOffset:   afterSpanOffset,
	}, nil
}
