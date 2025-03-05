package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var ContentAddressPrefix = "addr#"

type dContentAddress struct {
	PK    string `dynamodbav:"PK" json:"PK"`
	SK    string `dynamodbav:"SK" json:"SK"`
	Bytes []byte `dynamodbav:"payload" json:"payload"`
}

func (db *DB) GetContentAddressIDs(docID string, params *PaginationParams) ([]string, *PaginationParams, error) {
	pk := ContentAddressPrefix + docID

	if params == nil {
		params = &PaginationParams{
			Limit:             100,
			ExclusiveStartKey: nil,
		}
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
		},
		Limit:             aws.Int64(params.Limit),
		ExclusiveStartKey: params.ExclusiveStartKey,
		ScanIndexForward:  aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return nil, nil, fmt.Errorf("Query(%q, %+v): %s", docID, params, err)
	}

	addressIDs := make([]string, 0, len(results.Items))
	for _, item := range results.Items {
		ca := dContentAddress{}
		err := dynamodbattribute.UnmarshalMap(item, &ca)
		if err != nil {
			return nil, nil, fmt.Errorf("Query(%q, %+v): %s", docID, params, err)
		}
		addressIDs = append(addressIDs, ca.SK)
	}

	newParams := PaginationParams{
		Limit:             params.Limit,
		ExclusiveStartKey: results.LastEvaluatedKey,
	}

	return addressIDs, &newParams, nil
}

func (db *DB) GetContentAddress(documentID, addressID string) ([]byte, error) {
	item, err := db.GetItem(
		ContentAddressPrefix+documentID,
		addressID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetItem(%v, %v): %w", ContentAddressPrefix+documentID, addressID, err)
	}

	dca := dContentAddress{}
	err = dynamodbattribute.UnmarshalMap(item, &dca)
	if err != nil {
		return nil, fmt.Errorf("dynamodbattribute.UnmarshalMap(%v): %w", item, err)
	}

	return dca.Bytes, nil
}

func (db *DB) CreateContentAddress(documentID string, address []byte) (string, error) {
	addressID := hashAndEncode(address, 16)
	ca := dContentAddress{
		PK:    ContentAddressPrefix + documentID,
		SK:    addressID,
		Bytes: address,
	}

	av, err := dynamodbattribute.Marshal(ca)
	if err != nil {
		return "", fmt.Errorf("dynamodbattribute.Marshal(%v): %w", ca, err)
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      av.M,
	}

	_, err = db.Client.PutItem(putInput)
	if err != nil {
		return "", fmt.Errorf("PutItem(%v): %w", putInput, err)
	}
	return addressID, nil
}
