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

var (
	ErrDynamoDBNoEncryptionKey = errors.New("cannot create DynamoDB: poln's encryption key is missing")
	ErrDynamoDBNoTableName     = errors.New("cannot create DynamoDB: no table name")
	ErrBadMax                  = errors.New("incorrect max")
	ErrInvalidUser             = errors.New("nil user or missing required field")
)

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

func (db *dynamoDB) IsPresent(a string) (bool, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	r, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.tn),
		Key: map[string]*dynamodb.AttributeValue{
			"address": {
				S: aws.String(a),
			},
		},
	})
	if err != nil {
		return false, err
	}
	return r.Item != nil, nil
}

func (db *dynamoDB) Save(u *User) error {
	if u == nil || !u.IsSet() {
		return ErrInvalidUser
	}
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		return errors.New("cannot create dynamodb client")
	}

	encEmail, err := cipher.Encrypt(u.Email, db.ek)
	if err != nil {
		return err
	}
	u2 := NewUser(u.Address, encEmail, u.Type, u.Sponsor)
	av, err := dynamodbattribute.MarshalMap(*u2)
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
	fmt.Printf("💾 User saved in DB: [%v]\n", *u2)
	*u = *u2 // copy saved user
	return nil
}

func (db *dynamoDB) Count() (map[string]int, error) {
	m := map[string]int{
		"advisor":     0,
		"agent":       0,
		"initiator":   0,
		"contributor": 0,
		"investor":    0,
		"mentor":      0,
		"contractor":  0,
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

	var max *int64
	if len(options) == 2 {
		// offset ignored
		max = new(int64)
		*max = int64(options[1])
	}
	if max != nil && *max < 0 {
		return nil, ErrBadMax
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(db.tn),
		Limit:     max,
	}
	for {
		if input.Limit != nil && *input.Limit == 0 {
			break
		}
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
		if input.Limit != nil {
			*input.Limit = *input.Limit - *result.ScannedCount
		}
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	return users, nil
}
