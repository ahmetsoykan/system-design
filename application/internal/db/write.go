package db

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"url-shortener/internal/hash"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName = "urls"
)

// DDBInit works only on startup
// It puts an item to database to track IDs for new URLs and auto-increment the IDs.
func (d *DynamoDB) DDBInit() {

	item := Item{
		ID:      "Metadata",
		LongURL: "Metadata",
		LastID:  100000, // pre-defined start ID
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("error: DDBPutItem: %s", err.Error())
	}
	// Create Metadata Item
	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}

	_, err = d.client.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.Printf("main : database has initialized already")
			default:
				log.Fatalf("error: database couldn't initialized: %s", err.Error())
			}
		}
	} else {
		log.Printf("main : database initialized")
	}

}

// DDBPutItem inserts an URL item and increments the lastid of the metadata item.
func (d *DynamoDB) DDBPutItem(item Item, hostValue string) (Item, error) {

	// Existance Check
	existingItem, err := d.GetItembyGSI(item.LongURL)
	if err != nil {
		return item, err
	} else {
		if existingItem.LongURL == item.LongURL {
			return existingItem, nil
		}
	}

	//* New Record Creation
	// Find the last inserted item's ID, and encode the next ID for the new URL.
	lastid, _ := d.GetLastIDbyMetadataPK()
	item.ID = fmt.Sprintf("%d", lastid+1)
	item.ShortURL = fmt.Sprintf("%s/%s", hostValue, hash.Encode(lastid+1))
	item.Short = fmt.Sprintf("%s", hash.Encode(lastid+1))

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("error: DDBPutItem: %s", err.Error())
	}
	// Put Item
	// Increment the ID count in Metadata item
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName: &tableName,
					Item:      av,
				},
			},
			{
				Update: &dynamodb.Update{
					TableName: &tableName,
					Key: map[string]*dynamodb.AttributeValue{
						"id": {
							S: aws.String("Metadata"),
						},
					},
					UpdateExpression:    aws.String("SET lastid = lastid + :incr"),
					ConditionExpression: aws.String("lastid = :current_id"),
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":incr": {
							N: aws.String("1"),
						},
						":current_id": {
							N: aws.String(strconv.FormatUint(lastid, 10)),
						},
					},
				},
			},
		},
	}

	_, err = d.client.TransactWriteItems(input)
	if err != nil {
		log.Fatalf("error: DDBTransactWriteItems: %s", err.Error())
	}

	return item, nil

}

// When URL visited, extend the TTL 90 days more
func (d *DynamoDB) ExtendTTL(short string) error {

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Update: &dynamodb.Update{
					TableName: &tableName,
					Key: map[string]*dynamodb.AttributeValue{
						"id": {
							S: aws.String(short),
						},
					},
					UpdateExpression:    aws.String("SET #ttl = :_ttl"),
					ConditionExpression: aws.String("defaultttl = :defaultttl"),
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":_ttl": {
							N: aws.String(strconv.FormatInt(time.Now().UTC().Add(time.Hour*24*90).Unix(), 10)),
						},
						":defaultttl": {
							BOOL: aws.Bool(true),
						},
					},
					ExpressionAttributeNames: map[string]*string{
						"#ttl": aws.String("ttl"),
					},
				},
			},
		},
	}

	// TODO: error handling
	_, _ = d.client.TransactWriteItems(input)

	return nil
}
