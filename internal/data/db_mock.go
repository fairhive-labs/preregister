package data

import (
	"errors"
	"fmt"

	key "github.com/fairhive-labs/ethkeygen/pkg"
)

var (
	UsersMapMock = map[string]int{
		"advisor":     1,
		"agent":       5,
		"client":      7,
		"contributor": 0,
		"investor":    10,
		"mentor":      5,
		"talent":      31,
	}
	UsersCountMock = 0
)

func init() {
	for _, v := range UsersMapMock {
		UsersCountMock += v
	}
}

type mockDB struct {
}

func (db mockDB) Save(u *User) (err error) {
	fmt.Printf("💾 User [ %v ] saved in DB\n", *u)
	return
}

func (db mockDB) Count() (map[string]int, error) {
	m := UsersMapMock
	return m, nil
}

func (db mockDB) List(options ...int) ([]*User, error) {
	m := UsersMapMock
	users := []*User{}
	for k, v := range m {
		for i := 0; i < v; i++ {
			_, a, _ := key.Generate() // user's address
			_, s, _ := key.Generate() // user's sponsor
			u := NewUser(a, fmt.Sprintf("%s_%d@domain.com", k, (i+1)), k, s)
			users = append(users, u)
		}
	}

	offset, max := 0, len(users)
	if len(options) >= 1 {
		offset = options[0]
		max = len(users) - offset
	}
	if len(options) == 2 {
		max = options[1]
	}
	if offset < 0 || offset > len(users) {
		return nil, fmt.Errorf("incorrect offset")
	}
	if max < 0 {
		return nil, fmt.Errorf("incorrect max")
	}
	if max > len(users) {
		max = len(users) - offset
	}

	if offset+max > len(users) {
		return nil, fmt.Errorf("ouf of bounds [%d:%d]", offset, offset+max)
	}

	return users[offset : offset+max], nil
}

func (db mockDB) IsPresent(a string) (bool, error) {
	return true, nil
}

var MockDB = mockDB{}

type mockErrDB struct {
	mockDB
}

func (db mockErrDB) Save(u *User) (err error) {
	m := fmt.Sprintf("🔥 Error saving User [ %v ] in DB\n", *u)
	fmt.Print(m)
	return errors.New(m)
}

func (db mockErrDB) Count() (map[string]int, error) {
	m := "🔥 Error counting Users in DB"
	fmt.Println(m)
	return nil, errors.New(m)
}

func (db mockErrDB) List(options ...int) ([]*User, error) {
	m := "🔥 Error listing Users in DB"
	fmt.Println(m)
	return nil, errors.New(m)
}

var MockErrDB = mockErrDB{MockDB}
