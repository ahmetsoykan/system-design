package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	client *dynamodb.DynamoDB
}

func NewDynamoDBClient(region string) *DynamoDB {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.New(sess)
	return &DynamoDB{client: svc}
}

type Item struct {
	ID         string `json:"id"`
	LongURL    string `json:"longurl"`
	ShortURL   string `json:"shorturl"`
	LastID     uint64 `json:"lastid"`
	TTL        int64  `json:"ttl"`
	DefaultTTL bool   `json:"defaultttl"`
}
