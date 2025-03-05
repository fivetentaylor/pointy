package dynamo

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Pagination parameters
type PaginationParams struct {
	Limit             int64                               // Number of items to be returned in each page
	ExclusiveStartKey map[string]*dynamodb.AttributeValue // Key to start with for the next page
}
