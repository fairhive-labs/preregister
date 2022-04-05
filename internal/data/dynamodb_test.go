package data

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "Users_TEST"

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
