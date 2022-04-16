package data

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	tableName = "Waitlist_TEST"
	ek        = "4e8e7d24d3a991f9e83005d96f8d5d69b4763143a48cf5bdf7941726a26a69ab"
)

func TestListTables(t *testing.T) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	if svc == nil {
		t.Errorf("cannot create dynamodb client")
		t.FailNow()
	}

	tables := map[string]struct{}{}

	input := &dynamodb.ListTablesInput{}
	for {
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					t.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
					t.FailNow()
				default:
					t.Error(aerr.Error())
					t.FailNow()
				}
			} else {
				t.Error(aerr.Error())
				t.FailNow()
			}
			return
		}
		for _, n := range result.TableNames {
			tables[*n] = struct{}{}
		}
		input.ExclusiveStartTableName = result.LastEvaluatedTableName
		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	if _, ok := tables[tableName]; !ok {
		t.Errorf("%s is not listed in DynamoDB tables", tableName)
		t.FailNow()
	}
}

func TestNewDynamoDB(t *testing.T) {
	tt := []struct {
		name   string
		tn, ek string
		err    error
	}{
		{"normal", tableName, ek, nil},
		{"no table name", "", ek, ErrDynamoDBNoTableName},
		{"no encryption key", tableName, "", ErrDynamoDBNoEncryptionKey},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDynamoDB(tc.tn, tc.ek)
			if !errors.Is(err, tc.err) {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.FailNow()
			}
		})
	}
}

func TestSave(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email := "john.doe@mailservice.com"
	utype := "talent"
	u := &User{
		Address: address,
		Email:   email,
		Type:    utype,
	}

	db, _ := NewDynamoDB(tableName, ek)
	if err := db.Save(u); err != nil {
		t.Errorf("cannot save user %v: %v", *u, err)
		t.FailNow()
	}

	u = &User{
		Address: "",
		Email:   email,
		Type:    utype,
	}
	if err := db.Save(u); err == nil {
		t.Errorf("impossible to save user with an empty string Address: ValidationException")
		t.FailNow()
	}
}

func TestCount(t *testing.T) {
	db, _ := NewDynamoDB(tableName, ek)
	mc, err := db.Count()
	if err != nil {
		t.Errorf("cannot count users: %v", err)
		t.FailNow()
	}
	if mc["talent"] == 0 {
		t.Errorf("incorrect talent count: must be greater than 0")
		t.FailNow()
	}
}
