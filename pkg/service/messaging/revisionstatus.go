package messaging

import (
	"context"
	"time"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/graph/model"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UpdateRevisionStatus(ctx context.Context, containerID, messageID string, status model.MessageRevisionStatus, contentAddress string) (*dynamo.Message, error) {
	msg, err := env.Dynamo(ctx).GetMessage(containerID, messageID)
	if err != nil {
		return nil, err
	}

	switch status {
	case model.MessageRevisionStatusAccepted:
		msg.MessageMetadata.RevisionStatus = models.RevisionStatus_REVISION_STATUS_ACCEPTED
	case model.MessageRevisionStatusDeclined:
		msg.MessageMetadata.RevisionStatus = models.RevisionStatus_REVISION_STATUS_DECLINED
	}

	msg.MessageMetadata.ContentAddressAfter = contentAddress
	msg.MessageMetadata.ContentAddressAfterTimestamp = timestamppb.New(time.Now())

	err = env.Dynamo(ctx).UpdateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
