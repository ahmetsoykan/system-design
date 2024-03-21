package db

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Returning an item given the partition key
func (d *DynamoDB) GetItembyPK(pk string) (Item, error) {
	result, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(pk),
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("error: DDBGetItm: %s", err.Error())
	}

	// Unmarhsall
	item := Item{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
		log.Fatalf("error: DDBGetItm: %s", err.Error())
	}

	return item, nil
}

// findLastID returns the last inserted item's id.
// Using it before putting the next item's id.
func (d *DynamoDB) GetLastIDbyMetadataPK() (uint64, error) {
	result, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("Metadata"),
			},
		},
	})
	if err != nil {
		log.Fatalf("error: DDBGetItm: %s", err.Error())
	}

	// Unmarhsall
	item := Item{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
		log.Fatalf("error: DDBGetItm: %s", err.Error())
	}

	return item.LastID, nil
}

// Returning existing item if the longurl exists in DDB
func (d *DynamoDB) GetItembyGSI(rk string) (Item, error) {

	item := Item{}

	result, err := d.client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("gsi1"),
		Limit:                  aws.Int64(1),
		KeyConditionExpression: aws.String("longurl = :longurl"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":longurl": {
				S: aws.String(rk),
			},
		},
	})
	if err != nil {
		return item, errors.New("error: DDBGSIGetItem: " + err.Error())
	}

	// Unmarhsall
	if len(result.Items) > 0 {
		if err := dynamodbattribute.UnmarshalMap(result.Items[0], &item); err != nil {
			return item, errors.New("error: DDBGSIGetItem: " + err.Error())
		}
	}

	return item, nil
}
