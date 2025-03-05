package dag

import (
	"context"
	"encoding/json"
	"fmt"

	v3 "github.com/teamreviso/code/rogue/v3"
)

func PublishOp(ctx context.Context, docId string, op v3.Op) error {
	docStore, err := GetDocStore(ctx)
	if err != nil {
		return fmt.Errorf("error getting doc store: %s", err)
	}

	realtime, err := GetDocumentRealtimeConnection(ctx, docId)
	if err != nil {
		return fmt.Errorf("error getting realtime: %s", err)
	}

	_, err = docStore.AddUpdate(ctx, docId, op)
	if err != nil {
		return fmt.Errorf("error adding update: %s", err)
	}

	bts, err := json.Marshal(op)
	err = realtime.PublishOp(docId, bts)
	if err != nil {
		return fmt.Errorf("error applying op: %s", err)
	}

	return nil
}
