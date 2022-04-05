package data

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type dynamoDB struct {
	tn string
	ek string
}

var ErrDynamoDBNoEncryptionKey = errors.New("cannot create DynamoDB: fairhive's encryption key is missing")
var ErrDynamoDBNoTableName = errors.New("cannot create DynamoDB: no table name")

func NewDynamoDB(tn, ek string) (db *dynamoDB, err error) {
	if tn == "" {
		return nil, ErrDynamoDBNoTableName
	}
	if ek == "" {
		return nil, ErrDynamoDBNoEncryptionKey
	}
	db = &dynamoDB{
		tn: tn,
		ek: ek,
	}
	return
}

func (db *dynamoDB) Save(u *User) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		return errors.New("cannot create dynamodb client")
	}

	av, err := dynamodbattribute.MarshalMap(*u)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.tn),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}
	fmt.Printf("💾 User [ %v ] saved in DB\n", *u)
	return nil
}
