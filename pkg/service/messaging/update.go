package messaging

import (
	"context"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func UpdateMessage(ctx context.Context, message *dynamo.Message) error {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)

	err := dydb.UpdateMessage(message)
	if err != nil {
		log.Errorf("error updating message: %s", err)
		return err
	}

	err = PublishMessage(ctx, message)
	if err != nil {
		log.Errorf("error publishing message: %s", err)
		return err
	}

	return nil
}

func UpdateThread(ctx context.Context, thread *dynamo.Thread) error {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)

	err := dydb.UpdateThread(thread)
	if err != nil {
		log.Errorf("error updating message: %s", err)
		return err
	}

	err = PublishThread(ctx, thread)
	if err != nil {
		log.Errorf("error publishing message: %s", err)
		return err
	}

	return nil
}
