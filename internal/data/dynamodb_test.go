package data

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
	db, _ := NewDynamoDB(tableName, ek)

	u := NewUser(sponsor, "jsie@trendev.fr", "mentor", sponsor) // sponsor is first user and its own sponsor
	if err := db.Save(u); err != nil {
		t.Errorf("impossible to save default sponsor: %v", err)
		t.FailNow()
	}

	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email := "john.doe@mailservice.com"
	utype := "talent"
	u = &User{
		Address: address,
		Email:   email,
		Type:    utype,
		Sponsor: sponsor,
	}

	if err := db.Save(u); err != nil {
		t.Errorf("cannot save user %v: %v", *u, err)
		t.FailNow()
	}

	u = &User{
		Address: "",
		Email:   email,
		Type:    utype,
		Sponsor: sponsor,
	}
	if err := db.Save(u); err == nil {
		t.Errorf("impossible to save user with an empty string Address: ValidationException")
		t.FailNow()
	}

	u = &User{
		Address: address,
		Email:   email,
		Type:    utype,
		Sponsor: "",
	}
	if err := db.Save(u); err == nil {
		t.Errorf("impossible to save user with an empty string Sponsor: ValidationException")
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

func TestList(t *testing.T) {
	db, _ := NewDynamoDB(tableName, ek)
	t.Run("no option", func(t *testing.T) {
		users, err := db.List()
		if err != nil {
			t.Errorf("cannot list users: %v", err)
			t.FailNow()
		}
		if len(users) == 0 {
			t.Errorf("users list cannot be nil or empty")
			t.FailNow()
		}

		sort.Slice(users, func(i, j int) bool {
			if users[i].Timestamp == users[j].Timestamp {
				return strings.Compare(users[i].Email, users[i].Email) < 1
			}
			return users[i].Timestamp > users[j].Timestamp // more recent first
		})
		var johndoe *User
		for i, u := range users {
			fmt.Printf("%0.2d - %s\n", i+1, u)
			if u.Address == "0x8ba1f109551bD432803012645Ac136ddd64DBA72" &&
				u.Email == "john.doe@mailservice.com" &&
				u.Type == "talent" {
				johndoe = u
			}
		}
		if johndoe == nil {
			t.Errorf("johndoe must be present in DB")
			t.FailNow()
		}
	})

	tt := []struct {
		offset, max int
		len         int
		err         error
	}{
		{0, 5, 5, nil},
		{0, 10, 10, nil},
		{10, 5, 5, nil},
		{2, 3, 3, nil},
		{0, -1, 0, ErrBadMax},
		{0, 1000, 61, nil}, //  2023-05-13: 61 items in the test db
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("offset=%d max=%d", tc.offset, tc.max), func(t *testing.T) {
			users, err := db.List(tc.offset, tc.max)
			if err != tc.err {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.FailNow()
			}
			if len(users) != tc.len {
				t.Errorf("incorrect len(users), got %v, want %v", len(users), tc.len)
				t.FailNow()
			}
		})
	}
}

func TestIsPresent(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	db, _ := NewDynamoDB(tableName, ek)

	tt := []struct {
		a string
		r bool
	}{
		{address, true},
		{sponsor, true},
		{"fake4adr3ss", false},
	}

	for _, tc := range tt {
		t.Run(tc.a, func(t *testing.T) {
			r, err := db.IsPresent(tc.a)
			if err != nil {
				t.Errorf("cannot test if user %s is present or not: %v", tc.a, err)
				t.FailNow()
			}
			if r != tc.r {
				t.Errorf("incorrect IsPresent(%s) result, got %v, want %v", tc.a, r, tc.r)
				t.FailNow()
			}
		})
	}
}
