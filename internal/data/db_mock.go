package data

import (
	"errors"
	"fmt"

	key "github.com/fairhive-labs/ethkeygen/pkg"
)

type mockDB struct {
}

func (db mockDB) Save(u *User) (err error) {
	fmt.Printf("💾 User [ %v ] saved in DB\n", *u)
	return
}

func (db mockDB) Count() (map[string]int, error) {
	m := map[string]int{
		"advisor":     0,
		"agent":       2,
		"client":      0,
		"contributor": 0,
		"investor":    0,
		"mentor":      1,
		"talent":      3,
	}
	return m, nil
}

func (db mockDB) List(offset, max int) ([]*User, error) {
	m := map[string]int{
		"advisor":     1,
		"agent":       5,
		"client":      7,
		"contributor": 0,
		"investor":    10,
		"mentor":      5,
		"talent":      31,
	}

	users := []*User{}
	for k, v := range m {
		for i := 0; i < v; i++ {
			_, a, _ := key.Generate()
			u := NewUser(a, fmt.Sprintf("%s_%d@domain.com", k, (i+1)), k)
			users = append(users, u)
		}
	}
	if offset < 0 || offset > len(users) {
		return nil, fmt.Errorf("incorrect offset")
	}
	if max < 0 || max > len(users) {
		return nil, fmt.Errorf("incorrect max")
	}
	if offset+max > len(users) {
		return nil, fmt.Errorf("ouf of bounds [%d:%d]", offset, offset+max)
	}
	return users[offset : offset+max], nil
}

var MockDB = mockDB{}

type mockErrDB struct {
}

func (db mockErrDB) Save(u *User) (err error) {
	m := fmt.Sprintf("🔥 Error saving User [ %v ] in DB\n", *u)
	fmt.Print(m)
	return errors.New(m)
}

func (db mockErrDB) Count() (map[string]int, error) {
	m := fmt.Sprintf("🔥 Error counting Users in DB\n")
	fmt.Print(m)
	return nil, errors.New(m)
}

func (db mockErrDB) List(offset, max int) ([]*User, error) {
	m := fmt.Sprintf("🔥 Error listing Users in DB\n")
	fmt.Print(m)
	return nil, errors.New(m)
}

var MockErrDB = mockErrDB{}
