package data

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fairhive-labs/preregister/internal/crypto/cipher"
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

func (db *dynamoDB) Save(user *User) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		return errors.New("cannot create dynamodb client")
	}

	encEmail, err := cipher.Encrypt(user.Email, db.ek)
	if err != nil {
		return err
	}
	u := NewUser(user.Address, encEmail, user.Type)
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
	fmt.Printf("💾 User [%v] saved in DB\n", *u)
	return nil
}

func (db *dynamoDB) Count() (map[string]int, error) {
	m := map[string]int{
		"advisor":     0,
		"agent":       0,
		"client":      0,
		"contributor": 0,
		"investor":    0,
		"mentor":      0,
		"talent":      0,
	}
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		return nil, errors.New("cannot create dynamodb client")
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(db.tn),
	}
	for {
		result, err := svc.Scan(input)
		if err != nil {
			return nil, err
		}

		for _, u := range result.Items {
			user := User{}
			err = dynamodbattribute.UnmarshalMap(u, &user)
			if err != nil {
				return nil, err
			}
			m[user.Type]++
		}
		// pagination
		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	return m, nil
}

func (db *dynamoDB) List(options ...int) ([]*User, error) {
	users := []*User{}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		return nil, errors.New("cannot create dynamodb client")
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(db.tn),
	}
	for {
		result, err := svc.Scan(input)
		if err != nil {
			return nil, err
		}

		for _, u := range result.Items {
			user := User{}
			err = dynamodbattribute.UnmarshalMap(u, &user)
			if err != nil {
				return nil, err
			}
			e, err := cipher.Decrypt(user.Email, db.ek)
			if err != nil {
				return nil, err
			}
			user.Email = e
			users = append(users, &user)
		}
		// pagination
		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	return users, nil
}
